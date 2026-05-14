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

type applicationDataSourceModel struct {
	UUID            types.String `tfsdk:"uuid"`
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	GitRepository   types.String `tfsdk:"git_repository"`
	GitBranch       types.String `tfsdk:"git_branch"`
	BuildPack       types.String `tfsdk:"build_pack"`
	Status          types.String `tfsdk:"status"`
	Domains         types.String `tfsdk:"domains"`
}

type applicationDataSource struct{ client *coolify.Client }

func NewApplicationDataSource() datasource.DataSource { return &applicationDataSource{} }

func (d *applicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *applicationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid":           schema.StringAttribute{Required: true},
			"id":             schema.StringAttribute{Computed: true},
			"name":           schema.StringAttribute{Computed: true},
			"description":    schema.StringAttribute{Computed: true},
			"git_repository": schema.StringAttribute{Computed: true},
			"git_branch":     schema.StringAttribute{Computed: true},
			"build_pack":     schema.StringAttribute{Computed: true},
			"status":         schema.StringAttribute{Computed: true},
			"domains":        schema.StringAttribute{Computed: true},
		},
	}
}

func (d *applicationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config applicationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out map[string]any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/applications/%s", config.UUID.ValueString()), &out); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Application not found", "No application with UUID "+config.UUID.ValueString())
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify application", err.Error())
		return
	}
	data := objectPayload(out)
	config.ID = types.StringValue(config.UUID.ValueString())
	for key, ptr := range map[string]*types.String{
		"name":           &config.Name,
		"description":    &config.Description,
		"git_repository": &config.GitRepository,
		"git_branch":     &config.GitBranch,
		"build_pack":     &config.BuildPack,
		"status":         &config.Status,
		"domains":        &config.Domains,
	} {
		if v, ok := data[key]; ok {
			*ptr = types.StringValue(stringFromAny(v))
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
