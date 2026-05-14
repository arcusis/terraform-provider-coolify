package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ── Field descriptor ─────────────────────────────────────────────────────────

type fieldKind string

const (
	kindString fieldKind = "string"
	kindBool   fieldKind = "bool"
	kindInt64  fieldKind = "int64"
)

type resourceField struct {
	Name      string
	Kind      fieldKind
	Required  bool
	Optional  bool
	Computed  bool
	Sensitive bool
	Send      bool
}

func stringField(name string, required, optional, sensitive bool) resourceField {
	return resourceField{Name: name, Kind: kindString, Required: required, Optional: optional, Sensitive: sensitive, Send: true}
}

func boolField(name string, required, optional, sensitive bool) resourceField {
	return resourceField{Name: name, Kind: kindBool, Required: required, Optional: optional, Sensitive: sensitive, Send: true}
}

func int64Field(name string, required, optional, sensitive bool) resourceField {
	return resourceField{Name: name, Kind: kindInt64, Required: required, Optional: optional, Sensitive: sensitive, Send: true}
}

func computedStringField(name string) resourceField {
	return resourceField{Name: name, Kind: kindString, Computed: true, Send: false}
}

func (f resourceField) schemaAttribute() schema.Attribute {
	computed := f.Computed || (f.Optional && !f.Required)
	switch f.Kind {
	case kindBool:
		return schema.BoolAttribute{Required: f.Required, Optional: f.Optional, Computed: computed, Sensitive: f.Sensitive}
	case kindInt64:
		return schema.Int64Attribute{Required: f.Required, Optional: f.Optional, Computed: computed, Sensitive: f.Sensitive}
	default:
		return schema.StringAttribute{Required: f.Required, Optional: f.Optional, Computed: computed, Sensitive: f.Sensitive}
	}
}

// ── Generic resource ──────────────────────────────────────────────────────────

type genericResource struct {
	client              *coolify.Client
	typeName            string
	displayName         string
	createPath          func(vals map[string]string) string
	readPath            func(id string) string
	updatePath          func(id string) string
	deletePath          func(id string) string
	fields              []resourceField
	updateBodyTransform func(id string, body map[string]any) map[string]any
}

func newGenericResource(typeName, displayName, createPath, itemPath string, fields []resourceField) resource.Resource {
	return &genericResource{
		typeName:    typeName,
		displayName: displayName,
		createPath:  func(map[string]string) string { return createPath },
		readPath:    func(id string) string { return fmt.Sprintf(itemPath, id) },
		updatePath:  func(id string) string { return fmt.Sprintf(itemPath, id) },
		deletePath:  func(id string) string { return fmt.Sprintf(itemPath, id) },
		fields:      fields,
	}
}

func (r *genericResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName
}

func (r *genericResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{Computed: true, Description: "Coolify resource UUID."},
	}
	for _, f := range r.fields {
		attrs[f.Name] = f.schemaAttribute()
	}
	resp.Schema = schema.Schema{Attributes: attrs}
}

