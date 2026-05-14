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

type githubAppModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	APIURL         types.String `tfsdk:"api_url"`
	HTMLURL        types.String `tfsdk:"html_url"`
	AppID          types.Int64  `tfsdk:"app_id"`
	InstallationID types.Int64  `tfsdk:"installation_id"`
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	PrivateKeyUUID types.String `tfsdk:"private_key_uuid"`
	Organization   types.String `tfsdk:"organization"`
	WebhookSecret  types.String `tfsdk:"webhook_secret"`
	IsSystemWide   types.Bool   `tfsdk:"is_system_wide"`
}

type githubAppResource struct{ client *coolify.Client }

func NewGitHubAppResource() resource.Resource { return &githubAppResource{} }

func (r *githubAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_github_app"
}

func (r *githubAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Computed: true},
			"name":             schema.StringAttribute{Required: true},
			"api_url":          schema.StringAttribute{Required: true},
			"html_url":         schema.StringAttribute{Required: true},
			"app_id":           schema.Int64Attribute{Required: true},
			"installation_id":  schema.Int64Attribute{Required: true},
			"client_id":        schema.StringAttribute{Required: true},
			"client_secret":    schema.StringAttribute{Required: true, Sensitive: true},
			"private_key_uuid": schema.StringAttribute{Required: true},
			"organization":     schema.StringAttribute{Optional: true, Computed: true},
			"webhook_secret":   schema.StringAttribute{Required: true, Sensitive: true},
			"is_system_wide":   schema.BoolAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *githubAppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *githubAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan githubAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]any{
		"name":             plan.Name.ValueString(),
		"api_url":          plan.APIURL.ValueString(),
		"html_url":         plan.HTMLURL.ValueString(),
		"app_id":           plan.AppID.ValueInt64(),
		"installation_id":  plan.InstallationID.ValueInt64(),
		"client_id":        plan.ClientID.ValueString(),
		"client_secret":    plan.ClientSecret.ValueString(),
		"private_key_uuid": plan.PrivateKeyUUID.ValueString(),
	}
	if !plan.Organization.IsNull() && !plan.Organization.IsUnknown() {
		body["organization"] = plan.Organization.ValueString()
	}
	body["webhook_secret"] = plan.WebhookSecret.ValueString()
	if !plan.IsSystemWide.IsNull() && !plan.IsSystemWide.IsUnknown() {
		body["is_system_wide"] = plan.IsSystemWide.ValueBool()
	}

	var created map[string]any
	if err := r.client.Post(ctx, "/api/v1/github-apps", body, &created); err != nil {
		resp.Diagnostics.AddError("Unable to create Coolify GitHub App", err.Error())
		return
	}

	id := firstStringFromMap(created, "uuid", "id")
	if id == "" {
		resp.Diagnostics.AddError("Create response missing ID", "GitHub App create did not return an ID.")
		return
	}

	plan.ID = types.StringValue(id)
	if plan.Organization.IsNull() || plan.Organization.IsUnknown() {
		plan.Organization = types.StringValue("")
	}
	if plan.IsSystemWide.IsNull() || plan.IsSystemWide.IsUnknown() {
		plan.IsSystemWide = types.BoolValue(false)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *githubAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state githubAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No single-item GET — list and filter
	var out any
	if err := r.client.Get(ctx, "/api/v1/github-apps", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify GitHub Apps", err.Error())
		return
	}

	found := r.findByID(out, state.ID.ValueString())
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if v, ok := found["name"]; ok {
		state.Name = types.StringValue(stringFromAny(v))
	}
	if v, ok := found["api_url"]; ok {
		state.APIURL = types.StringValue(stringFromAny(v))
	}
	if v, ok := found["html_url"]; ok {
		state.HTMLURL = types.StringValue(stringFromAny(v))
	}
	if v, ok := found["app_id"]; ok {
		state.AppID = types.Int64Value(int64FromAny(v))
	}
	if v, ok := found["installation_id"]; ok {
		state.InstallationID = types.Int64Value(int64FromAny(v))
	}
	if v, ok := found["client_id"]; ok {
		state.ClientID = types.StringValue(stringFromAny(v))
	}
	if v, ok := found["private_key_uuid"]; ok {
		state.PrivateKeyUUID = types.StringValue(stringFromAny(v))
	}
	if v, ok := found["organization"]; ok && v != nil {
		state.Organization = types.StringValue(stringFromAny(v))
	}
	if v, ok := found["is_system_wide"]; ok {
		state.IsSystemWide = types.BoolValue(boolFromAny(v))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *githubAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state githubAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]any{
		"name":             plan.Name.ValueString(),
		"client_id":        plan.ClientID.ValueString(),
		"client_secret":    plan.ClientSecret.ValueString(),
		"private_key_uuid": plan.PrivateKeyUUID.ValueString(),
	}
	if err := r.client.Patch(ctx, fmt.Sprintf("/api/v1/github-apps/%s", state.ID.ValueString()), body, nil); err != nil {
		resp.Diagnostics.AddError("Unable to update Coolify GitHub App", err.Error())
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *githubAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state githubAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, fmt.Sprintf("/api/v1/github-apps/%s", state.ID.ValueString()), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Coolify GitHub App", err.Error())
	}
}

func (r *githubAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *githubAppResource) findByID(out any, id string) map[string]any {
	for _, item := range listFromAny(out) {
		if firstStringFromMap(item, "uuid", "id") == id {
			return item
		}
	}
	return nil
}
