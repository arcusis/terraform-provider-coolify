package resources

import (
	"context"
	"fmt"
	"net/http"
	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_scheduled_task_execution manages a single execution record for a
// scheduled task attached to an application or service.
// The resource is a deletion gate: creating it records the execution UUID in
// state; destroying it calls DELETE to purge the execution log.
//
// Endpoint: DELETE /api/v1/{parentType}s/{parentUUID}/scheduled-tasks/{taskUUID}/executions/{executionUUID}

type scheduledTaskExecutionModel struct {
	ParentUUID    types.String `tfsdk:"parent_uuid"`
	TaskUUID      types.String `tfsdk:"task_uuid"`
	ExecutionUUID types.String `tfsdk:"execution_uuid"`
}

type scheduledTaskExecutionResource struct {
	client     *coolify.Client
	parentType string // "application" or "service"
}

func NewApplicationScheduledTaskExecutionResource() resource.Resource {
	return &scheduledTaskExecutionResource{parentType: "application"}
}

func NewServiceScheduledTaskExecutionResource() resource.Resource {
	return &scheduledTaskExecutionResource{parentType: "service"}
}

func (r *scheduledTaskExecutionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.parentType + "_scheduled_task_execution"
}

func (r *scheduledTaskExecutionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: fmt.Sprintf(
			"Tracks a single execution record for a %s scheduled task. "+
				"Destroying this resource deletes the execution log entry via "+
				"`DELETE /api/v1/%ss/{parent_uuid}/scheduled-tasks/{task_uuid}/executions/{execution_uuid}`. "+
				"A 404 on delete is treated as success (record already gone).",
			r.parentType, r.parentType,
		),
		Attributes: map[string]schema.Attribute{
			"parent_uuid": schema.StringAttribute{
				Required:    true,
				Description: fmt.Sprintf("UUID of the %s that owns the scheduled task.", r.parentType),
			},
			"task_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the scheduled task.",
			},
			"execution_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the specific execution record to manage.",
			},
		},
	}
}

func (r *scheduledTaskExecutionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *scheduledTaskExecutionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan scheduledTaskExecutionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Nothing to create on the API — we're just tracking an existing execution.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scheduledTaskExecutionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state scheduledTaskExecutionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// No single-item GET for executions; keep state as-is.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *scheduledTaskExecutionResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *scheduledTaskExecutionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state scheduledTaskExecutionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := fmt.Sprintf("/api/v1/%ss/%s/scheduled-tasks/%s/executions/%s",
		r.parentType,
		state.ParentUUID.ValueString(),
		state.TaskUUID.ValueString(),
		state.ExecutionUUID.ValueString(),
	)
	if err := r.client.Delete(ctx, path, nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete scheduled task execution", err.Error())
	}
}
