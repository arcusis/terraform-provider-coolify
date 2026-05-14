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

// coolify_resource_action covers the start/stop/restart endpoints for
// applications, services, and databases.
// Create → POST action. Delete → stop. Update → new action if changed.

type resourceActionModel struct {
	ID           types.String `tfsdk:"id"`
	ResourceType types.String `tfsdk:"resource_type"`
	ResourceUUID types.String `tfsdk:"resource_uuid"`
	Action       types.String `tfsdk:"action"`
	Status       types.String `tfsdk:"status"`
}

type resourceActionResource struct{ client *coolify.Client }

func NewResourceActionResource() resource.Resource { return &resourceActionResource{} }

func (r *resourceActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_action"
}

func (r *resourceActionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true},
			"resource_type": schema.StringAttribute{Required: true, Description: "Type of resource: application, service, or database."},
			"resource_uuid": schema.StringAttribute{Required: true},
			"action":        schema.StringAttribute{Required: true, Description: "Action to perform: start, stop, or restart."},
			"status":        schema.StringAttribute{Computed: true},
		},
	}
}

func (r *resourceActionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *resourceActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.callAction(ctx, plan.ResourceType.ValueString(), plan.ResourceUUID.ValueString(), plan.Action.ValueString()); err != nil {
		resp.Diagnostics.AddError("Unable to perform Coolify resource action", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.ResourceType.ValueString() + "/" + plan.ResourceUUID.ValueString() + "/" + plan.Action.ValueString())
	plan.Status = types.StringValue("action_triggered")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	status := r.readStatus(ctx, state.ResourceType.ValueString(), state.ResourceUUID.ValueString())
	state.Status = types.StringValue(status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.callAction(ctx, plan.ResourceType.ValueString(), plan.ResourceUUID.ValueString(), plan.Action.ValueString()); err != nil {
		resp.Diagnostics.AddError("Unable to perform Coolify resource action", err.Error())
		return
	}

	plan.Status = types.StringValue("action_triggered")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// On destroy, stop the resource regardless of current action
	_ = r.callAction(ctx, state.ResourceType.ValueString(), state.ResourceUUID.ValueString(), "stop")
}

func (r *resourceActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *resourceActionResource) callAction(ctx context.Context, resourceType, resourceUUID, action string) error {
	actionPath := fmt.Sprintf("/api/v1/%ss/%s/%s", resourceType, resourceUUID, action)
	var out map[string]any
	err := r.client.Get(ctx, actionPath, &out)
	if err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return nil
		}
	}
	return err
}

func (r *resourceActionResource) readStatus(ctx context.Context, resourceType, resourceUUID string) string {
	var out map[string]any
	if err := r.client.Get(ctx, fmt.Sprintf("/api/v1/%ss/%s", resourceType, resourceUUID), &out); err != nil {
		return "unknown"
	}
	data := objectPayload(out)
	if v, ok := data["status"]; ok {
		return stringFromAny(v)
	}
	return "unknown"
}
