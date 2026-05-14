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

type serviceDataSourceModel struct {
	UUID        types.String `tfsdk:"uuid"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Status      types.String `tfsdk:"status"`
	ServiceType types.String `tfsdk:"service_type"`
}

type serviceDataSource struct{ client *coolify.Client }

func NewServiceDataSource() datasource.DataSource { return &serviceDataSource{} }

func (d *serviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *serviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid":         schema.StringAttribute{Required: true},
			"id":           schema.StringAttribute{Computed: true},
			"name":         schema.StringAttribute{Computed: true},
			"description":  schema.StringAttribute{Computed: true},
			"status":       schema.StringAttribute{Computed: true},
			"service_type": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *serviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config serviceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out map[string]any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/services/%s", config.UUID.ValueString()), &out); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Service not found", "No service with UUID "+config.UUID.ValueString())
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify service", err.Error())
		return
	}
	data := objectPayload(out)
	config.ID = types.StringValue(config.UUID.ValueString())
	for key, ptr := range map[string]*types.String{
		"name":        &config.Name,
		"description": &config.Description,
		"status":      &config.Status,
		"type":        &config.ServiceType,
	} {
		if v, ok := data[key]; ok {
			*ptr = types.StringValue(stringFromAny(v))
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
