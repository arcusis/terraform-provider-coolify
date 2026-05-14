package datasources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type applicationLogsDataSourceModel struct {
	ApplicationUUID types.String `tfsdk:"application_uuid"`
	Logs            types.String `tfsdk:"logs"`
}

type applicationLogsDataSource struct{ client *coolify.Client }

func NewApplicationLogsDataSource() datasource.DataSource { return &applicationLogsDataSource{} }

func (d *applicationLogsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_logs"
}

func (d *applicationLogsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches recent log output for a Coolify application. Useful for checking startup errors or last-run output without leaving Terraform.",
		Attributes: map[string]schema.Attribute{
			"application_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the application to fetch logs for.",
			},
			"logs": schema.StringAttribute{
				Computed:    true,
				Description: "Raw log output from the application container.",
			},
		},
	}
}

func (d *applicationLogsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *applicationLogsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config applicationLogsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	path := fmt.Sprintf("/api/v1/applications/%s/logs", config.ApplicationUUID.ValueString())
	if err := d.client.Get(ctx, path, &out); err != nil {
		// Logs endpoint may return 404 when app has never started; treat as empty
		config.Logs = types.StringValue("")
		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}

	switch v := out.(type) {
	case map[string]any:
		if logs, ok := v["logs"]; ok {
			config.Logs = types.StringValue(stringFromAny(logs))
		} else {
			config.Logs = types.StringValue(fmt.Sprint(v))
		}
	case string:
		config.Logs = types.StringValue(v)
	default:
		config.Logs = types.StringValue(fmt.Sprint(v))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
