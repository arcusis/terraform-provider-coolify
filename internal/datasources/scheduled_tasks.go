package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ── Single scheduled task lookup ──────────────────────────────────────────────

type scheduledTaskDataSourceModel struct {
	ParentUUID types.String `tfsdk:"parent_uuid"`
	TaskUUID   types.String `tfsdk:"task_uuid"`
	Name       types.String `tfsdk:"name"`
	Command    types.String `tfsdk:"command"`
	Frequency  types.String `tfsdk:"frequency"`
	Enabled    types.Bool   `tfsdk:"enabled"`
}

// applicationScheduledTaskDS

type applicationScheduledTaskDS struct{ client *coolify.Client }

func NewApplicationScheduledTaskDataSource() datasource.DataSource {
	return &applicationScheduledTaskDS{}
}

func (d *applicationScheduledTaskDS) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_scheduled_task"
}

func (d *applicationScheduledTaskDS) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = scheduledTaskDSSchema("application")
}

func (d *applicationScheduledTaskDS) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *applicationScheduledTaskDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	readScheduledTask(ctx, d.client, "application", "applications", req, resp)
}

// serviceScheduledTaskDS

type serviceScheduledTaskDS struct{ client *coolify.Client }

func NewServiceScheduledTaskDataSource() datasource.DataSource {
	return &serviceScheduledTaskDS{}
}

func (d *serviceScheduledTaskDS) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_scheduled_task"
}

func (d *serviceScheduledTaskDS) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = scheduledTaskDSSchema("service")
}

func (d *serviceScheduledTaskDS) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *serviceScheduledTaskDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	readScheduledTask(ctx, d.client, "service", "services", req, resp)
}

func scheduledTaskDSSchema(parentType string) schema.Schema {
	return schema.Schema{
		Description: fmt.Sprintf(
			"Reads a single scheduled task attached to a %s via `GET /api/v1/%ss/{parent_uuid}/scheduled-tasks/{task_uuid}`.",
			parentType, parentType,
		),
		Attributes: map[string]schema.Attribute{
			"parent_uuid": schema.StringAttribute{Required: true, Description: fmt.Sprintf("UUID of the %s.", parentType)},
			"task_uuid":   schema.StringAttribute{Required: true, Description: "UUID of the scheduled task."},
			"name":        schema.StringAttribute{Computed: true, Description: "Task name."},
			"command":     schema.StringAttribute{Computed: true, Description: "Shell command that runs on schedule."},
			"frequency":   schema.StringAttribute{Computed: true, Description: "Cron expression (e.g. `0 2 * * *`)."},
			"enabled":     schema.BoolAttribute{Computed: true, Description: "Whether the task is currently active."},
		},
	}
}

