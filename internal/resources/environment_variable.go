package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type envVarModel struct {
	ID              types.String `tfsdk:"id"`
	ApplicationUUID types.String `tfsdk:"application_uuid"`
	Key             types.String `tfsdk:"key"`
	Value           types.String `tfsdk:"value"`
	IsPreview       types.Bool   `tfsdk:"is_preview"`
	IsLiteral       types.Bool   `tfsdk:"is_literal"`
	IsMultiline     types.Bool   `tfsdk:"is_multiline"`
	IsShownOnce     types.Bool   `tfsdk:"is_shown_once"`
}

type environmentVariableResource struct {
	client *coolify.Client
}

func NewEnvironmentVariableResource() resource.Resource {
	return &environmentVariableResource{}
}

func (r *environmentVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_variable"
}

func (r *environmentVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Computed: true},
			"application_uuid": schema.StringAttribute{Required: true},
			"key":              schema.StringAttribute{Required: true},
			"value":            schema.StringAttribute{Required: true, Sensitive: true},
			"is_preview":       schema.BoolAttribute{Optional: true, Computed: true},
			"is_literal":       schema.BoolAttribute{Optional: true, Computed: true},
			"is_multiline":     schema.BoolAttribute{Optional: true, Computed: true},
			"is_shown_once":    schema.BoolAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *environmentVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *environmentVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan envVarModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var created map[string]any
	if err := r.client.Post(ctx, envsPath(plan.ApplicationUUID.ValueString()), envVarBody(plan), &created); err != nil {
		resp.Diagnostics.AddError("Unable to create Coolify environment variable", err.Error())
		return
	}

	applyEnvVarData(&plan, created)
	id := firstStringFromMap(created, "uuid", "id")
	if id == "" {
		found, err := r.find(ctx, plan.ApplicationUUID.ValueString(), "", plan.Key.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Unable to read env var after create", err.Error())
			return
		}
		id = firstStringFromMap(found, "uuid", "id")
		applyEnvVarData(&plan, found)
	}
	if id == "" {
		id = plan.Key.ValueString()
	}
	// Ensure computed bool fields are known (default to false if not returned by API)
	if plan.IsPreview.IsNull() || plan.IsPreview.IsUnknown() {
		plan.IsPreview = types.BoolValue(false)
	}
	if plan.IsLiteral.IsNull() || plan.IsLiteral.IsUnknown() {
		plan.IsLiteral = types.BoolValue(false)
	}
	if plan.IsMultiline.IsNull() || plan.IsMultiline.IsUnknown() {
		plan.IsMultiline = types.BoolValue(false)
	}
	if plan.IsShownOnce.IsNull() || plan.IsShownOnce.IsUnknown() {
		plan.IsShownOnce = types.BoolValue(false)
	}
	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state envVarModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found, err := r.find(ctx, state.ApplicationUUID.ValueString(), state.ID.ValueString(), state.Key.ValueString())
	if err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify environment variable", err.Error())
		return
	}
	if len(found) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}
	applyEnvVarData(&state, found)
	if id := firstStringFromMap(found, "uuid", "id"); id != "" {
		state.ID = types.StringValue(id)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan envVarModel
	var state envVarModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Patch(ctx, envsPath(plan.ApplicationUUID.ValueString()), envVarBody(plan), nil); err != nil {
		resp.Diagnostics.AddError("Unable to update Coolify environment variable", err.Error())
		return
	}
	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}
	found, err := r.find(ctx, plan.ApplicationUUID.ValueString(), plan.ID.ValueString(), plan.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read env var after update", err.Error())
		return
	}
	applyEnvVarData(&plan, found)
	if id := firstStringFromMap(found, "uuid", "id"); id != "" {
		plan.ID = types.StringValue(id)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state envVarModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.ApplicationUUID.ValueString() == "" || state.ID.ValueString() == "" {
		return
	}
	if err := r.client.Delete(ctx, fmt.Sprintf("/api/v1/applications/%s/envs/%s",
		state.ApplicationUUID.ValueString(), state.ID.ValueString()), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Coolify environment variable", err.Error())
	}
}

func (r *environmentVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *environmentVariableResource) find(ctx context.Context, appUUID, envUUID, key string) (map[string]any, error) {
	var out any
	if err := r.client.Get(ctx, envsPath(appUUID), &out); err != nil {
		return nil, err
	}
	for _, item := range envList(out) {
		if envUUID != "" && firstStringFromMap(item, "uuid", "id") == envUUID {
			return item, nil
		}
		if key != "" && stringFromAny(item["key"]) == key {
			return item, nil
		}
	}
	return map[string]any{}, nil
}

func envVarBody(m envVarModel) map[string]any {
	body := map[string]any{
		"key":   m.Key.ValueString(),
		"value": m.Value.ValueString(),
	}
	if !m.IsPreview.IsNull() && !m.IsPreview.IsUnknown() {
		body["is_preview"] = m.IsPreview.ValueBool()
	}
	if !m.IsLiteral.IsNull() && !m.IsLiteral.IsUnknown() {
		body["is_literal"] = m.IsLiteral.ValueBool()
	}
	if !m.IsMultiline.IsNull() && !m.IsMultiline.IsUnknown() {
		body["is_multiline"] = m.IsMultiline.ValueBool()
	}
	if !m.IsShownOnce.IsNull() && !m.IsShownOnce.IsUnknown() {
		body["is_shown_once"] = m.IsShownOnce.ValueBool()
	}
	return body
}

func applyEnvVarData(m *envVarModel, data map[string]any) {
	if v, ok := data["key"]; ok {
		m.Key = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["value"]; ok {
		m.Value = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["is_preview"]; ok {
		m.IsPreview = types.BoolValue(boolFromAny(v))
	}
	if v, ok := data["is_literal"]; ok {
		m.IsLiteral = types.BoolValue(boolFromAny(v))
	}
	if v, ok := data["is_multiline"]; ok {
		m.IsMultiline = types.BoolValue(boolFromAny(v))
	}
	if v, ok := data["is_shown_once"]; ok {
		m.IsShownOnce = types.BoolValue(boolFromAny(v))
	}
}

func envsPath(appUUID string) string {
	return fmt.Sprintf("/api/v1/applications/%s/envs", appUUID)
}

func envList(out any) []map[string]any {
	switch v := out.(type) {
	case []any:
		return mapListEnv(v)
	case map[string]any:
		for _, key := range []string{"data", "envs", "environment_variables"} {
			if arr, ok := v[key].([]any); ok {
				return mapListEnv(arr)
			}
		}
		return []map[string]any{objectPayload(v)}
	default:
		return nil
	}
}

func mapListEnv(items []any) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			result = append(result, m)
		}
	}
	return result
}
