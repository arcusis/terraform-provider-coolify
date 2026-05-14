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

type dbEnvVarModel struct {
	ID           types.String `tfsdk:"id"`
	DatabaseUUID types.String `tfsdk:"database_uuid"`
	Key          types.String `tfsdk:"key"`
	Value        types.String `tfsdk:"value"`
	IsLiteral    types.Bool   `tfsdk:"is_literal"`
	IsMultiline  types.Bool   `tfsdk:"is_multiline"`
	IsShownOnce  types.Bool   `tfsdk:"is_shown_once"`
}

type databaseEnvironmentVariableResource struct{ client *coolify.Client }

func NewDatabaseEnvironmentVariableResource() resource.Resource {
	return &databaseEnvironmentVariableResource{}
}

func (r *databaseEnvironmentVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_environment_variable"
}

func (r *databaseEnvironmentVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true},
			"database_uuid": schema.StringAttribute{Required: true},
			"key":           schema.StringAttribute{Required: true},
			"value":         schema.StringAttribute{Required: true, Sensitive: true},
			"is_literal":    schema.BoolAttribute{Optional: true, Computed: true},
			"is_multiline":  schema.BoolAttribute{Optional: true, Computed: true},
			"is_shown_once": schema.BoolAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *databaseEnvironmentVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *databaseEnvironmentVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dbEnvVarModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := dbEnvBody(plan)
	var created map[string]any
	if err := r.client.Post(ctx, dbEnvsPath(plan.DatabaseUUID.ValueString()), body, &created); err != nil {
		resp.Diagnostics.AddError("Unable to create Coolify database environment variable", err.Error())
		return
	}

	applyDBEnvData(&plan, created)
	id := firstStringFromMap(created, "uuid", "id")
	if id == "" {
		found, err := r.findEnv(ctx, plan.DatabaseUUID.ValueString(), "", plan.Key.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Unable to read db env var after create", err.Error())
			return
		}
		id = firstStringFromMap(found, "uuid", "id")
		applyDBEnvData(&plan, found)
	}
	if id == "" {
		id = plan.Key.ValueString()
	}
	for _, ptr := range []*types.Bool{&plan.IsLiteral, &plan.IsMultiline, &plan.IsShownOnce} {
		if ptr.IsNull() || ptr.IsUnknown() {
			*ptr = types.BoolValue(false)
		}
	}
	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *databaseEnvironmentVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dbEnvVarModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	found, err := r.findEnv(ctx, state.DatabaseUUID.ValueString(), state.ID.ValueString(), state.Key.ValueString())
	if err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify database environment variable", err.Error())
		return
	}
	if len(found) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}
	applyDBEnvData(&state, found)
	if id := firstStringFromMap(found, "uuid", "id"); id != "" {
		state.ID = types.StringValue(id)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *databaseEnvironmentVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state dbEnvVarModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Patch(ctx, dbEnvsPath(plan.DatabaseUUID.ValueString()), dbEnvBody(plan), nil); err != nil {
		resp.Diagnostics.AddError("Unable to update Coolify database environment variable", err.Error())
		return
	}
	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}
	found, err := r.findEnv(ctx, plan.DatabaseUUID.ValueString(), plan.ID.ValueString(), plan.Key.ValueString())
	if err == nil {
		applyDBEnvData(&plan, found)
		if id := firstStringFromMap(found, "uuid", "id"); id != "" {
			plan.ID = types.StringValue(id)
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *databaseEnvironmentVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dbEnvVarModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.DatabaseUUID.ValueString() == "" || state.ID.ValueString() == "" {
		return
	}
	if err := r.client.Delete(ctx, fmt.Sprintf("/api/v1/databases/%s/envs/%s",
		state.DatabaseUUID.ValueString(), state.ID.ValueString()), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Coolify database environment variable", err.Error())
	}
}

func (r *databaseEnvironmentVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *databaseEnvironmentVariableResource) findEnv(ctx context.Context, dbUUID, envUUID, key string) (map[string]any, error) {
	var out any
	if err := r.client.Get(ctx, dbEnvsPath(dbUUID), &out); err != nil {
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

func dbEnvsPath(uuid string) string { return fmt.Sprintf("/api/v1/databases/%s/envs", uuid) }

func dbEnvBody(m dbEnvVarModel) map[string]any {
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

func applyDBEnvData(m *dbEnvVarModel, data map[string]any) {
	if v, ok := data["key"]; ok {
		m.Key = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["value"]; ok {
		m.Value = types.StringValue(stringFromAny(v))
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
