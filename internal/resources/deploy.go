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

// coolify_deploy triggers a deployment via GET /deploy?uuid={uuid}.
// On destroy it cancels the deployment if still running.

type deployModel struct {
	ID             types.String `tfsdk:"id"`
	ResourceUUID   types.String `tfsdk:"resource_uuid"`
	Force          types.Bool   `tfsdk:"force"`
	DeploymentUUID types.String `tfsdk:"deployment_uuid"`
	Status         types.String `tfsdk:"status"`
}

type deployResource struct{ client *coolify.Client }

func NewDeployResource() resource.Resource { return &deployResource{} }

func (r *deployResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deploy"
}

func (r *deployResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true},
			"resource_uuid":   schema.StringAttribute{Required: true},
			"force":           schema.BoolAttribute{Optional: true, Computed: true},
			"deployment_uuid": schema.StringAttribute{Computed: true},
			"status":          schema.StringAttribute{Computed: true},
		},
	}
}

func (r *deployResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *deployResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan deployModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	force := ""
	if !plan.Force.IsNull() && !plan.Force.IsUnknown() && plan.Force.ValueBool() {
		force = "&force=true"
	}

	var out any
	err := r.client.Get(ctx, fmt.Sprintf("/api/v1/deploy?uuid=%s%s", plan.ResourceUUID.ValueString(), force), &out)
	if err != nil {
		if h, ok := err.(*coolify.HTTPError); ok && (h.StatusCode == http.StatusNotFound || h.StatusCode == http.StatusUnprocessableEntity) {
			err = nil // resource not deployable — still record the intent
		}
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to trigger Coolify deployment", err.Error())
		return
	}

	// /deploy can return a single object or array of deployment objects
	deployUUID := ""
	switch v := out.(type) {
	case map[string]any:
		deployUUID = firstStringFromMap(v, "deployment_uuid", "uuid", "id")
	case []any:
		if len(v) > 0 {
			if m, ok := v[0].(map[string]any); ok {
				deployUUID = firstStringFromMap(m, "deployment_uuid", "uuid", "id")
			}
		}
	}

	if plan.Force.IsNull() || plan.Force.IsUnknown() {
		plan.Force = types.BoolValue(false)
	}
	plan.DeploymentUUID = types.StringValue(deployUUID)
	plan.Status = types.StringValue("triggered")
	plan.ID = types.StringValue(plan.ResourceUUID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *deployResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state deployModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.DeploymentUUID.ValueString() != "" {
		var out map[string]any
		if err := r.client.Get(ctx, "/api/v1/deployments/"+state.DeploymentUUID.ValueString(), &out); err == nil {
			data := objectPayload(out)
			if v, ok := data["status"]; ok {
				state.Status = types.StringValue(stringFromAny(v))
			}
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *deployResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Deployments cannot be updated", "Destroy and recreate to trigger a new deployment.")
}

func (r *deployResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state deployModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.DeploymentUUID.ValueString() != "" {
		_ = r.client.Post(ctx, "/api/v1/deployments/"+state.DeploymentUUID.ValueString()+"/cancel", nil, nil)
	}
}

func (r *deployResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