func (r *genericResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *genericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	vals := r.readAttrsFromPlan(ctx, req.Plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	body := r.bodyFromVals(vals)
	var created map[string]any
	if err := r.client.Post(ctx, r.createPath(r.stringVals(vals)), body, &created); err != nil {
		resp.Diagnostics.AddError("Unable to create Coolify "+r.displayName, err.Error())
		return
	}

	id := firstStringFromMap(created, "uuid", "id")
	if id == "" {
		resp.Diagnostics.AddError("Create response missing UUID", "Coolify did not return uuid for "+r.displayName)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
	r.writeAttrsToState(ctx, &resp.State, vals)

	apiData := r.fetchByID(ctx, id, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
	r.writeAPIDataToState(ctx, &resp.State, apiData, &resp.Diagnostics)
}

func (r *genericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() || id.IsNull() || id.ValueString() == "" {
		return
	}

	apiData, diags := r.fetchByIDWithDiags(ctx, id.ValueString())
	if isNotFoundDiag(diags) {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id.ValueString()))...)
	r.writeAPIDataToState(ctx, &resp.State, apiData, &resp.Diagnostics)
}

func (r *genericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vals := r.readAttrsFromPlan(ctx, req.Plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	body := r.bodyFromVals(vals)
	if r.updateBodyTransform != nil {
		body = r.updateBodyTransform(id.ValueString(), body)
	}
	if err := r.client.Patch(ctx, r.updatePath(id.ValueString()), body, nil); err != nil {
		resp.Diagnostics.AddError("Unable to update Coolify "+r.displayName, err.Error())
		return
	}

	apiData := r.fetchByID(ctx, id.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	r.writeAPIDataToState(ctx, &resp.State, apiData, &resp.Diagnostics)
}

func (r *genericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() || id.IsNull() || id.ValueString() == "" {
		return
	}
	if err := r.client.Delete(ctx, r.deletePath(id.ValueString()), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Coolify "+r.displayName, err.Error())
	}
}

func (r *genericResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// ── Plan/state helpers ────────────────────────────────────────────────────────

type planOrState interface {
	GetAttribute(ctx context.Context, p path.Path, target interface{}) diag.Diagnostics
}

func (r *genericResource) readAttrsFromPlan(ctx context.Context, plan planOrState, diags *diag.Diagnostics) map[string]any {
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

func (r *genericResource) bodyFromVals(vals map[string]any) map[string]any {
	body := make(map[string]any)
	for _, f := range r.fields {
		if !f.Send {
			continue
		}
		if v, ok := vals[f.Name]; ok {
			body[f.Name] = v
		}
	}
	return body
}

func (r *genericResource) stringVals(vals map[string]any) map[string]string {
	out := make(map[string]string)
	for k, v := range vals {
		out[k] = fmt.Sprint(v)
	}
	return out
}

type stateTarget interface {
	SetAttribute(ctx context.Context, p path.Path, val interface{}) diag.Diagnostics
}


func (r *genericResource) writeAttrsToState(ctx context.Context, state stateTarget, vals map[string]any) {
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

func (r *genericResource) writeAPIDataToState(ctx context.Context, state stateTarget, data map[string]any, diags *diag.Diagnostics) {
	for _, f := range r.fields {
		v, ok := data[f.Name]
		if !ok || v == nil {
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

func (r *genericResource) fetchByID(ctx context.Context, id string, diags *diag.Diagnostics) map[string]any {
	data, d := r.fetchByIDWithDiags(ctx, id)
	diags.Append(d...)
	return data
}

func (r *genericResource) fetchByIDWithDiags(ctx context.Context, id string) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics
	var out map[string]any
	if err := r.client.Get(ctx, r.readPath(id), &out); err != nil {
		diags.AddError("Unable to read Coolify "+r.displayName, err.Error())
		return nil, diags
	}
	return objectPayload(out), diags
}

// ── Shared helpers ────────────────────────────────────────────────────────────

func objectPayload(data map[string]any) map[string]any {
	if nested, ok := data["data"].(map[string]any); ok {
		return nested
	}
	return data
}

func firstStringFromMap(data map[string]any, keys ...string) string {
	data = objectPayload(data)
	for _, key := range keys {
		if value, ok := data[key]; ok {
			if text := stringFromAny(value); text != "" {
				return text
			}
		}
	}
	return ""
}

func stringFromAny(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		return strconv.FormatBool(v)
	case []any, map[string]any:
		payload, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprint(v)
		}
		return string(payload)
	default:
		return fmt.Sprint(v)
	}
}

func boolFromAny(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		parsed, _ := strconv.ParseBool(v)
		return parsed
	case float64:
		return v != 0
	default:
		return false
	}
}

func int64FromAny(value any) int64 {
	switch v := value.(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	case string:
		parsed, _ := strconv.ParseInt(v, 10, 64)
		return parsed
	default:
		return 0
	}
}

func isNotFoundDiag(diags diag.Diagnostics) bool {
	for _, d := range diags {
		if strings.Contains(d.Detail(), "HTTP 404") {
			return true
		}
	}
	return false
}

// ── Project resource ──────────────────────────────────────────────────────────

func NewProjectResource() resource.Resource {
	return newGenericResource("project", "project", "/api/v1/projects", "/api/v1/projects/%s", []resourceField{
		stringField("name", true, false, false),
		stringField("description", false, true, false),
	})
}
