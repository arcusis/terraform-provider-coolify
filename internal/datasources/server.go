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

type serverDataSourceModel struct {
	UUID        types.String `tfsdk:"uuid"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IP          types.String `tfsdk:"ip"`
	Port        types.Int64  `tfsdk:"port"`
	User        types.String `tfsdk:"user"`
	ProxyType   types.String `tfsdk:"proxy_type"`
	IsReachable types.Bool   `tfsdk:"is_reachable"`
	IsUsable    types.Bool   `tfsdk:"is_usable"`
}

type serverDataSource struct{ client *coolify.Client }

func NewServerDataSource() datasource.DataSource { return &serverDataSource{} }

func (d *serverDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (d *serverDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid":         schema.StringAttribute{Required: true},
			"id":           schema.StringAttribute{Computed: true},
			"name":         schema.StringAttribute{Computed: true},
			"description":  schema.StringAttribute{Computed: true},
			"ip":           schema.StringAttribute{Computed: true},
			"port":         schema.Int64Attribute{Computed: true},
			"user":         schema.StringAttribute{Computed: true},
			"proxy_type":   schema.StringAttribute{Computed: true},
			"is_reachable": schema.BoolAttribute{Computed: true},
			"is_usable":    schema.BoolAttribute{Computed: true},
		},
	}
}

func (d *serverDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serverDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config serverDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out map[string]any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/servers/%s", config.UUID.ValueString()), &out); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Server not found", "No server with UUID "+config.UUID.ValueString())
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify server", err.Error())
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
	if v, ok := data["ip"]; ok {
		config.IP = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["port"]; ok {
		config.Port = types.Int64Value(int64FromAny(v))
	}
	if v, ok := data["user"]; ok {
		config.User = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["proxy_type"]; ok {
		config.ProxyType = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["is_reachable"]; ok {
		config.IsReachable = types.BoolValue(boolFromAny(v))
	}
	if v, ok := data["is_usable"]; ok {
		config.IsUsable = types.BoolValue(boolFromAny(v))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
