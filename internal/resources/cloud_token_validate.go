package resources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_cloud_token_validate triggers POST /cloud-tokens/{uuid}/validate.
// Creating the resource runs the validation; deleting removes it from state.

type cloudTokenValidateModel struct {
	CloudTokenUUID types.String `tfsdk:"cloud_token_uuid"`
	Status         types.String `tfsdk:"status"`
}

type cloudTokenValidateResource struct{ client *coolify.Client }

func NewCloudTokenValidateResource() resource.Resource { return &cloudTokenValidateResource{} }

func (r *cloudTokenValidateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_token_validate"
}

func (r *cloudTokenValidateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Validates a Coolify cloud token against the provider API (Hetzner, DigitalOcean, etc.) via `POST /api/v1/cloud-tokens/{uuid}/validate`. The validation result is stored in `status`. Use this as a precondition before provisioning cloud resources.",
		Attributes: map[string]schema.Attribute{
			"cloud_token_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the cloud token to validate (from `coolify_cloud_token`).",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Validation result returned by Coolify.",
			},
		},
	}
}

func (r *cloudTokenValidateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *cloudTokenValidateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cloudTokenValidateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	path := fmt.Sprintf("/api/v1/cloud-tokens/%s/validate", plan.CloudTokenUUID.ValueString())
	if err := r.client.Post(ctx, path, nil, &out); err != nil {
		// Non-fatal: store the error as status
		plan.Status = types.StringValue("invalid: " + err.Error())
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}
	status := "valid"
	if m, ok := out.(map[string]any); ok {
		if v := m["status"]; v != nil {
			status = fmt.Sprint(v)
		} else if v := m["message"]; v != nil {
			status = fmt.Sprint(v)
		}
	}
	plan.Status = types.StringValue(status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *cloudTokenValidateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cloudTokenValidateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *cloudTokenValidateResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *cloudTokenValidateResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
