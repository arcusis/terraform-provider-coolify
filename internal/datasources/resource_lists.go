package datasources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ── Current team members ──────────────────────────────────────────────────────

type currentTeamMembersDataSourceModel struct {
	Members []teamMemberModel `tfsdk:"members"`
}

type currentTeamMembersDataSource struct{ client *coolify.Client }

func NewCurrentTeamMembersDataSource() datasource.DataSource {
	return &currentTeamMembersDataSource{}
}

func (d *currentTeamMembersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_current_team_members"
}

func (d *currentTeamMembersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists members of the team associated with the current API token via `GET /api/v1/teams/current/members`. Equivalent to `coolify_team_members` with the current team ID, but more convenient.",
		Attributes: map[string]schema.Attribute{
			"members": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Members of the current team.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":    schema.StringAttribute{Computed: true, Description: "User ID."},
						"name":  schema.StringAttribute{Computed: true, Description: "User display name."},
						"email": schema.StringAttribute{Computed: true, Description: "User email address."},
						"role":  schema.StringAttribute{Computed: true, Description: "Team role (owner, member)."},
					},
				},
			},
		},
	}
}

func (d *currentTeamMembersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *currentTeamMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config currentTeamMembersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/teams/current/members", &out); err != nil {
		resp.Diagnostics.AddError("Unable to read current team members", err.Error())
		return
	}
	config.Members = []teamMemberModel{}
	for _, item := range dataList(out) {
		config.Members = append(config.Members, teamMemberModel{
			ID:    types.StringValue(firstString(item, "id", "uuid")),
			Name:  types.StringValue(stringFromAny(item["name"])),
			Email: types.StringValue(stringFromAny(item["email"])),
			Role:  types.StringValue(stringFromAny(item["role"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Project environments list ─────────────────────────────────────────────────

type projectEnvironmentModel struct {
	UUID        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type projectEnvironmentsDataSourceModel struct {
	ProjectUUID  types.String               `tfsdk:"project_uuid"`
	Environments []projectEnvironmentModel  `tfsdk:"environments"`
}

type projectEnvironmentsDataSource struct{ client *coolify.Client }

func NewProjectEnvironmentsDataSource() datasource.DataSource {
	return &projectEnvironmentsDataSource{}
}

func (d *projectEnvironmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_environments"
}

func (d *projectEnvironmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all environments in a Coolify project via `GET /api/v1/projects/{uuid}/environments`. Use this to discover environment names before creating resources inside them.",
		Attributes: map[string]schema.Attribute{
			"project_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the project.",
			},
			"environments": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All environments in the project.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":        schema.StringAttribute{Computed: true, Description: "Environment UUID."},
						"name":        schema.StringAttribute{Computed: true, Description: "Environment name."},
						"description": schema.StringAttribute{Computed: true, Description: "Environment description."},
					},
				},
			},
		},
	}
}

func (d *projectEnvironmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *projectEnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config projectEnvironmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	path := fmt.Sprintf("/api/v1/projects/%s/environments", config.ProjectUUID.ValueString())
	if err := d.client.Get(ctx, path, &out); err != nil {
		resp.Diagnostics.AddError("Unable to list project environments", err.Error())
		return
	}
	config.Environments = []projectEnvironmentModel{}
	for _, item := range dataList(out) {
		config.Environments = append(config.Environments, projectEnvironmentModel{
			UUID:        types.StringValue(firstString(item, "uuid", "id")),
			Name:        types.StringValue(stringFromAny(item["name"])),
			Description: types.StringValue(stringFromAny(item["description"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Application scheduled tasks list ─────────────────────────────────────────

type scheduledTaskListItemModel struct {
	UUID      types.String `tfsdk:"uuid"`
	Name      types.String `tfsdk:"name"`
	Command   types.String `tfsdk:"command"`
	Frequency types.String `tfsdk:"frequency"`
	Enabled   types.Bool   `tfsdk:"enabled"`
}

type applicationScheduledTasksListDataSourceModel struct {
	ApplicationUUID types.String                 `tfsdk:"application_uuid"`
	Tasks           []scheduledTaskListItemModel  `tfsdk:"tasks"`
}

type applicationScheduledTasksListDataSource struct{ client *coolify.Client }

func NewApplicationScheduledTasksListDataSource() datasource.DataSource {
	return &applicationScheduledTasksListDataSource{}
}

func (d *applicationScheduledTasksListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_scheduled_tasks"
}

func (d *applicationScheduledTasksListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all scheduled tasks for a Coolify application via `GET /api/v1/applications/{uuid}/scheduled-tasks`.",
		Attributes: map[string]schema.Attribute{
			"application_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the application.",
			},
			"tasks": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All scheduled tasks for this application.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: scheduledTaskListItemAttrs(),
				},
			},
		},
	}
}

func (d *applicationScheduledTasksListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *applicationScheduledTasksListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config applicationScheduledTasksListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	path := fmt.Sprintf("/api/v1/applications/%s/scheduled-tasks", config.ApplicationUUID.ValueString())
	config.Tasks = readScheduledTaskList(ctx, d.client, path, resp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Service scheduled tasks list ──────────────────────────────────────────────

type serviceScheduledTasksListDataSourceModel struct {
	ServiceUUID types.String                 `tfsdk:"service_uuid"`
	Tasks       []scheduledTaskListItemModel  `tfsdk:"tasks"`
}

type serviceScheduledTasksListDataSource struct{ client *coolify.Client }

func NewServiceScheduledTasksListDataSource() datasource.DataSource {
	return &serviceScheduledTasksListDataSource{}
}

func (d *serviceScheduledTasksListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_scheduled_tasks"
}

func (d *serviceScheduledTasksListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all scheduled tasks for a Coolify service via `GET /api/v1/services/{uuid}/scheduled-tasks`.",
		Attributes: map[string]schema.Attribute{
			"service_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the service.",
			},
			"tasks": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All scheduled tasks for this service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: scheduledTaskListItemAttrs(),
				},
			},
		},
	}
}

func (d *serviceScheduledTasksListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *serviceScheduledTasksListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config serviceScheduledTasksListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	path := fmt.Sprintf("/api/v1/services/%s/scheduled-tasks", config.ServiceUUID.ValueString())
	config.Tasks = readScheduledTaskList(ctx, d.client, path, resp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func scheduledTaskListItemAttrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"uuid":      schema.StringAttribute{Computed: true, Description: "Task UUID."},
		"name":      schema.StringAttribute{Computed: true, Description: "Task name."},
		"command":   schema.StringAttribute{Computed: true, Description: "Shell command."},
		"frequency": schema.StringAttribute{Computed: true, Description: "Cron expression."},
		"enabled":   schema.BoolAttribute{Computed: true, Description: "Whether the task is active."},
	}
}

func readScheduledTaskList(ctx context.Context, client *coolify.Client, path string, resp *datasource.ReadResponse) []scheduledTaskListItemModel {
	var out any
	if err := client.Get(ctx, path, &out); err != nil {
		return []scheduledTaskListItemModel{}
	}
	tasks := []scheduledTaskListItemModel{}
	for _, item := range dataList(out) {
		tasks = append(tasks, scheduledTaskListItemModel{
			UUID:      types.StringValue(firstString(item, "uuid", "id")),
			Name:      types.StringValue(stringFromAny(item["name"])),
			Command:   types.StringValue(stringFromAny(item["command"])),
			Frequency: types.StringValue(stringFromAny(item["frequency"])),
			Enabled:   types.BoolValue(boolFromAny(item["enabled"])),
		})
	}
	return tasks
}

// ── Database backups list ─────────────────────────────────────────────────────

type databaseBackupListItemModel struct {
	UUID      types.String `tfsdk:"uuid"`
	Frequency types.String `tfsdk:"frequency"`
	Enabled   types.Bool   `tfsdk:"enabled"`
}

type databaseBackupsListDataSourceModel struct {
	DatabaseUUID types.String                  `tfsdk:"database_uuid"`
	Backups      []databaseBackupListItemModel  `tfsdk:"backups"`
}

type databaseBackupsListDataSource struct{ client *coolify.Client }

func NewDatabaseBackupsListDataSource() datasource.DataSource {
	return &databaseBackupsListDataSource{}
}

func (d *databaseBackupsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_backups"
}

func (d *databaseBackupsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all scheduled backup configurations for a database via `GET /api/v1/databases/{uuid}/backups`.",
		Attributes: map[string]schema.Attribute{
			"database_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the database.",
			},
			"backups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All backup schedules configured for this database.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":      schema.StringAttribute{Computed: true, Description: "Backup schedule UUID."},
						"frequency": schema.StringAttribute{Computed: true, Description: "Backup frequency (daily, weekly, etc.)."},
						"enabled":   schema.BoolAttribute{Computed: true, Description: "Whether the backup schedule is active."},
					},
				},
			},
		},
	}
}

func (d *databaseBackupsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *databaseBackupsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config databaseBackupsListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	path := fmt.Sprintf("/api/v1/databases/%s/backups", config.DatabaseUUID.ValueString())
	if err := d.client.Get(ctx, path, &out); err != nil {
		config.Backups = []databaseBackupListItemModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}
	config.Backups = []databaseBackupListItemModel{}
	for _, item := range dataList(out) {
		config.Backups = append(config.Backups, databaseBackupListItemModel{
			UUID:      types.StringValue(firstString(item, "uuid", "id")),
			Frequency: types.StringValue(stringFromAny(item["frequency"])),
			Enabled:   types.BoolValue(boolFromAny(item["enabled"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
