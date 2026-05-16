package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/arcusis/terraform-provider-coolify/internal/datasources"
	"github.com/arcusis/terraform-provider-coolify/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type CoolifyProvider struct {
	version string
}

type configModel struct {
	Endpoint             types.String `tfsdk:"endpoint"`
	Token                types.String `tfsdk:"token"`
	CFAccessClientID     types.String `tfsdk:"cf_access_client_id"`
	CFAccessClientSecret types.String `tfsdk:"cf_access_client_secret"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CoolifyProvider{version: version}
	}
}

func (p *CoolifyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "coolify"
	resp.Version = p.version
}

func (p *CoolifyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "Coolify API endpoint, for example https://coolify.example.com. Can also be set with COOLIFY_ENDPOINT.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Coolify API bearer token. Can also be set with COOLIFY_TOKEN.",
			},
			"cf_access_client_id": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Cloudflare Access service token client ID. Required when the Coolify endpoint is protected by Cloudflare Access. Can also be set with COOLIFY_CF_ACCESS_CLIENT_ID.",
			},
			"cf_access_client_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Cloudflare Access service token client secret. Required when the Coolify endpoint is protected by Cloudflare Access. Can also be set with COOLIFY_CF_ACCESS_CLIENT_SECRET.",
			},
		},
	}
}

func (p *CoolifyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config configModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := strings.TrimSpace(os.Getenv("COOLIFY_ENDPOINT"))
	token := strings.TrimSpace(os.Getenv("COOLIFY_TOKEN"))
	cfAccessClientID := strings.TrimSpace(os.Getenv("COOLIFY_CF_ACCESS_CLIENT_ID"))
	cfAccessClientSecret := strings.TrimSpace(os.Getenv("COOLIFY_CF_ACCESS_CLIENT_SECRET"))

	if !config.Endpoint.IsNull() && !config.Endpoint.IsUnknown() {
		endpoint = strings.TrimSpace(config.Endpoint.ValueString())
	}
	if !config.Token.IsNull() && !config.Token.IsUnknown() {
		token = strings.TrimSpace(config.Token.ValueString())
	}
	if !config.CFAccessClientID.IsNull() && !config.CFAccessClientID.IsUnknown() {
		cfAccessClientID = strings.TrimSpace(config.CFAccessClientID.ValueString())
	}
	if !config.CFAccessClientSecret.IsNull() && !config.CFAccessClientSecret.IsUnknown() {
		cfAccessClientSecret = strings.TrimSpace(config.CFAccessClientSecret.ValueString())
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Coolify endpoint",
			"Set the provider endpoint attribute or the COOLIFY_ENDPOINT environment variable.",
		)
	}
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Coolify token",
			"Set the provider token attribute or the COOLIFY_TOKEN environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client := coolify.NewClient(endpoint, token, cfAccessClientID, cfAccessClientSecret)
	tflog.Debug(ctx, "configured Coolify API client", map[string]any{"endpoint": endpoint})
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CoolifyProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Projects & environments
		resources.NewProjectResource,
		resources.NewEnvironmentResource,
		// Infrastructure
		resources.NewServerResource,
		resources.NewPrivateKeyResource,
		resources.NewCloudTokenResource,
		resources.NewGitHubAppResource,
		// Applications
		resources.NewApplicationResource,
		resources.NewApplicationStorageResource,
		resources.NewApplicationScheduledTaskResource,
		// Services
		resources.NewServiceResource,
		resources.NewServiceStorageResource,
		resources.NewServiceScheduledTaskResource,
		resources.NewServiceEnvironmentVariableResource,
		// Databases
		resources.NewPostgreSQLDatabaseResource,
		resources.NewMySQLDatabaseResource,
		resources.NewMariaDBDatabaseResource,
		resources.NewRedisDatabaseResource,
		resources.NewMongoDBDatabaseResource,
		resources.NewKeyDBDatabaseResource,
		resources.NewDragonflyDatabaseResource,
		resources.NewClickhouseDatabaseResource,
		resources.NewDatabaseBackupResource,
		resources.NewDatabaseStorageResource,
		resources.NewDatabaseEnvironmentVariableResource,
		// Environment variables
		resources.NewEnvironmentVariableResource,
		// Lifecycle & operations
		resources.NewResourceActionResource,
		resources.NewDeployResource,
		resources.NewEnvsBulkResource,
		resources.NewPRPreviewResource,
		// Server operations
		resources.NewServerValidateResource,
		resources.NewServerHetznerResource,
		// API settings
		resources.NewAPISettingsResource,
		// Instance settings & profile
		resources.NewInstanceSettingsResource,
		resources.NewProfileResource,
		// Team management
		resources.NewTeamResource,
		resources.NewTeamMemberResource,
		// Cloud token operations
		resources.NewCloudTokenValidateResource,
		// Scheduled task executions
		resources.NewApplicationScheduledTaskExecutionResource,
		resources.NewServiceScheduledTaskExecutionResource,
		// Backup executions
		resources.NewBackupExecutionResource,
	}
}

func (p *CoolifyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewTeamDataSource,
		datasources.NewProjectDataSource,
		datasources.NewServerDataSource,
		datasources.NewPrivateKeyDataSource,
		datasources.NewApplicationDataSource,
		datasources.NewServiceDataSource,
		datasources.NewDatabaseDataSource,
		datasources.NewDeploymentDataSource,
		datasources.NewHetznerLocationsDataSource,
		datasources.NewHetznerServerTypesDataSource,
		datasources.NewServerResourcesDataSource,
		datasources.NewServerDomainsDataSource,
		datasources.NewCoolifyResourcesDataSource,
		datasources.NewSystemInfoDataSource,
		datasources.NewDeploymentsDataSource,
		datasources.NewApplicationDeploymentsDataSource,
		// List datasources
		datasources.NewProjectsListDataSource,
		datasources.NewApplicationsListDataSource,
		datasources.NewServicesListDataSource,
		datasources.NewDatabasesListDataSource,
		datasources.NewServersListDataSource,
		datasources.NewPrivateKeysListDataSource,
		datasources.NewGitHubAppsListDataSource,
		datasources.NewCloudTokensListDataSource,
		datasources.NewTeamsListDataSource,
		datasources.NewTeamMembersDataSource,
		// Application data
		datasources.NewApplicationLogsDataSource,
		// Backup executions
		datasources.NewBackupExecutionsDataSource,
		// Hetzner extended
		datasources.NewHetznerImagesDataSource,
		datasources.NewHetznerSSHKeysDataSource,
		// Profile & settings
		datasources.NewProfileDataSource,
		datasources.NewInstanceSettingsDataSource,
		// Current team
		datasources.NewCurrentTeamDataSource,
		// Scheduled task lookups
		datasources.NewApplicationScheduledTaskDataSource,
		datasources.NewServiceScheduledTaskDataSource,
		datasources.NewApplicationScheduledTaskExecutionsDataSource,
		datasources.NewServiceScheduledTaskExecutionsDataSource,
		// Scheduled task lists
		datasources.NewApplicationScheduledTasksListDataSource,
		datasources.NewServiceScheduledTasksListDataSource,
		// Database backups list
		datasources.NewDatabaseBackupsListDataSource,
		// GitHub App extended
		datasources.NewGitHubAppRepositoriesDataSource,
		datasources.NewGitHubAppBranchesDataSource,
		// Current team members
		datasources.NewCurrentTeamMembersDataSource,
		// Project environments list
		datasources.NewProjectEnvironmentsDataSource,
	}
}

func (p *CoolifyProvider) ValidateConfig(_ context.Context, _ provider.ValidateConfigRequest, _ *provider.ValidateConfigResponse) {
}

func (p *CoolifyProvider) String() string {
	return fmt.Sprintf("CoolifyProvider(%s)", p.version)
}
