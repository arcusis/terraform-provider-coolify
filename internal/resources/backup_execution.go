package resources

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_backup_execution tracks a single backup execution record.
// Destroying it calls DELETE /databases/{uuid}/backups/{sbUUID}/executions/{execUUID}.

type backupExecutionResourceModel struct {
	DatabaseUUID        types.String `tfsdk:"database_uuid"`
	ScheduledBackupUUID types.String `tfsdk:"scheduled_backup_uuid"`
	ExecutionUUID       types.String `tfsdk:"execution_uuid"`
}

type backupExecutionResource struct{ client *coolify.Client }

func NewBackupExecutionResource() resource.Resource { return &backupExecutionResource{} }

func (r *backupExecutionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_execution"
}

func (r *backupExecutionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Tracks a database backup execution record. Destroying this resource purges the execution log entry via `DELETE /api/v1/databases/{database_uuid}/backups/{scheduled_backup_uuid}/executions/{execution_uuid}`. A 404 on delete is treated as success.",
		Attributes: map[string]schema.Attribute{
			"database_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the database that owns the backup schedule.",
			},
			"scheduled_backup_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the scheduled backup configuration.",
			},
			"execution_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the specific execution record to track.",
			},
		},
	}
}

func (r *backupExecutionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *backupExecutionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan backupExecutionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupExecutionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state backupExecutionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupExecutionResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *backupExecutionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state backupExecutionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Extract child UUID from composite ID if needed
	backupID := state.ScheduledBackupUUID.ValueString()
	if idx := strings.LastIndex(backupID, "/"); idx >= 0 {
		backupID = backupID[idx+1:]
	}
	path := fmt.Sprintf("/api/v1/databases/%s/backups/%s/executions/%s",
		state.DatabaseUUID.ValueString(),
		backupID,
		state.ExecutionUUID.ValueString(),
	)
	if err := r.client.Delete(ctx, path, nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete backup execution", err.Error())
	}
}
