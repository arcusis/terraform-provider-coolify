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

type serviceEnvVarModel struct {
	ID          types.String `tfsdk:"id"`
	ServiceUUID types.String `tfsdk:"service_uuid"`
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
	IsLiteral   types.Bool   `tfsdk:"is_literal"`
	IsMultiline types.Bool   `tfsdk:"is_multiline"`
	IsShownOnce types.Bool   `tfsdk:"is_shown_once"`
}

type serviceEnvironmentVariableResource struct{ client *coolify.Client }

func NewServiceEnvironmentVariableResource() resource.Resource {
	return &serviceEnvironmentVariableResource{}
}

func (r *serviceEnvironmentVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_environment_variable"
}

func (r *serviceEnvironmentVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true},
			"service_uuid": schema.StringAttribute{Required: true},
			"key":          schema.StringAttribute{Required: true},
			"value":        schema.StringAttribute{Required: true, Sensitive: true},
			"is_literal":   schema.BoolAttribute{Optional: true, Computed: true},
			"is_multiline": schema.BoolAttribute{Optional: true, Computed: true},
			"is_shown_once": schema.BoolAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *serviceEnvironmentVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *serviceEnvironmentVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceEnvVarModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := serviceEnvBody(plan)
	var created map[string]any
	if err := r.client.Post(ctx, serviceEnvsPath(plan.ServiceUUID.ValueString()), body, &created); err != nil {
		// 409 means Coolify pre-populated this env var from the service template.
		// Fall back to upsert: find the existing one and PATCH it.
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusConflict {
			existing, findErr := r.findEnv(ctx, plan.ServiceUUID.ValueString(), "", plan.Key.ValueString())
			if findErr != nil || len(existing) == 0 {
				resp.Diagnostics.AddError("Unable to create Coolify service environment variable", err.Error())
				return
			}
			if patchErr := r.client.Patch(ctx, serviceEnvsPath(plan.ServiceUUID.ValueString()), serviceEnvBody(plan), nil); patchErr != nil {
				resp.Diagnostics.AddError("Unable to upsert Coolify service environment variable", patchErr.Error())
				return
			}
			// Use existing ID but keep plan's sensitive value — the API masks it on read.
			if id := firstStringFromMap(existing, "uuid", "id"); id != "" {
				created = map[string]any{"uuid": id}
			}
		} else {
			resp.Diagnostics.AddError("Unable to create Coolify service environment variable", err.Error())
			return
		}
	}

	applyServiceEnvData(&plan, created)
	id := firstStringFromMap(created, "uuid", "id")
	if id == "" {
		found, err := r.findEnv(ctx, plan.ServiceUUID.ValueString(), "", plan.Key.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Unable to read service env var after create", err.Error())
			return
		}
		id = firstStringFromMap(found, "uuid", "id")
		applyServiceEnvData(&plan, found)
	}
	if id == "" {
		id = plan.Key.ValueString()
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

func (r *serviceEnvironmentVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceEnvVarModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	found, err := r.findEnv(ctx, state.ServiceUUID.ValueString(), state.ID.ValueString(), state.Key.ValueString())
	if err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify service environment variable", err.Error())
		return
	}
	if len(found) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}
	applyServiceEnvData(&state, found)
	if id := firstStringFromMap(found, "uuid", "id"); id != "" {
		state.ID = types.StringValue(id)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *serviceEnvironmentVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state serviceEnvVarModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Patch(ctx, serviceEnvsPath(plan.ServiceUUID.ValueString()), serviceEnvBody(plan), nil); err != nil {
		resp.Diagnostics.AddError("Unable to update Coolify service environment variable", err.Error())
		return
	}
	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}
	found, err := r.findEnv(ctx, plan.ServiceUUID.ValueString(), plan.ID.ValueString(), plan.Key.ValueString())
	if err == nil {
		applyServiceEnvData(&plan, found)
		if id := firstStringFromMap(found, "uuid", "id"); id != "" {
			plan.ID = types.StringValue(id)
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serviceEnvironmentVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceEnvVarModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.ServiceUUID.ValueString() == "" || state.ID.ValueString() == "" {
		return
	}
	if err := r.client.Delete(ctx, fmt.Sprintf("/api/v1/services/%s/envs/%s",
		state.ServiceUUID.ValueString(), state.ID.ValueString()), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Coolify service environment variable", err.Error())
	}
}

func (r *serviceEnvironmentVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *serviceEnvironmentVariableResource) findEnv(ctx context.Context, svcUUID, envUUID, key string) (map[string]any, error) {
	var out any
	if err := r.client.Get(ctx, serviceEnvsPath(svcUUID), &out); err != nil {
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

func serviceEnvsPath(uuid string) string { return fmt.Sprintf("/api/v1/services/%s/envs", uuid) }

func serviceEnvBody(m serviceEnvVarModel) map[string]any {
	body := map[string]any{"key": m.Key.ValueString(), "value": m.Value.ValueString()}
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

func applyServiceEnvData(m *serviceEnvVarModel, data map[string]any) {
	if v, ok := data["key"]; ok {
		m.Key = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["value"]; ok {
		s := stringFromAny(v)
		if s != "" {
			m.Value = types.StringValue(s)
		}
		// else: keep plan value — API masks sensitive values after they're set
	}
	if v, ok := data["is_literal"]; ok {
		m.IsLiteral = types.BoolValue(boolFromAny(v))
	}
	if v, ok := data["is_multiline"]; ok {
		m.IsMultiline = types.BoolValue(boolFromAny(v))
	}
	if v, ok := data["is_shown_once"]; ok {
		// is_shown_once=true means the value was hidden after write;
		// the API returns false afterwards. Keep true if we set it.
		if !m.IsShownOnce.IsNull() && m.IsShownOnce.ValueBool() {
			return
		}
		m.IsShownOnce = types.BoolValue(boolFromAny(v))
	}
}
