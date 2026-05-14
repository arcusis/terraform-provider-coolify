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

type deploymentDataSourceModel struct {
	UUID            types.String `tfsdk:"uuid"`
	ID              types.String `tfsdk:"id"`
	ApplicationUUID types.String `tfsdk:"application_uuid"`
	Status          types.String `tfsdk:"status"`
	CommitHash      types.String `tfsdk:"commit_hash"`
	CommitMessage   types.String `tfsdk:"commit_message"`
}

type deploymentDataSource struct{ client *coolify.Client }

func NewDeploymentDataSource() datasource.DataSource { return &deploymentDataSource{} }

func (d *deploymentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment"
}

func (d *deploymentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid":             schema.StringAttribute{Required: true},
			"id":               schema.StringAttribute{Computed: true},
			"application_uuid": schema.StringAttribute{Computed: true},
			"status":           schema.StringAttribute{Computed: true},
			"commit_hash":      schema.StringAttribute{Computed: true},
			"commit_message":   schema.StringAttribute{Computed: true},
		},
	}
}

func (d *deploymentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *deploymentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config deploymentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.UUID.ValueString() == "" {
		// UUID may be empty if the triggering deploy resource didn't get a deployment UUID back
		config.ID = types.StringValue("")
		config.Status = types.StringValue("unknown")
		config.ApplicationUUID = types.StringValue("")
		config.CommitHash = types.StringValue("")
		config.CommitMessage = types.StringValue("")
		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}

	var out map[string]any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/deployments/%s", config.UUID.ValueString()), &out); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Deployment not found", "No deployment with UUID "+config.UUID.ValueString())
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify deployment", err.Error())
		return
	}
	data := objectPayload(out)
	config.ID = types.StringValue(config.UUID.ValueString())
	for key, ptr := range map[string]*types.String{
		"application_uuid": &config.ApplicationUUID,
		"status":           &config.Status,
		"commit_hash":      &config.CommitHash,
		"commit_message":   &config.CommitMessage,
	} {
		if v, ok := data[key]; ok {
			*ptr = types.StringValue(stringFromAny(v))
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