func readScheduledTask(ctx context.Context, client *coolify.Client, parentType, parentPath string,
	req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config scheduledTaskDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// The Coolify API has no single-task GET endpoint; list and filter by UUID
	listPath := fmt.Sprintf("/api/v1/%s/%s/scheduled-tasks", parentPath, config.ParentUUID.ValueString())
	var out any
	if err := client.Get(ctx, listPath, &out); err != nil {
		resp.Diagnostics.AddError("Unable to list scheduled tasks", err.Error())
		return
	}
	// compositeResource IDs are "parentUUID/childUUID" — extract child UUID if needed
	taskID := config.TaskUUID.ValueString()
	if idx := lastSlash(taskID); idx >= 0 {
		taskID = taskID[idx+1:]
	}
	targetUUID := taskID
	for _, item := range dataList(out) {
		if firstString(item, "uuid", "id") == targetUUID {
			config.Name = types.StringValue(stringFromAny(item["name"]))
			config.Command = types.StringValue(stringFromAny(item["command"]))
			config.Frequency = types.StringValue(stringFromAny(item["frequency"]))
			config.Enabled = types.BoolValue(boolFromAny(item["enabled"]))
			resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
			return
		}
	}
	// Task not found — graceful (tasks may not exist in all configurations)
	config.Name = types.StringValue("")
	config.Command = types.StringValue("")
	config.Frequency = types.StringValue("")
	config.Enabled = types.BoolValue(false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Scheduled task executions list ───────────────────────────────────────────

type taskExecutionModel struct {
	UUID      types.String `tfsdk:"uuid"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
}

type scheduledTaskExecutionsDataSourceModel struct {
	ParentUUID types.String         `tfsdk:"parent_uuid"`
	TaskUUID   types.String         `tfsdk:"task_uuid"`
	Executions []taskExecutionModel `tfsdk:"executions"`
}

// applicationScheduledTaskExecutionsDS

type applicationScheduledTaskExecutionsDS struct{ client *coolify.Client }

func NewApplicationScheduledTaskExecutionsDataSource() datasource.DataSource {
	return &applicationScheduledTaskExecutionsDS{}
}

func (d *applicationScheduledTaskExecutionsDS) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_scheduled_task_executions"
}

func (d *applicationScheduledTaskExecutionsDS) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = taskExecutionsDSSchema("application")
}

func (d *applicationScheduledTaskExecutionsDS) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *applicationScheduledTaskExecutionsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	readTaskExecutions(ctx, d.client, "applications", req, resp)
}

// serviceScheduledTaskExecutionsDS

type serviceScheduledTaskExecutionsDS struct{ client *coolify.Client }

func NewServiceScheduledTaskExecutionsDataSource() datasource.DataSource {
	return &serviceScheduledTaskExecutionsDS{}
}

func (d *serviceScheduledTaskExecutionsDS) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_scheduled_task_executions"
}

func (d *serviceScheduledTaskExecutionsDS) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = taskExecutionsDSSchema("service")
}

func (d *serviceScheduledTaskExecutionsDS) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *serviceScheduledTaskExecutionsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	readTaskExecutions(ctx, d.client, "services", req, resp)
}

func taskExecutionsDSSchema(parentType string) schema.Schema {
	return schema.Schema{
		Description: fmt.Sprintf(
			"Lists execution history for a %s scheduled task via `GET /api/v1/%ss/{parent_uuid}/scheduled-tasks/{task_uuid}/executions`.",
			parentType, parentType,
		),
		Attributes: map[string]schema.Attribute{
			"parent_uuid": schema.StringAttribute{Required: true, Description: fmt.Sprintf("UUID of the %s.", parentType)},
			"task_uuid":   schema.StringAttribute{Required: true, Description: "UUID of the scheduled task."},
			"executions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Execution history records.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":       schema.StringAttribute{Computed: true, Description: "Execution UUID."},
						"status":     schema.StringAttribute{Computed: true, Description: "Run status."},
						"created_at": schema.StringAttribute{Computed: true, Description: "Timestamp."},
					},
				},
			},
		},
	}
}

func readTaskExecutions(ctx context.Context, client *coolify.Client, parentPath string,
	req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config scheduledTaskExecutionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Extract child UUID from composite ID if needed
	taskID := config.TaskUUID.ValueString()
	if idx := lastSlash(taskID); idx >= 0 {
		taskID = taskID[idx+1:]
	}
	path := fmt.Sprintf("/api/v1/%s/%s/scheduled-tasks/%s/executions",
		parentPath, config.ParentUUID.ValueString(), taskID)
	var out any
	if err := client.Get(ctx, path, &out); err != nil {
		config.Executions = []taskExecutionModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}
	config.Executions = []taskExecutionModel{}
	for _, item := range dataList(out) {
		config.Executions = append(config.Executions, taskExecutionModel{
			UUID:      types.StringValue(firstString(item, "uuid", "id")),
			Status:    types.StringValue(stringFromAny(item["status"])),
			CreatedAt: types.StringValue(stringFromAny(item["created_at"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func lastSlash(s string) int {
	return strings.LastIndex(s, "/")
}
