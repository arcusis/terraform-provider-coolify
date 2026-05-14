package datasources

import (
	"context"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_resources lists all resources across all projects via GET /resources.

type coolifyResourceItemModel struct {
	UUID        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Status      types.String `tfsdk:"status"`
	ProjectUUID types.String `tfsdk:"project_uuid"`
}

type coolifyResourcesDataSourceModel struct {
	Resources []coolifyResourceItemModel `tfsdk:"resources"`
}

type coolifyResourcesDataSource struct{ client *coolify.Client }

func NewCoolifyResourcesDataSource() datasource.DataSource { return &coolifyResourcesDataSource{} }

func (d *coolifyResourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resources"
}

func (d *coolifyResourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
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

func (d *coolifyResourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *coolifyResourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config coolifyResourcesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	if err := d.client.Get(ctx, "/api/v1/resources", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify resources", err.Error())
		return
	}

	config.Resources = []coolifyResourceItemModel{}
	for _, item := range dataList(out) {
		config.Resources = append(config.Resources, coolifyResourceItemModel{
			UUID:        types.StringValue(firstString(item, "uuid", "id")),
			Name:        types.StringValue(stringFromAny(item["name"])),
			Type:        types.StringValue(stringFromAny(item["type"])),
			Status:      types.StringValue(stringFromAny(item["status"])),
			ProjectUUID: types.StringValue(stringFromAny(item["project_uuid"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
