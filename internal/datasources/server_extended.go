package datasources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ── Server resources ─────────────────────────────────────────────────────────

type serverResourceItemModel struct {
	UUID        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Status      types.String `tfsdk:"status"`
	ProjectUUID types.String `tfsdk:"project_uuid"`
}

type serverResourcesDataSourceModel struct {
	ServerUUID types.String              `tfsdk:"server_uuid"`
	Resources  []serverResourceItemModel `tfsdk:"resources"`
}

type serverResourcesDataSource struct{ client *coolify.Client }

func NewServerResourcesDataSource() datasource.DataSource { return &serverResourcesDataSource{} }

func (d *serverResourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_resources"
}

func (d *serverResourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_uuid": schema.StringAttribute{Required: true},
			"resources": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":         schema.StringAttribute{Computed: true},
						"name":         schema.StringAttribute{Computed: true},
						"type":         schema.StringAttribute{Computed: true},
						"status":       schema.StringAttribute{Computed: true},
						"project_uuid": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *serverResourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *serverResourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config serverResourcesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/servers/%s/resources", config.ServerUUID.ValueString()), &out); err != nil {
		resp.Diagnostics.AddError("Unable to read Coolify server resources", err.Error())
		return
	}

	config.Resources = []serverResourceItemModel{}
	for _, item := range dataList(out) {
		config.Resources = append(config.Resources, serverResourceItemModel{
			UUID:        types.StringValue(firstString(item, "uuid", "id")),
			Name:        types.StringValue(stringFromAny(item["name"])),
			Type:        types.StringValue(stringFromAny(item["type"])),
			Status:      types.StringValue(stringFromAny(item["status"])),
			ProjectUUID: types.StringValue(stringFromAny(item["project_uuid"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Server domains ────────────────────────────────────────────────────────────

type serverDomainModel struct {
	Domain types.String `tfsdk:"domain"`
}

type serverDomainsDataSourceModel struct {
	ServerUUID types.String        `tfsdk:"server_uuid"`
	Domains    []serverDomainModel `tfsdk:"domains"`
}

type serverDomainsDataSource struct{ client *coolify.Client }

func NewServerDomainsDataSource() datasource.DataSource { return &serverDomainsDataSource{} }

func (d *serverDomainsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_domains"
}

func (d *serverDomainsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_uuid": schema.StringAttribute{Required: true},
			"domains": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *serverDomainsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *serverDomainsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config serverDomainsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/servers/%s/domains", config.ServerUUID.ValueString()), &out); err != nil {
		resp.Diagnostics.AddError("Unable to read Coolify server domains", err.Error())
		return
	}

	config.Domains = []serverDomainModel{}
	for _, item := range dataList(out) {
		if domain := stringFromAny(item["domain"]); domain != "" {
			config.Domains = append(config.Domains, serverDomainModel{
				Domain: types.StringValue(domain),
			})
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
