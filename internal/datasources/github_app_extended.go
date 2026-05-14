package datasources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ── GitHub App repositories ───────────────────────────────────────────────────

type githubRepoModel struct {
	FullName    types.String `tfsdk:"full_name"`
	Description types.String `tfsdk:"description"`
	Private     types.Bool   `tfsdk:"private"`
	DefaultBranch types.String `tfsdk:"default_branch"`
}

type githubAppRepositoriesDataSourceModel struct {
	GitHubAppUUID types.String      `tfsdk:"github_app_uuid"`
	Repositories  []githubRepoModel `tfsdk:"repositories"`
}

type githubAppRepositoriesDataSource struct{ client *coolify.Client }

func NewGitHubAppRepositoriesDataSource() datasource.DataSource {
	return &githubAppRepositoriesDataSource{}
}

func (d *githubAppRepositoriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_github_app_repositories"
}

func (d *githubAppRepositoriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all GitHub repositories accessible via a Coolify GitHub App integration. Requires a GitHub App that has been installed on one or more repositories. Use this to discover repository names before creating applications.",
		Attributes: map[string]schema.Attribute{
			"github_app_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the `coolify_github_app` integration.",
			},
			"repositories": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Repositories the GitHub App has access to.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"full_name":      schema.StringAttribute{Computed: true, Description: "Full repository name in `owner/repo` format."},
						"description":    schema.StringAttribute{Computed: true, Description: "Repository description."},
						"private":        schema.BoolAttribute{Computed: true, Description: "Whether the repository is private."},
						"default_branch": schema.StringAttribute{Computed: true, Description: "Default branch name (e.g. `main`)."},
					},
				},
			},
		},
	}
}

func (d *githubAppRepositoriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *githubAppRepositoriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config githubAppRepositoriesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	path := fmt.Sprintf("/api/v1/github-apps/%s/repositories", config.GitHubAppUUID.ValueString())
	if err := d.client.Get(ctx, path, &out); err != nil {
		config.Repositories = []githubRepoModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}
	config.Repositories = []githubRepoModel{}
	for _, item := range dataList(out) {
		config.Repositories = append(config.Repositories, githubRepoModel{
			FullName:      types.StringValue(stringFromAny(item["full_name"])),
			Description:   types.StringValue(stringFromAny(item["description"])),
			Private:       types.BoolValue(boolFromAny(item["private"])),
			DefaultBranch: types.StringValue(stringFromAny(item["default_branch"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── GitHub App repository branches ───────────────────────────────────────────

type githubBranchModel struct {
	Name      types.String `tfsdk:"name"`
	Protected types.Bool   `tfsdk:"protected"`
}

type githubAppBranchesDataSourceModel struct {
	GitHubAppUUID types.String        `tfsdk:"github_app_uuid"`
	Owner         types.String        `tfsdk:"owner"`
	Repository    types.String        `tfsdk:"repository"`
	Branches      []githubBranchModel `tfsdk:"branches"`
}

type githubAppBranchesDataSource struct{ client *coolify.Client }

func NewGitHubAppBranchesDataSource() datasource.DataSource {
	return &githubAppBranchesDataSource{}
}

func (d *githubAppBranchesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_github_app_branches"
}

func (d *githubAppBranchesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all branches of a GitHub repository accessible via a Coolify GitHub App. Use this to discover branch names before configuring application deployments.",
		Attributes: map[string]schema.Attribute{
			"github_app_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the `coolify_github_app` integration.",
			},
			"owner": schema.StringAttribute{
				Required:    true,
				Description: "GitHub repository owner (organization or user login).",
			},
			"repository": schema.StringAttribute{
				Required:    true,
				Description: "Repository name (without the owner prefix).",
			},
			"branches": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Branches available in the repository.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":      schema.StringAttribute{Computed: true, Description: "Branch name."},
						"protected": schema.BoolAttribute{Computed: true, Description: "Whether the branch has branch protection rules."},
					},
				},
			},
		},
	}
}

func (d *githubAppBranchesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *githubAppBranchesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config githubAppBranchesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	path := fmt.Sprintf("/api/v1/github-apps/%s/repositories/%s/%s/branches",
		config.GitHubAppUUID.ValueString(),
		config.Owner.ValueString(),
		config.Repository.ValueString(),
	)
	if err := d.client.Get(ctx, path, &out); err != nil {
		config.Branches = []githubBranchModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}
	config.Branches = []githubBranchModel{}
	for _, item := range dataList(out) {
		config.Branches = append(config.Branches, githubBranchModel{
			Name:      types.StringValue(stringFromAny(item["name"])),
			Protected: types.BoolValue(boolFromAny(item["protected"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
