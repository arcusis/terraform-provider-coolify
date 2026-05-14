package resources

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// compositeResource handles sub-resources whose IDs are "parentUUID/childUUID".
// CRUD paths: list=GET /parent/parentUUID/child, item=GET|PATCH|DELETE /parent/parentUUID/child/childUUID
type compositeResource struct {
	client      *coolify.Client
	typeName    string
	displayName string
	parentField string // the field name in the schema that holds the parent UUID
	listPath    func(parentUUID string) string
	createPath  func(parentUUID string) string
	updatePath  func(parentUUID, childUUID string) string
	deletePath  func(parentUUID, childUUID string) string
	fields      []resourceField
}

func (r *compositeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName
}

func (r *compositeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{Computed: true, Description: "Composite ID: parentUUID/childUUID."},
	}
	for _, f := range r.fields {
		attrs[f.Name] = f.schemaAttribute()
	}
	resp.Schema = schema.Schema{Attributes: attrs}
}

func (r *compositeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*coolify.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *coolify.Client.")
		return
	}
	r.client = client
}

func (r *compositeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	vals := r.readAttrs(ctx, req.Plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	parentUUID := r.parentUUID(vals)
	body := r.bodyFromVals(vals)

	var created map[string]any
	if err := r.client.Post(ctx, r.createPath(parentUUID), body, &created); err != nil {
		resp.Diagnostics.AddError("Unable to create Coolify "+r.displayName, err.Error())
		return
	}

	childUUID := firstStringFromMap(created, "uuid", "id")
	if childUUID == "" {
		resp.Diagnostics.AddError("Create response missing UUID", r.displayName+" create did not return uuid")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(parentUUID+"/"+childUUID))...)
	r.writeAttrsToState(ctx, &resp.State, vals)

	// Re-read from API to get computed fields
	apiData, found := r.findByUUID(ctx, parentUUID, childUUID)
	if found {
		r.writeAPIDataToState(ctx, &resp.State, apiData, &resp.Diagnostics)
	}
}

func (r *compositeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentUUID, childUUID := splitCompositeID2(id.ValueString())
	if parentUUID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	apiData, found := r.findByUUID(ctx, parentUUID, childUUID)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	r.writeAPIDataToState(ctx, &resp.State, apiData, &resp.Diagnostics)
}

func (r *compositeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentUUID, childUUID := splitCompositeID2(id.ValueString())
	vals := r.readAttrs(ctx, req.Plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Patch(ctx, r.updatePath(parentUUID, childUUID), r.bodyFromVals(vals), nil); err != nil {
		resp.Diagnostics.AddError("Unable to update Coolify "+r.displayName, err.Error())
		return
	}

	apiData, found := r.findByUUID(ctx, parentUUID, childUUID)
	if found {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
		r.writeAPIDataToState(ctx, &resp.State, apiData, &resp.Diagnostics)
	}
}

func (r *compositeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentUUID, childUUID := splitCompositeID2(id.ValueString())
	if err := r.client.Delete(ctx, r.deletePath(parentUUID, childUUID), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Coolify "+r.displayName, err.Error())
	}
}

func (r *compositeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func (r *compositeResource) parentUUID(vals map[string]any) string {
	if v, ok := vals[r.parentField]; ok {
		return fmt.Sprint(v)
	}
	return ""
}

func (r *compositeResource) readAttrs(ctx context.Context, plan planOrState, diags *diag.Diagnostics) map[string]any {
	vals := make(map[string]any)
	for _, f := range r.fields {
		switch f.Kind {
		case kindString:
			var v types.String
			diags.Append(plan.GetAttribute(ctx, path.Root(f.Name), &v)...)
			if !v.IsNull() && !v.IsUnknown() {
				vals[f.Name] = v.ValueString()
			}
		case kindBool:
			var v types.Bool
			diags.Append(plan.GetAttribute(ctx, path.Root(f.Name), &v)...)
			if !v.IsNull() && !v.IsUnknown() {
				vals[f.Name] = v.ValueBool()
			}
		case kindInt64:
			var v types.Int64
			diags.Append(plan.GetAttribute(ctx, path.Root(f.Name), &v)...)
			if !v.IsNull() && !v.IsUnknown() {
				vals[f.Name] = v.ValueInt64()
			}
		}
	}
	return vals
}

func (r *compositeResource) bodyFromVals(vals map[string]any) map[string]any {
	body := make(map[string]any)
	for _, f := range r.fields {
		if !f.Send || f.Name == r.parentField {
			continue
		}
		if v, ok := vals[f.Name]; ok {
			body[f.Name] = v
		}
	}
	return body
}

func (r *compositeResource) writeAttrsToState(ctx context.Context, state stateTarget, vals map[string]any) {
	for _, f := range r.fields {
		v, ok := vals[f.Name]
		if !ok {
			continue
		}
		switch f.Kind {
		case kindString:
			state.SetAttribute(ctx, path.Root(f.Name), types.StringValue(fmt.Sprint(v)))
		case kindBool:
			if b, ok := v.(bool); ok {
				state.SetAttribute(ctx, path.Root(f.Name), types.BoolValue(b))
			}
		case kindInt64:
			state.SetAttribute(ctx, path.Root(f.Name), types.Int64Value(int64FromAny(v)))
		}
	}
}

func (r *compositeResource) writeAPIDataToState(ctx context.Context, state stateTarget, data map[string]any, diags *diag.Diagnostics) {
	for _, f := range r.fields {
		if f.SkipAPIRead {
			continue
		}
		v, ok := data[f.Name]
		if !ok || v == nil {
			continue
		}
		if f.Sensitive && stringFromAny(v) == "" {
			continue
		}
		switch f.Kind {
		case kindString:
			diags.Append(state.SetAttribute(ctx, path.Root(f.Name), types.StringValue(stringFromAny(v)))...)
		case kindBool:
			diags.Append(state.SetAttribute(ctx, path.Root(f.Name), types.BoolValue(boolFromAny(v)))...)
		case kindInt64:
			diags.Append(state.SetAttribute(ctx, path.Root(f.Name), types.Int64Value(int64FromAny(v)))...)
		}
	}
}

func (r *compositeResource) findByUUID(ctx context.Context, parentUUID, childUUID string) (map[string]any, bool) {
	var out any
	if err := r.client.Get(ctx, r.listPath(parentUUID), &out); err != nil {
		return nil, false
	}
	for _, item := range listFromAny(out) {
		if firstStringFromMap(item, "uuid", "id") == childUUID {
			return item, true
		}
	}
	return nil, false
}

func splitCompositeID2(id string) (string, string) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func listFromAny(out any) []map[string]any {
	switch v := out.(type) {
	case []any:
		return mapListFromAny(v)
	case map[string]any:
		if nested, ok := v["data"].([]any); ok {
			return mapListFromAny(nested)
		}
		return []map[string]any{objectPayload(v)}
	default:
		return nil
	}
}

func mapListFromAny(items []any) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			result = append(result, m)
		}
	}
	return result
}
