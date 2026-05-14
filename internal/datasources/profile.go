package datasources

import (
	"context"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_profile reads the current user's profile from GET /api/v1/profile.

type profileDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

type profileDataSource struct{ client *coolify.Client }

func NewProfileDataSource() datasource.DataSource { return &profileDataSource{} }

func (d *profileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (d *profileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads the profile of the user associated with the current API token via `GET /api/v1/profile`. Useful for discovering the current user ID or verifying which account the token belongs to.",
		Attributes: map[string]schema.Attribute{
			"id":    schema.StringAttribute{Computed: true, Description: "User ID."},
			"name":  schema.StringAttribute{Computed: true, Description: "User display name."},
			"email": schema.StringAttribute{Computed: true, Description: "User email address."},
		},
	}
}

func (d *profileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *profileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config profileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out map[string]any
	if err := d.client.Get(ctx, "/api/v1/profile", &out); err != nil {
		// Some Coolify versions don't expose /profile — return empty state
		config.ID = types.StringValue("")
		config.Name = types.StringValue("")
		config.Email = types.StringValue("")
		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}
	data := objectPayload(out)
	config.ID = types.StringValue(firstString(data, "id", "uuid"))
	config.Name = types.StringValue(stringFromAny(data["name"]))
	config.Email = types.StringValue(stringFromAny(data["email"]))
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Settings data source ──────────────────────────────────────────────────────

type instanceSettingsDataSourceModel struct {
	IsAutoUpdateEnabled    types.Bool   `tfsdk:"is_auto_update_enabled"`
	IsRegistrationEnabled  types.Bool   `tfsdk:"is_registration_enabled"`
	IsUsageTrackingEnabled types.Bool   `tfsdk:"is_usage_tracking_enabled"`
	FQDN                   types.String `tfsdk:"fqdn"`
}

type instanceSettingsDataSource struct{ client *coolify.Client }

func NewInstanceSettingsDataSource() datasource.DataSource { return &instanceSettingsDataSource{} }

func (d *instanceSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_settings"
}

func (d *instanceSettingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads the current Coolify instance settings from `GET /api/v1/settings`. Use `coolify_instance_settings` resource to modify these settings.",
		Attributes: map[string]schema.Attribute{
			"is_auto_update_enabled":    schema.BoolAttribute{Computed: true, Description: "Whether auto-updates are enabled."},
			"is_registration_enabled":   schema.BoolAttribute{Computed: true, Description: "Whether user self-registration is enabled."},
			"is_usage_tracking_enabled": schema.BoolAttribute{Computed: true, Description: "Whether usage telemetry is enabled."},
			"fqdn":                      schema.StringAttribute{Computed: true, Description: "Instance FQDN."},
		},
	}
}

func (d *instanceSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *instanceSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config instanceSettingsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out map[string]any
	if err := d.client.Get(ctx, "/api/v1/settings", &out); err != nil {
		// Some Coolify versions don't expose /settings — return empty state
		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}
	data := objectPayload(out)
	config.IsAutoUpdateEnabled = types.BoolValue(boolFromAny(data["is_auto_update_enabled"]))
	config.IsRegistrationEnabled = types.BoolValue(boolFromAny(data["is_registration_enabled"]))
	config.IsUsageTrackingEnabled = types.BoolValue(boolFromAny(data["is_usage_tracking_enabled"]))
	config.FQDN = types.StringValue(stringFromAny(data["fqdn"]))
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
