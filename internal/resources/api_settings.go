package resources

import (
	"context"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_api_settings toggles the Coolify API on/off via
// GET /api/v1/enable and GET /api/v1/disable.
// Creating with enabled=true calls /enable; enabled=false calls /disable.
// Destroying calls /enable to restore the API.

type apiSettingsModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type apiSettingsResource struct{ client *coolify.Client }

func NewAPISettingsResource() resource.Resource { return &apiSettingsResource{} }

func (r *apiSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_settings"
}

func (r *apiSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages whether the Coolify REST API is enabled or disabled. Setting `enabled = false` calls `GET /api/v1/disable`, which prevents all API access until re-enabled. The resource restores `enabled = true` on destroy to avoid locking yourself out.",
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Whether the Coolify API is enabled. Set to `false` to disable, `true` to enable.",
			},
		},
	}
}

func (r *apiSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *apiSettingsResource) applyState(ctx context.Context, enabled bool) error {
	path := "/api/v1/enable"
	if !enabled {
		path = "/api/v1/disable"
	}
	var out any
	err := r.client.Get(ctx, path, &out)
	if err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && (httpErr.StatusCode == 403 || httpErr.StatusCode == 404) {
			// Non-admin tokens may get 403; some versions may not have these endpoints
			return nil
		}
	}
	return err
}

func (r *apiSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.applyState(ctx, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to set Coolify API state", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// The API has no GET for the current enabled/disabled state; keep existing state.
	var state apiSettingsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *apiSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan apiSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.applyState(ctx, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to update Coolify API state", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiSettingsResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Restore enabled state on destroy to avoid self-lockout
	_ = r.applyState(ctx, true)
}
