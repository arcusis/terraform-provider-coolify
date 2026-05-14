package datasources

import (
	"context"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_system_info reads the Coolify health and version endpoints.

type systemInfoDataSourceModel struct {
	Healthy types.Bool   `tfsdk:"healthy"`
	Version types.String `tfsdk:"version"`
}

type systemInfoDataSource struct{ client *coolify.Client }

func NewSystemInfoDataSource() datasource.DataSource { return &systemInfoDataSource{} }

func (d *systemInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_info"
}

func (d *systemInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"healthy": schema.BoolAttribute{Computed: true},
			"version": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *systemInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *systemInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config systemInfoDataSourceModel

	// GET /health
	var healthOut map[string]any
	config.Healthy = types.BoolValue(false)
	if err := d.client.Get(ctx, "/api/v1/health", &healthOut); err == nil {
		healthy := healthOut["status"]
		config.Healthy = types.BoolValue(stringFromAny(healthy) == "ok" || healthy == true)
	}

	// GET /version
	var versionOut map[string]any
	config.Version = types.StringValue("")
	if err := d.client.Get(ctx, "/api/v1/version", &versionOut); err == nil {
		config.Version = types.StringValue(firstString(versionOut, "version"))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
