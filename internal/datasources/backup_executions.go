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

type backupExecutionModel struct {
	UUID      types.String `tfsdk:"uuid"`
	Status    types.String `tfsdk:"status"`
	StartedAt types.String `tfsdk:"started_at"`
	EndedAt   types.String `tfsdk:"ended_at"`
}

type backupExecutionsDataSourceModel struct {
	DatabaseUUID         types.String           `tfsdk:"database_uuid"`
	ScheduledBackupUUID  types.String           `tfsdk:"scheduled_backup_uuid"`
	Executions           []backupExecutionModel `tfsdk:"executions"`
}

type backupExecutionsDataSource struct{ client *coolify.Client }

func NewBackupExecutionsDataSource() datasource.DataSource { return &backupExecutionsDataSource{} }

func (d *backupExecutionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_executions"
}

func (d *backupExecutionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists execution history for a scheduled database backup. Each execution record shows whether the backup succeeded and when it ran.",
		Attributes: map[string]schema.Attribute{
			"database_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the database that owns the backup schedule.",
			},
			"scheduled_backup_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the scheduled backup configuration (from `coolify_database_backup`).",
			},
			"executions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Ordered list of backup run records (newest first).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":       schema.StringAttribute{Computed: true, Description: "Execution UUID."},
						"status":     schema.StringAttribute{Computed: true, Description: "Run status (success, failed, running)."},
						"started_at": schema.StringAttribute{Computed: true, Description: "ISO 8601 start timestamp."},
						"ended_at":   schema.StringAttribute{Computed: true, Description: "ISO 8601 end timestamp."},
					},
				},
			},
		},
	}
}

func (d *backupExecutionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *backupExecutionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config backupExecutionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// compositeResource IDs are "parentUUID/childUUID" — extract child UUID if needed
	backupID := config.ScheduledBackupUUID.ValueString()
	if idx := strings.LastIndex(backupID, "/"); idx >= 0 {
		backupID = backupID[idx+1:]
	}

	path := fmt.Sprintf("/api/v1/databases/%s/backups/%s/executions",
		config.DatabaseUUID.ValueString(),
		backupID,
	)
	var out any
	if err := d.client.Get(ctx, path, &out); err != nil {
		// Graceful: no executions yet is not an error
		config.Executions = []backupExecutionModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}

	config.Executions = []backupExecutionModel{}
	for _, item := range dataList(out) {
		config.Executions = append(config.Executions, backupExecutionModel{
			UUID:      types.StringValue(firstString(item, "uuid", "id")),
			Status:    types.StringValue(stringFromAny(item["status"])),
			StartedAt: types.StringValue(stringFromAny(item["created_at"])),
			EndedAt:   types.StringValue(stringFromAny(item["updated_at"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
