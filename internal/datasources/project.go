package datasources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type projectDataSource struct {
	client *coolify.Client
}

type projectDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	UUID        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Terraform data source ID.",
			},
			"uuid": schema.StringAttribute{
				Required:    true,
				Description: "Coolify project UUID.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Coolify project name.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Coolify project description.",
			},
		},
	}
}

func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*coolify.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *coolify.Client from provider configuration.")
		return
	}
	d.client = client
}

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config projectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := config.UUID.ValueString()
	var out map[string]any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/projects/%s", uuid), &out); err != nil {
		resp.Diagnostics.AddError("Unable to read Coolify project data source", err.Error())
		return
	}
	data := objectPayload(out)
	if id := firstString(data, "uuid", "id"); id != "" {
		config.ID = types.StringValue(id)
		config.UUID = types.StringValue(id)
	} else {
		config.ID = types.StringValue(uuid)
	}
	if value, ok := data["name"]; ok {
		config.Name = types.StringValue(stringFromAny(value))
	}
	if value, ok := data["description"]; ok {
		config.Description = types.StringValue(stringFromAny(value))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
