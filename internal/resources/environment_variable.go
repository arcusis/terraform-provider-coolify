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
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Coolify environment variable UUID.",
			},
			"application_uuid": schema.StringAttribute{
				Required:    true,
				Description: "Coolify application UUID.",
			},
			"key": schema.StringAttribute{
				Required:    true,
				Description: "Environment variable key.",
			},
			"value": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Environment variable value.",
			},
			"is_preview": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether this environment variable applies to preview deployments.",
			},
			"is_literal": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether this environment variable is literal.",
			},
			"is_multiline": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether this environment variable is multiline.",
			},
			"is_shown_once": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether this environment variable is shown once.",
			},
		},
	}
}

func (r *environmentVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*coolify.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *coolify.Client from provider configuration.")
		return
	}
	r.client = client
}

func (r *environmentVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var patched map[string]any
	if err := r.client.Patch(ctx, envsPath(plan.ApplicationUUID.ValueString()), r.body(plan), &patched); err != nil {
		resp.Diagnostics.AddError("Unable to create Coolify environment variable", err.Error())
		return
	}

	id := firstString(patched, "uuid", "id")
	if id == "" {
		found, err := r.find(ctx, plan.ApplicationUUID.ValueString(), "", plan.Key.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Unable to read Coolify environment variable after create", err.Error())
			return
		}
		id = firstString(found, "uuid", "id")
		applyEnvironmentVariable(&plan, found)
	}
	if id == "" {
		id = plan.Key.ValueString()
	}
	plan.ID = types.StringValue(id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appUUID := state.ApplicationUUID.ValueString()
	if appUUID == "" {
		resp.Diagnostics.AddError("Missing application UUID", "Cannot read Coolify environment variable because application_uuid is empty.")
		return
	}
	found, err := r.find(ctx, appUUID, state.ID.ValueString(), state.Key.ValueString())
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

	applyEnvironmentVariable(&state, found)
	if id := firstString(found, "uuid", "id"); id != "" {
		state.ID = types.StringValue(id)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceModel
	var state resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Patch(ctx, envsPath(plan.ApplicationUUID.ValueString()), r.body(plan), nil); err != nil {
		resp.Diagnostics.AddError("Unable to update Coolify environment variable", err.Error())
		return
	}
	if plan.ID.IsNull() || plan.ID.IsUnknown() || plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}
	found, err := r.find(ctx, plan.ApplicationUUID.ValueString(), plan.ID.ValueString(), plan.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Coolify environment variable after update", err.Error())
		return
	}
	applyEnvironmentVariable(&plan, found)
	if id := firstString(found, "uuid", "id"); id != "" {
		plan.ID = types.StringValue(id)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appUUID := state.ApplicationUUID.ValueString()
	envUUID := state.ID.ValueString()
	if appUUID == "" || envUUID == "" {
		return
	}
	if err := r.client.Delete(ctx, fmt.Sprintf("/api/v1/applications/%s/envs/%s", appUUID, envUUID), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Coolify environment variable", err.Error())
	}
}

func (r *environmentVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *environmentVariableResource) body(model resourceModel) map[string]any {
	body := make(map[string]any)
	for _, name := range []string{"key", "value", "is_preview", "is_literal", "is_multiline", "is_shown_once"} {
		if value, ok := modelFieldValue(model, name); ok {
			body[name] = value
		}
	}
	return body
}

func (r *environmentVariableResource) find(ctx context.Context, appUUID, envUUID, key string) (map[string]any, error) {
	var out any
	if err := r.client.Get(ctx, envsPath(appUUID), &out); err != nil {
		return nil, err
	}
	for _, item := range envList(out) {
		if envUUID != "" && (firstString(item, "uuid", "id") == envUUID) {
			return item, nil
		}
		if key != "" && stringFromAny(item["key"]) == key {
			return item, nil
		}
	}
	return map[string]any{}, nil
}

func envsPath(appUUID string) string {
	return fmt.Sprintf("/api/v1/applications/%s/envs", appUUID)
}

func envList(out any) []map[string]any {
	switch v := out.(type) {
	case []any:
		return mapList(v)
	case map[string]any:
		if data, ok := v["data"].([]any); ok {
			return mapList(data)
		}
		if envs, ok := v["envs"].([]any); ok {
			return mapList(envs)
		}
		if envs, ok := v["environment_variables"].([]any); ok {
			return mapList(envs)
		}
		return []map[string]any{objectPayload(v)}
	default:
		return nil
	}
}

func mapList(items []any) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if mapped, ok := item.(map[string]any); ok {
			result = append(result, mapped)
		}
	}
	return result
}

func applyEnvironmentVariable(model *resourceModel, data map[string]any) {
	for _, name := range []string{"application_uuid", "key", "value", "is_preview", "is_literal", "is_multiline", "is_shown_once"} {
		if value, ok := data[name]; ok {
			setModelField(model, name, value)
		}
	}
}
