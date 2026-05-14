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

type privateKeyDataSourceModel struct {
	UUID        types.String `tfsdk:"uuid"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type privateKeyDataSource struct{ client *coolify.Client }

func NewPrivateKeyDataSource() datasource.DataSource { return &privateKeyDataSource{} }

func (d *privateKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_key"
}

func (d *privateKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid":        schema.StringAttribute{Required: true},
			"id":          schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Computed: true},
			"description": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *privateKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *privateKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config privateKeyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out map[string]any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/security/keys/%s", config.UUID.ValueString()), &out); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Private key not found", "No private key with UUID "+config.UUID.ValueString())
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify private key", err.Error())
		return
	}
	data := objectPayload(out)
	config.ID = types.StringValue(config.UUID.ValueString())
	if v, ok := data["name"]; ok {
		config.Name = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["description"]; ok {
		config.Description = types.StringValue(stringFromAny(v))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
