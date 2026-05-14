package datasources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type deploymentItemModel struct {
	UUID            types.String `tfsdk:"uuid"`
	ApplicationUUID types.String `tfsdk:"application_uuid"`
	Status          types.String `tfsdk:"status"`
	CommitMessage   types.String `tfsdk:"commit_message"`
}

// ── All deployments ───────────────────────────────────────────────────────────

type deploymentsDataSourceModel struct {
	Deployments []deploymentItemModel `tfsdk:"deployments"`
}

type deploymentsDataSource struct{ client *coolify.Client }

func NewDeploymentsDataSource() datasource.DataSource { return &deploymentsDataSource{} }

func (d *deploymentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployments"
}

func (d *deploymentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"deployments": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: deploymentAttrs(),
				},
			},
		},
	}
}

func (d *deploymentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *deploymentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config deploymentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	if err := d.client.Get(ctx, "/api/v1/deployments", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify deployments", err.Error())
		return
	}
	config.Deployments = parseDeployments(out)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Application deployments ───────────────────────────────────────────────────

type applicationDeploymentsDataSourceModel struct {
	ApplicationUUID types.String          `tfsdk:"application_uuid"`
	Deployments     []deploymentItemModel `tfsdk:"deployments"`
}

type applicationDeploymentsDataSource struct{ client *coolify.Client }

func NewApplicationDeploymentsDataSource() datasource.DataSource {
	return &applicationDeploymentsDataSource{}
}

func (d *applicationDeploymentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_deployments"
}

func (d *applicationDeploymentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"application_uuid": schema.StringAttribute{Required: true},
			"deployments": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: deploymentAttrs(),
				},
			},
		},
	}
}

func (d *applicationDeploymentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *applicationDeploymentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config applicationDeploymentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/deployments/applications/%s", config.ApplicationUUID.ValueString()), &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify application deployments", err.Error())
		return
	}
	config.Deployments = parseDeployments(out)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func deploymentAttrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"uuid":             schema.StringAttribute{Computed: true},
		"application_uuid": schema.StringAttribute{Computed: true},
		"status":           schema.StringAttribute{Computed: true},
		"commit_message":   schema.StringAttribute{Computed: true},
	}
}

func parseDeployments(out any) []deploymentItemModel {
	items := []deploymentItemModel{}
	for _, item := range dataList(out) {
		items = append(items, deploymentItemModel{
			UUID:            types.StringValue(firstString(item, "deployment_uuid", "uuid", "id")),
			ApplicationUUID: types.StringValue(stringFromAny(item["application_uuid"])),
			Status:          types.StringValue(stringFromAny(item["status"])),
			CommitMessage:   types.StringValue(stringFromAny(item["commit_message"])),
		})
	}
	return items
}
