package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type databaseDataSourceModel struct {
	UUID        types.String `tfsdk:"uuid"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Status      types.String `tfsdk:"status"`
	DBType      types.String `tfsdk:"db_type"`
}

type databaseDataSource struct{ client *coolify.Client }

func NewDatabaseDataSource() datasource.DataSource { return &databaseDataSource{} }

func (d *databaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (d *databaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid":        schema.StringAttribute{Required: true},
			"id":          schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Computed: true},
			"description": schema.StringAttribute{Computed: true},
			"status":      schema.StringAttribute{Computed: true},
			"db_type":     schema.StringAttribute{Computed: true},
		},
	}
}

func (d *databaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*coolify.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *coolify.Client.")
		return
	}
	d.client = client
}

func (d *databaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config databaseDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out map[string]any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/databases/%s", config.UUID.ValueString()), &out); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Database not found", "No database with UUID "+config.UUID.ValueString())
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify database", err.Error())
		return
	}
	data := objectPayload(out)
	config.ID = types.StringValue(config.UUID.ValueString())
	for key, ptr := range map[string]*types.String{
		"name":        &config.Name,
		"description": &config.Description,
		"status":      &config.Status,
		"type":        &config.DBType,
	} {
		if v, ok := data[key]; ok {
			*ptr = types.StringValue(stringFromAny(v))
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
