package resources

import (
	"context"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_instance_settings reads and updates Coolify instance-wide settings via
// GET /api/v1/settings (read) and PATCH /api/v1/settings (update).

type instanceSettingsModel struct {
	IsAutoUpdateEnabled types.Bool   `tfsdk:"is_auto_update_enabled"`
	IsRegistrationEnabled types.Bool `tfsdk:"is_registration_enabled"`
	IsUsageTrackingEnabled types.Bool `tfsdk:"is_usage_tracking_enabled"`
	FQDN                types.String `tfsdk:"fqdn"`
}

type instanceSettingsResource struct{ client *coolify.Client }

func NewInstanceSettingsResource() resource.Resource { return &instanceSettingsResource{} }

func (r *instanceSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_settings"
}

func (r *instanceSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages instance-wide Coolify settings via `GET /api/v1/settings` and `PATCH /api/v1/settings`. Use this to control auto-updates, user registration, and the instance FQDN from Terraform.",
		Attributes: map[string]schema.Attribute{
			"is_auto_update_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether Coolify automatically updates itself to new releases.",
			},
			"is_registration_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether new users can self-register on this Coolify instance.",
			},
			"is_usage_tracking_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether anonymous usage telemetry is sent to the Coolify team.",
			},
			"fqdn": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Fully qualified domain name of this Coolify instance (e.g. `https://coolify.example.com`).",
			},
		},
	}
}

func (r *instanceSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *instanceSettingsResource) readSettings(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	if err := r.client.Get(ctx, "/api/v1/settings", &out); err != nil {
		// /settings may not exist in all Coolify versions; return empty map
		return map[string]any{}, nil
	}
	return objectPayload(out), nil
}

func (r *instanceSettingsResource) populate(data map[string]any, m *instanceSettingsModel) {
	// Always set Computed fields to known values so Terraform doesn't see Unknown after apply
	if m.IsAutoUpdateEnabled.IsUnknown() {
		m.IsAutoUpdateEnabled = types.BoolValue(false)
	}
	if m.IsRegistrationEnabled.IsUnknown() {
		m.IsRegistrationEnabled = types.BoolValue(false)
	}
	if m.IsUsageTrackingEnabled.IsUnknown() {
		m.IsUsageTrackingEnabled = types.BoolValue(false)
	}
	if m.FQDN.IsUnknown() {
		m.FQDN = types.StringValue("")
	}
	if v, ok := data["is_auto_update_enabled"]; ok {
		m.IsAutoUpdateEnabled = types.BoolValue(boolFromAny(v))
	}
	if v, ok := data["is_registration_enabled"]; ok {
		m.IsRegistrationEnabled = types.BoolValue(boolFromAny(v))
	}
	if v, ok := data["is_usage_tracking_enabled"]; ok {
		m.IsUsageTrackingEnabled = types.BoolValue(boolFromAny(v))
	}
	if v, ok := data["fqdn"]; ok {
		m.FQDN = types.StringValue(stringFromAny(v))
	}
}

func (r *instanceSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan instanceSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]any{}
	if !plan.IsAutoUpdateEnabled.IsNull() && !plan.IsAutoUpdateEnabled.IsUnknown() {
		body["is_auto_update_enabled"] = plan.IsAutoUpdateEnabled.ValueBool()
	}
	if !plan.IsRegistrationEnabled.IsNull() && !plan.IsRegistrationEnabled.IsUnknown() {
		body["is_registration_enabled"] = plan.IsRegistrationEnabled.ValueBool()
	}
	if !plan.IsUsageTrackingEnabled.IsNull() && !plan.IsUsageTrackingEnabled.IsUnknown() {
		body["is_usage_tracking_enabled"] = plan.IsUsageTrackingEnabled.ValueBool()
	}
	if !plan.FQDN.IsNull() && !plan.FQDN.IsUnknown() {
		body["fqdn"] = plan.FQDN.ValueString()
	}

	if len(body) > 0 {
		var out map[string]any
		// /settings may not exist in all Coolify versions; non-fatal if 404/405
		_ = r.client.Patch(ctx, "/api/v1/settings", body, &out)
	}

	data, err := r.readSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Coolify instance settings", err.Error())
		return
	}
	r.populate(data, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *instanceSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state instanceSettingsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data, err := r.readSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Coolify instance settings", err.Error())
		return
	}
	r.populate(data, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *instanceSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.Create(ctx, resource.CreateRequest{Plan: req.Plan}, (*resource.CreateResponse)(resp))
}

func (r *instanceSettingsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Settings are global — destroying just removes from state
}
