package datasources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ── Projects list ─────────────────────────────────────────────────────────────

type projectListItemModel struct {
	UUID        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type projectsListDataSourceModel struct {
	Projects []projectListItemModel `tfsdk:"projects"`
}

type projectsListDataSource struct{ client *coolify.Client }

func NewProjectsListDataSource() datasource.DataSource { return &projectsListDataSource{} }

func (d *projectsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *projectsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all Coolify projects. Use this data source to enumerate projects and reference their UUIDs in other resources.",
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All projects available on this Coolify instance.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":        schema.StringAttribute{Computed: true, Description: "Project UUID."},
						"name":        schema.StringAttribute{Computed: true, Description: "Project name."},
						"description": schema.StringAttribute{Computed: true, Description: "Project description."},
					},
				},
			},
		},
	}
}

func (d *projectsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *projectsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config projectsListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/projects", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify projects", err.Error())
		return
	}
	config.Projects = []projectListItemModel{}
	for _, item := range dataList(out) {
		config.Projects = append(config.Projects, projectListItemModel{
			UUID:        types.StringValue(firstString(item, "uuid", "id")),
			Name:        types.StringValue(stringFromAny(item["name"])),
			Description: types.StringValue(stringFromAny(item["description"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Applications list ─────────────────────────────────────────────────────────

type applicationListItemModel struct {
	UUID            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Status          types.String `tfsdk:"status"`
	ProjectUUID     types.String `tfsdk:"project_uuid"`
	EnvironmentName types.String `tfsdk:"environment_name"`
	ServerUUID      types.String `tfsdk:"server_uuid"`
}

type applicationsListDataSourceModel struct {
	Applications []applicationListItemModel `tfsdk:"applications"`
}

type applicationsListDataSource struct{ client *coolify.Client }

func NewApplicationsListDataSource() datasource.DataSource { return &applicationsListDataSource{} }

func (d *applicationsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_applications"
}

func (d *applicationsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all Coolify applications across all projects and environments.",
		Attributes: map[string]schema.Attribute{
			"applications": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All applications deployed on this Coolify instance.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":             schema.StringAttribute{Computed: true, Description: "Application UUID."},
						"name":             schema.StringAttribute{Computed: true, Description: "Application name."},
						"status":           schema.StringAttribute{Computed: true, Description: "Current status (e.g. running, stopped, exited)."},
						"project_uuid":     schema.StringAttribute{Computed: true, Description: "UUID of the project this application belongs to."},
						"environment_name": schema.StringAttribute{Computed: true, Description: "Name of the environment (e.g. production, staging)."},
						"server_uuid":      schema.StringAttribute{Computed: true, Description: "UUID of the server hosting this application."},
					},
				},
			},
		},
	}
}

func (d *applicationsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *applicationsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config applicationsListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/applications", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify applications", err.Error())
		return
	}
	config.Applications = []applicationListItemModel{}
	for _, item := range dataList(out) {
		config.Applications = append(config.Applications, applicationListItemModel{
			UUID:            types.StringValue(firstString(item, "uuid", "id")),
			Name:            types.StringValue(stringFromAny(item["name"])),
			Status:          types.StringValue(stringFromAny(item["status"])),
			ProjectUUID:     types.StringValue(stringFromAny(item["project_uuid"])),
			EnvironmentName: types.StringValue(stringFromAny(item["environment_name"])),
			ServerUUID:      types.StringValue(stringFromAny(item["server_uuid"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Services list ─────────────────────────────────────────────────────────────

type serviceListItemModel struct {
	UUID            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Status          types.String `tfsdk:"status"`
	ProjectUUID     types.String `tfsdk:"project_uuid"`
	EnvironmentName types.String `tfsdk:"environment_name"`
	ServerUUID      types.String `tfsdk:"server_uuid"`
}

type servicesListDataSourceModel struct {
	Services []serviceListItemModel `tfsdk:"services"`
}

type servicesListDataSource struct{ client *coolify.Client }

func NewServicesListDataSource() datasource.DataSource { return &servicesListDataSource{} }

func (d *servicesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

func (d *servicesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all Coolify services (one-click stacks such as WordPress, Ghost, Plausible) across all projects.",
		Attributes: map[string]schema.Attribute{
			"services": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All services deployed on this Coolify instance.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":             schema.StringAttribute{Computed: true, Description: "Service UUID."},
						"name":             schema.StringAttribute{Computed: true, Description: "Service name."},
						"status":           schema.StringAttribute{Computed: true, Description: "Current status."},
						"project_uuid":     schema.StringAttribute{Computed: true, Description: "UUID of the project."},
						"environment_name": schema.StringAttribute{Computed: true, Description: "Name of the environment."},
						"server_uuid":      schema.StringAttribute{Computed: true, Description: "UUID of the server."},
					},
				},
			},
		},
	}
}

func (d *servicesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *servicesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config servicesListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/services", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify services", err.Error())
		return
	}
	config.Services = []serviceListItemModel{}
	for _, item := range dataList(out) {
		config.Services = append(config.Services, serviceListItemModel{
			UUID:            types.StringValue(firstString(item, "uuid", "id")),
			Name:            types.StringValue(stringFromAny(item["name"])),
			Status:          types.StringValue(stringFromAny(item["status"])),
			ProjectUUID:     types.StringValue(stringFromAny(item["project_uuid"])),
			EnvironmentName: types.StringValue(stringFromAny(item["environment_name"])),
			ServerUUID:      types.StringValue(stringFromAny(item["server_uuid"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Databases list ────────────────────────────────────────────────────────────

type databaseListItemModel struct {
	UUID            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Status          types.String `tfsdk:"status"`
	Type            types.String `tfsdk:"type"`
	ProjectUUID     types.String `tfsdk:"project_uuid"`
	EnvironmentName types.String `tfsdk:"environment_name"`
	ServerUUID      types.String `tfsdk:"server_uuid"`
}

type databasesListDataSourceModel struct {
	Databases []databaseListItemModel `tfsdk:"databases"`
}

type databasesListDataSource struct{ client *coolify.Client }

func NewDatabasesListDataSource() datasource.DataSource { return &databasesListDataSource{} }

func (d *databasesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_databases"
}

func (d *databasesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all Coolify databases (PostgreSQL, MySQL, Redis, MongoDB, etc.) across all projects.",
		Attributes: map[string]schema.Attribute{
			"databases": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All databases managed by this Coolify instance.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":             schema.StringAttribute{Computed: true, Description: "Database UUID."},
						"name":             schema.StringAttribute{Computed: true, Description: "Database name."},
						"status":           schema.StringAttribute{Computed: true, Description: "Current status."},
						"type":             schema.StringAttribute{Computed: true, Description: "Database engine type (postgresql, mysql, redis, mongodb, mariadb, keydb, dragonfly, clickhouse)."},
						"project_uuid":     schema.StringAttribute{Computed: true, Description: "UUID of the project."},
						"environment_name": schema.StringAttribute{Computed: true, Description: "Name of the environment."},
						"server_uuid":      schema.StringAttribute{Computed: true, Description: "UUID of the server."},
					},
				},
			},
		},
	}
}

func (d *databasesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *databasesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config databasesListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/databases", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify databases", err.Error())
		return
	}
	config.Databases = []databaseListItemModel{}
	for _, item := range dataList(out) {
		config.Databases = append(config.Databases, databaseListItemModel{
			UUID:            types.StringValue(firstString(item, "uuid", "id")),
			Name:            types.StringValue(stringFromAny(item["name"])),
			Status:          types.StringValue(stringFromAny(item["status"])),
			Type:            types.StringValue(stringFromAny(item["type"])),
			ProjectUUID:     types.StringValue(stringFromAny(item["project_uuid"])),
			EnvironmentName: types.StringValue(stringFromAny(item["environment_name"])),
			ServerUUID:      types.StringValue(stringFromAny(item["server_uuid"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Servers list ──────────────────────────────────────────────────────────────

type serverListItemModel struct {
	UUID        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	IP          types.String `tfsdk:"ip"`
	Status      types.String `tfsdk:"status"`
	Description types.String `tfsdk:"description"`
}

type serversListDataSourceModel struct {
	Servers []serverListItemModel `tfsdk:"servers"`
}

type serversListDataSource struct{ client *coolify.Client }

func NewServersListDataSource() datasource.DataSource { return &serversListDataSource{} }

func (d *serversListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servers"
}

func (d *serversListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all servers registered with this Coolify instance. Use this to discover server UUIDs for use in other resources.",
		Attributes: map[string]schema.Attribute{
			"servers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All servers registered with this Coolify instance.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":        schema.StringAttribute{Computed: true, Description: "Server UUID."},
						"name":        schema.StringAttribute{Computed: true, Description: "Server name."},
						"ip":          schema.StringAttribute{Computed: true, Description: "Server IP address or hostname."},
						"status":      schema.StringAttribute{Computed: true, Description: "Connection status (reachable, unreachable)."},
						"description": schema.StringAttribute{Computed: true, Description: "Optional description."},
					},
				},
			},
		},
	}
}

func (d *serversListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *serversListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config serversListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/servers", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify servers", err.Error())
		return
	}
	config.Servers = []serverListItemModel{}
	for _, item := range dataList(out) {
		ip := stringFromAny(item["ip"])
		if ip == "" {
			ip = stringFromAny(item["hostname"])
		}
		config.Servers = append(config.Servers, serverListItemModel{
			UUID:        types.StringValue(firstString(item, "uuid", "id")),
			Name:        types.StringValue(stringFromAny(item["name"])),
			IP:          types.StringValue(ip),
			Status:      types.StringValue(stringFromAny(item["status"])),
			Description: types.StringValue(stringFromAny(item["description"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Private keys list ─────────────────────────────────────────────────────────

type privateKeyListItemModel struct {
	UUID        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type privateKeysListDataSourceModel struct {
	PrivateKeys []privateKeyListItemModel `tfsdk:"private_keys"`
}

type privateKeysListDataSource struct{ client *coolify.Client }

func NewPrivateKeysListDataSource() datasource.DataSource { return &privateKeysListDataSource{} }

func (d *privateKeysListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_keys"
}

func (d *privateKeysListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all SSH private keys stored in Coolify. Note: the actual private key material is never returned by the API.",
		Attributes: map[string]schema.Attribute{
			"private_keys": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All private keys stored in this Coolify instance.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":        schema.StringAttribute{Computed: true, Description: "Private key UUID."},
						"name":        schema.StringAttribute{Computed: true, Description: "Key name."},
						"description": schema.StringAttribute{Computed: true, Description: "Key description."},
					},
				},
			},
		},
	}
}

func (d *privateKeysListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *privateKeysListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config privateKeysListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/security/keys", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify private keys", err.Error())
		return
	}
	config.PrivateKeys = []privateKeyListItemModel{}
	for _, item := range dataList(out) {
		config.PrivateKeys = append(config.PrivateKeys, privateKeyListItemModel{
			UUID:        types.StringValue(firstString(item, "uuid", "id")),
			Name:        types.StringValue(stringFromAny(item["name"])),
			Description: types.StringValue(stringFromAny(item["description"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── GitHub Apps list ──────────────────────────────────────────────────────────

type githubAppListItemModel struct {
	UUID         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	AppID        types.Int64  `tfsdk:"app_id"`
	IsSystemWide types.Bool   `tfsdk:"is_system_wide"`
}

type githubAppsListDataSourceModel struct {
	GitHubApps []githubAppListItemModel `tfsdk:"github_apps"`
}

type githubAppsListDataSource struct{ client *coolify.Client }

func NewGitHubAppsListDataSource() datasource.DataSource { return &githubAppsListDataSource{} }

func (d *githubAppsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_github_apps"
}

func (d *githubAppsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all GitHub App integrations configured in Coolify. Use this to look up existing app UUIDs for source control connections.",
		Attributes: map[string]schema.Attribute{
			"github_apps": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All GitHub App integrations.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":           schema.StringAttribute{Computed: true, Description: "GitHub App UUID."},
						"name":           schema.StringAttribute{Computed: true, Description: "GitHub App name."},
						"app_id":         schema.Int64Attribute{Computed: true, Description: "GitHub App ID."},
						"is_system_wide": schema.BoolAttribute{Computed: true, Description: "Whether this app is available to all teams."},
					},
				},
			},
		},
	}
}

func (d *githubAppsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *githubAppsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config githubAppsListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/github-apps", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify GitHub Apps", err.Error())
		return
	}
	config.GitHubApps = []githubAppListItemModel{}
	for _, item := range dataList(out) {
		config.GitHubApps = append(config.GitHubApps, githubAppListItemModel{
			UUID:         types.StringValue(firstString(item, "uuid", "id")),
			Name:         types.StringValue(stringFromAny(item["name"])),
			AppID:        types.Int64Value(int64FromAny(item["app_id"])),
			IsSystemWide: types.BoolValue(boolFromAny(item["is_system_wide"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Cloud tokens list ─────────────────────────────────────────────────────────

type cloudTokenListItemModel struct {
	UUID         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
}

type cloudTokensListDataSourceModel struct {
	CloudTokens []cloudTokenListItemModel `tfsdk:"cloud_tokens"`
}

type cloudTokensListDataSource struct{ client *coolify.Client }

func NewCloudTokensListDataSource() datasource.DataSource { return &cloudTokensListDataSource{} }

func (d *cloudTokensListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_tokens"
}

func (d *cloudTokensListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all cloud API tokens (Hetzner, DigitalOcean, etc.) stored in Coolify. Use their UUIDs with the Hetzner data sources.",
		Attributes: map[string]schema.Attribute{
			"cloud_tokens": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All cloud provider tokens stored in this Coolify instance.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":           schema.StringAttribute{Computed: true, Description: "Cloud token UUID."},
						"name":           schema.StringAttribute{Computed: true, Description: "Token name."},
						"cloud_provider": schema.StringAttribute{Computed: true, Description: "Cloud provider (hetzner, digitalocean)."},
					},
				},
			},
		},
	}
}

func (d *cloudTokensListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *cloudTokensListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config cloudTokensListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/cloud-tokens", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify cloud tokens", err.Error())
		return
	}
	config.CloudTokens = []cloudTokenListItemModel{}
	for _, item := range dataList(out) {
		provider := stringFromAny(item["provider"])
		if provider == "" {
			provider = stringFromAny(item["cloud_provider"])
		}
		config.CloudTokens = append(config.CloudTokens, cloudTokenListItemModel{
			UUID:          types.StringValue(firstString(item, "uuid", "id")),
			Name:          types.StringValue(stringFromAny(item["name"])),
			CloudProvider: types.StringValue(provider),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Teams list ────────────────────────────────────────────────────────────────

type teamListItemModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type teamsListDataSourceModel struct {
	Teams []teamListItemModel `tfsdk:"teams"`
}

type teamsListDataSource struct{ client *coolify.Client }

func NewTeamsListDataSource() datasource.DataSource { return &teamsListDataSource{} }

func (d *teamsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

func (d *teamsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all teams on this Coolify instance that the authenticated token has access to.",
		Attributes: map[string]schema.Attribute{
			"teams": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All teams accessible to the current API token.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true, Description: "Team ID."},
						"name":        schema.StringAttribute{Computed: true, Description: "Team name."},
						"description": schema.StringAttribute{Computed: true, Description: "Team description."},
					},
				},
			},
		},
	}
}

func (d *teamsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *teamsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config teamsListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/teams", &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify teams", err.Error())
		return
	}
	config.Teams = []teamListItemModel{}
	for _, item := range dataList(out) {
		config.Teams = append(config.Teams, teamListItemModel{
			ID:          types.StringValue(firstString(item, "id", "uuid")),
			Name:        types.StringValue(stringFromAny(item["name"])),
			Description: types.StringValue(stringFromAny(item["description"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Current team ─────────────────────────────────────────────────────────────

type currentTeamDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type currentTeamDataSource struct{ client *coolify.Client }

func NewCurrentTeamDataSource() datasource.DataSource { return &currentTeamDataSource{} }

func (d *currentTeamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_current_team"
}

func (d *currentTeamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns the team associated with the current API token via `GET /api/v1/teams/current`. Equivalent to `coolify_team` with `name = \"Root Team\"` on most installations, but more explicit.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true, Description: "Team ID."},
			"name":        schema.StringAttribute{Computed: true, Description: "Team name."},
			"description": schema.StringAttribute{Computed: true, Description: "Team description."},
		},
	}
}

func (d *currentTeamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *currentTeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config currentTeamDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out map[string]any
	if err := d.client.Get(ctx, "/api/v1/teams/current", &out); err != nil {
		resp.Diagnostics.AddError("Unable to read current Coolify team", err.Error())
		return
	}
	data := objectPayload(out)
	config.ID = types.StringValue(firstString(data, "id", "uuid"))
	config.Name = types.StringValue(stringFromAny(data["name"]))
	config.Description = types.StringValue(stringFromAny(data["description"]))
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Team members ──────────────────────────────────────────────────────────────

type teamMemberModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

type teamMembersDataSourceModel struct {
	TeamID  types.String      `tfsdk:"team_id"`
	Members []teamMemberModel `tfsdk:"members"`
}

type teamMembersDataSource struct{ client *coolify.Client }

func NewTeamMembersDataSource() datasource.DataSource { return &teamMembersDataSource{} }

func (d *teamMembersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_members"
}

func (d *teamMembersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all members of a Coolify team. Use the team ID from `coolify_team` or `coolify_teams`.",
		Attributes: map[string]schema.Attribute{
			"team_id": schema.StringAttribute{
				Required:    true,
				Description: "The team ID to list members for.",
			},
			"members": schema.ListNestedAttribute{
				Computed:    true,
				Description: "All members of the team.",
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

func (d *teamMembersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *teamMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config teamMembersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/teams/%s/members", config.TeamID.ValueString()), &out); err != nil {
		resp.Diagnostics.AddError("Unable to list Coolify team members", err.Error())
		return
	}
	config.Members = []teamMemberModel{}
	for _, item := range dataList(out) {
		role := stringFromAny(item["role"])
		if role == "" {
			role = stringFromAny(item["pivot"])
		}
		config.Members = append(config.Members, teamMemberModel{
			ID:    types.StringValue(firstString(item, "id", "uuid")),
			Name:  types.StringValue(stringFromAny(item["name"])),
			Email: types.StringValue(stringFromAny(item["email"])),
			Role:  types.StringValue(role),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
