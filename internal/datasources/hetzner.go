package datasources

import (
	"context"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ── Hetzner Locations ────────────────────────────────────────────────────────

type hetznerLocationModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	City        types.String `tfsdk:"city"`
	Country     types.String `tfsdk:"country"`
}

type hetznerLocationsDataSourceModel struct {
	CloudTokenUUID types.String           `tfsdk:"cloud_token_uuid"`
	Locations      []hetznerLocationModel `tfsdk:"locations"`
}

type hetznerLocationsDataSource struct{ client *coolify.Client }

func NewHetznerLocationsDataSource() datasource.DataSource { return &hetznerLocationsDataSource{} }

func (d *hetznerLocationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hetzner_locations"
}

func (d *hetznerLocationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cloud_token_uuid": schema.StringAttribute{Required: true},
			"locations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":        schema.StringAttribute{Computed: true},
						"description": schema.StringAttribute{Computed: true},
						"city":        schema.StringAttribute{Computed: true},
						"country":     schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *hetznerLocationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *hetznerLocationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config hetznerLocationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	if err := d.client.Get(ctx, "/api/v1/hetzner/locations?cloud_token_uuid="+config.CloudTokenUUID.ValueString(), &out); err != nil {
		resp.Diagnostics.AddError("Unable to read Hetzner locations", err.Error())
		return
	}
	config.Locations = []hetznerLocationModel{}
	for _, item := range dataList(out) {
		config.Locations = append(config.Locations, hetznerLocationModel{
			Name:        types.StringValue(stringFromAny(item["name"])),
			Description: types.StringValue(stringFromAny(item["description"])),
			City:        types.StringValue(stringFromAny(item["city"])),
			Country:     types.StringValue(stringFromAny(item["country"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Hetzner Server Types ──────────────────────────────────────────────────────

type hetznerServerTypeModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Cores       types.Int64  `tfsdk:"cores"`
	Memory      types.Int64  `tfsdk:"memory"`
}

type hetznerServerTypesDataSourceModel struct {
	CloudTokenUUID types.String             `tfsdk:"cloud_token_uuid"`
	ServerTypes    []hetznerServerTypeModel `tfsdk:"server_types"`
}

type hetznerServerTypesDataSource struct{ client *coolify.Client }

func NewHetznerServerTypesDataSource() datasource.DataSource {
	return &hetznerServerTypesDataSource{}
}

func (d *hetznerServerTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hetzner_server_types"
}

func (d *hetznerServerTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cloud_token_uuid": schema.StringAttribute{Required: true},
			"server_types": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":        schema.StringAttribute{Computed: true},
						"description": schema.StringAttribute{Computed: true},
						"cores":       schema.Int64Attribute{Computed: true},
						"memory":      schema.Int64Attribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *hetznerServerTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *hetznerServerTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config hetznerServerTypesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	if err := d.client.Get(ctx, "/api/v1/hetzner/server-types?cloud_token_uuid="+config.CloudTokenUUID.ValueString(), &out); err != nil {
		resp.Diagnostics.AddError("Unable to read Hetzner server types", err.Error())
		return
	}
	config.ServerTypes = []hetznerServerTypeModel{}
	for _, item := range dataList(out) {
		config.ServerTypes = append(config.ServerTypes, hetznerServerTypeModel{
			Name:        types.StringValue(stringFromAny(item["name"])),
			Description: types.StringValue(stringFromAny(item["description"])),
			Cores:       types.Int64Value(int64FromAny(item["cores"])),
			Memory:      types.Int64Value(int64FromAny(item["memory"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
