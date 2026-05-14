package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_application_preview manages PR preview deployments for an application.
// On delete it calls DELETE /applications/{uuid}/previews/{pull_request_id}.

type prPreviewModel struct {
	ID              types.String `tfsdk:"id"`
	ApplicationUUID types.String `tfsdk:"application_uuid"`
	PullRequestID   types.Int64  `tfsdk:"pull_request_id"`
}

type prPreviewResource struct{ client *coolify.Client }

func NewPRPreviewResource() resource.Resource { return &prPreviewResource{} }

func (r *prPreviewResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_preview"
}

func (r *prPreviewResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Computed: true},
			"application_uuid": schema.StringAttribute{Required: true},
			"pull_request_id":  schema.Int64Attribute{Required: true, Description: "GitHub pull request ID whose preview deployment to manage."},
		},
	}
}

func (r *prPreviewResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *prPreviewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan prPreviewModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s/%d", plan.ApplicationUUID.ValueString(), plan.PullRequestID.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *prPreviewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state prPreviewModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *prPreviewResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("PR previews cannot be updated", "Destroy and recreate to change pull_request_id.")
}

func (r *prPreviewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state prPreviewModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deletePath := fmt.Sprintf("/api/v1/applications/%s/previews/%d",
		state.ApplicationUUID.ValueString(), state.PullRequestID.ValueInt64())

	if err := r.client.Delete(ctx, deletePath, nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return // preview doesn't exist — already clean
		}
		resp.Diagnostics.AddError("Unable to delete Coolify PR preview", err.Error())
	}
}

func (r *prPreviewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
