package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type fieldKind string

const (
	kindString fieldKind = "string"
	kindBool   fieldKind = "bool"
	kindInt64  fieldKind = "int64"
)

type resourceField struct {
	Name        string
	Kind        fieldKind
	Required    bool
	Optional    bool
	Computed    bool
	Sensitive   bool
	Description string
	Send        bool
}

type genericResource struct {
	client      *coolify.Client
	typeName    string
	displayName string
	createPath  func(resourceModel) string
	readPath    func(string) string
	updatePath  func(string) string
	deletePath  func(string) string
	fields      []resourceField
}

type resourceModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	IP                         types.String `tfsdk:"ip"`
	Port                       types.Int64  `tfsdk:"port"`
	User                       types.String `tfsdk:"user"`
	PrivateKey                 types.String `tfsdk:"private_key"`
	PrivateKeyUUID             types.String `tfsdk:"private_key_uuid"`
	IsBuildServer              types.Bool   `tfsdk:"is_build_server"`
	InstantValidate            types.Bool   `tfsdk:"instant_validate"`
	ProxyType                  types.String `tfsdk:"proxy_type"`
	Type                       types.String `tfsdk:"type"`
	ProjectUUID                types.String `tfsdk:"project_uuid"`
	ServerUUID                 types.String `tfsdk:"server_uuid"`
	EnvironmentName            types.String `tfsdk:"environment_name"`
	PortsExposes               types.String `tfsdk:"ports_exposes"`
	GitRepository              types.String `tfsdk:"git_repository"`
	GitBranch                  types.String `tfsdk:"git_branch"`
	BuildPack                  types.String `tfsdk:"build_pack"`
	Dockerfile                 types.String `tfsdk:"dockerfile"`
	DockerRegistryImageName    types.String `tfsdk:"docker_registry_image_name"`
	GithubAppUUID              types.String `tfsdk:"github_app_uuid"`
	Domains                    types.String `tfsdk:"domains"`
	InstallCommand             types.String `tfsdk:"install_command"`
	BuildCommand               types.String `tfsdk:"build_command"`
	StartCommand               types.String `tfsdk:"start_command"`
	IsAutoDeployEnabled        types.Bool   `tfsdk:"is_auto_deploy_enabled"`
	IsForceHTTPSEnabled        types.Bool   `tfsdk:"is_force_https_enabled"`
	InstantDeploy              types.Bool   `tfsdk:"instant_deploy"`
	DockerComposeRaw           types.String `tfsdk:"docker_compose_raw"`
	PostgresUser               types.String `tfsdk:"postgres_user"`
	PostgresPassword           types.String `tfsdk:"postgres_password"`
	PostgresDB                 types.String `tfsdk:"postgres_db"`
	MySQLUser                  types.String `tfsdk:"mysql_user"`
	MySQLPassword              types.String `tfsdk:"mysql_password"`
	MySQLDatabase              types.String `tfsdk:"mysql_database"`
	MySQLRootPassword          types.String `tfsdk:"mysql_root_password"`
	RedisPassword              types.String `tfsdk:"redis_password"`
	MongoInitdbRootUsername    types.String `tfsdk:"mongo_initdb_root_username"`
	MongoInitdbRootPassword    types.String `tfsdk:"mongo_initdb_root_password"`
	MongoInitdbDatabase        types.String `tfsdk:"mongo_initdb_database"`
	ApplicationUUID            types.String `tfsdk:"application_uuid"`
	Key                        types.String `tfsdk:"key"`
	Value                      types.String `tfsdk:"value"`
	IsPreview                  types.Bool   `tfsdk:"is_preview"`
	IsLiteral                  types.Bool   `tfsdk:"is_literal"`
	IsMultiline                types.Bool   `tfsdk:"is_multiline"`
	IsShownOnce                types.Bool   `tfsdk:"is_shown_once"`
}

func NewProjectResource() resource.Resource {
	return newGenericResource("project", "project", "/api/v1/projects", "/api/v1/projects/%s", []resourceField{
		stringField("name", true, false, false),
		stringField("description", false, true, false),
	})
}

func newGenericResource(typeName, displayName, createPath, itemPath string, fields []resourceField) resource.Resource {
	return &genericResource{
		typeName:    typeName,
		displayName: displayName,
		createPath:  func(resourceModel) string { return createPath },
		readPath:    func(id string) string { return fmt.Sprintf(itemPath, id) },
		updatePath:  func(id string) string { return fmt.Sprintf(itemPath, id) },
		deletePath:  func(id string) string { return fmt.Sprintf(itemPath, id) },
		fields:      normalizeFields(fields),
	}
}

func (r *genericResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName
}

func (r *genericResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Coolify resource UUID.",
		},
	}
	for _, f := range r.fields {
		attrs[f.Name] = f.schemaAttribute()
	}
	resp.Schema = schema.Schema{Attributes: attrs}
}

func (r *genericResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*coolify.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *coolify.Client from provider configuration.")
		return
	}
	r.client = client
}

func (r *genericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(r.validate(plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := r.bodyFromModel(plan)
	var created map[string]any
	if err := r.client.Post(ctx, r.createPath(plan), body, &created); err != nil {
		resp.Diagnostics.AddError("Unable to create Coolify "+r.displayName, err.Error())
		return
	}

	id := firstString(created, "uuid", "id")
	if id == "" {
		resp.Diagnostics.AddError("Coolify create response missing UUID", "The Coolify API did not return uuid or id for the created "+r.displayName+".")
		return
	}
	plan.ID = types.StringValue(id)

	resp.Diagnostics.Append(r.readIntoState(ctx, id, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *genericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Missing Terraform resource ID", "Cannot read Coolify "+r.displayName+" because state has no ID.")
		return
	}

	diags := r.readIntoState(ctx, id, &state)
	if isNotFound(diags) {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *genericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceModel
	var state resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if id == "" {
		id = plan.ID.ValueString()
	}
	if id == "" {
		resp.Diagnostics.AddError("Missing Terraform resource ID", "Cannot update Coolify "+r.displayName+" because state has no ID.")
		return
	}
	resp.Diagnostics.Append(r.validate(plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Patch(ctx, r.updatePath(id), r.bodyFromModel(plan), nil); err != nil {
		resp.Diagnostics.AddError("Unable to update Coolify "+r.displayName, err.Error())
		return
	}
	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(r.readIntoState(ctx, id, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *genericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if id == "" {
		return
	}
	if err := r.client.Delete(ctx, r.deletePath(id), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Coolify "+r.displayName, err.Error())
	}
}

func (r *genericResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *genericResource) readIntoState(ctx context.Context, id string, state *resourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	var out map[string]any
	if err := r.client.Get(ctx, r.readPath(id), &out); err != nil {
		diags.AddError("Unable to read Coolify "+r.displayName, err.Error())
		return diags
	}

	data := objectPayload(out)
	if responseID := firstString(data, "uuid", "id"); responseID != "" {
		state.ID = types.StringValue(responseID)
	} else {
		state.ID = types.StringValue(id)
	}
	for _, f := range r.fields {
		if value, ok := data[f.Name]; ok {
			setModelField(state, f.Name, value)
		}
	}
	return diags
}

func (r *genericResource) bodyFromModel(model resourceModel) map[string]any {
	body := make(map[string]any)
	for _, f := range r.fields {
		if !f.Send {
			continue
		}
		value, ok := modelFieldValue(model, f.Name)
		if ok {
			body[f.Name] = value
		}
	}
	return body
}

func (r *genericResource) validate(model resourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if r.typeName != "application" {
		return diags
	}

	appType := model.Type.ValueString()
	switch appType {
	case "public":
		requireApplicationFields(&diags, model, "git_repository", "git_branch", "build_pack")
	case "private-github-app":
		requireApplicationFields(&diags, model, "git_repository", "git_branch", "build_pack", "github_app_uuid")
	case "private-deploy-key":
		requireApplicationFields(&diags, model, "git_repository", "git_branch", "build_pack", "private_key_uuid")
	case "dockerfile":
		requireApplicationFields(&diags, model, "dockerfile")
	case "dockerimage":
		requireApplicationFields(&diags, model, "docker_registry_image_name")
	default:
		diags.AddError("Invalid application type", "Application type must be one of public, private-github-app, private-deploy-key, dockerfile, or dockerimage.")
	}
	return diags
}

func requireApplicationFields(diags *diag.Diagnostics, model resourceModel, names ...string) {
	for _, name := range names {
		if _, ok := modelFieldValue(model, name); !ok {
			diags.AddError("Missing application field", fmt.Sprintf("Application type %q requires %s.", model.Type.ValueString(), name))
		}
	}
}

func normalizeFields(fields []resourceField) []resourceField {
	for i := range fields {
		if !fields[i].Required && !fields[i].Optional && !fields[i].Computed {
			fields[i].Optional = true
		}
		if !fields[i].Required && !fields[i].Computed {
			fields[i].Optional = true
		}
		if fields[i].Send == false && fields[i].Description == "no-send" {
			fields[i].Description = ""
		} else if !fields[i].Computed && fields[i].Send == false && fields[i].Description != "no-send" {
			fields[i].Send = true
		}
	}
	return fields
}

func (f resourceField) schemaAttribute() schema.Attribute {
	switch f.Kind {
	case kindBool:
		return schema.BoolAttribute{
			Required:    f.Required,
			Optional:    f.Optional,
			Computed:    f.Computed || (f.Optional && !f.Required),
			Sensitive:   f.Sensitive,
			Description: f.Description,
		}
	case kindInt64:
		return schema.Int64Attribute{
			Required:    f.Required,
			Optional:    f.Optional,
			Computed:    f.Computed || (f.Optional && !f.Required),
			Sensitive:   f.Sensitive,
			Description: f.Description,
		}
	default:
		return schema.StringAttribute{
			Required:    f.Required,
			Optional:    f.Optional,
			Computed:    f.Computed || (f.Optional && !f.Required),
			Sensitive:   f.Sensitive,
			Description: f.Description,
		}
	}
}

func stringField(name string, required, optional, sensitive bool) resourceField {
	return resourceField{Name: name, Kind: kindString, Required: required, Optional: optional, Sensitive: sensitive, Send: true}
}

func boolField(name string, required, optional, sensitive bool) resourceField {
	return resourceField{Name: name, Kind: kindBool, Required: required, Optional: optional, Sensitive: sensitive, Send: true}
}

func int64Field(name string, required, optional, sensitive bool) resourceField {
	return resourceField{Name: name, Kind: kindInt64, Required: required, Optional: optional, Sensitive: sensitive, Send: true}
}

func noSendStringField(name string, required bool) resourceField {
	return resourceField{Name: name, Kind: kindString, Required: required, Send: false, Description: "no-send"}
}

func modelFieldValue(model resourceModel, name string) (any, bool) {
	switch name {
	case "name":
		return stringValue(model.Name)
	case "description":
		return stringValue(model.Description)
	case "ip":
		return stringValue(model.IP)
	case "port":
		return int64Value(model.Port)
	case "user":
		return stringValue(model.User)
	case "private_key":
		return stringValue(model.PrivateKey)
	case "private_key_uuid":
		return stringValue(model.PrivateKeyUUID)
	case "is_build_server":
		return boolValue(model.IsBuildServer)
	case "instant_validate":
		return boolValue(model.InstantValidate)
	case "proxy_type":
		return stringValue(model.ProxyType)
	case "type":
		return stringValue(model.Type)
	case "project_uuid":
		return stringValue(model.ProjectUUID)
	case "server_uuid":
		return stringValue(model.ServerUUID)
	case "environment_name":
		return stringValue(model.EnvironmentName)
	case "ports_exposes":
		return stringValue(model.PortsExposes)
	case "git_repository":
		return stringValue(model.GitRepository)
	case "git_branch":
		return stringValue(model.GitBranch)
	case "build_pack":
		return stringValue(model.BuildPack)
	case "dockerfile":
		return stringValue(model.Dockerfile)
	case "docker_registry_image_name":
		return stringValue(model.DockerRegistryImageName)
	case "github_app_uuid":
		return stringValue(model.GithubAppUUID)
	case "domains":
		return stringValue(model.Domains)
	case "install_command":
		return stringValue(model.InstallCommand)
	case "build_command":
		return stringValue(model.BuildCommand)
	case "start_command":
		return stringValue(model.StartCommand)
	case "is_auto_deploy_enabled":
		return boolValue(model.IsAutoDeployEnabled)
	case "is_force_https_enabled":
		return boolValue(model.IsForceHTTPSEnabled)
	case "instant_deploy":
		return boolValue(model.InstantDeploy)
	case "docker_compose_raw":
		return stringValue(model.DockerComposeRaw)
	case "postgres_user":
		return stringValue(model.PostgresUser)
	case "postgres_password":
		return stringValue(model.PostgresPassword)
	case "postgres_db":
		return stringValue(model.PostgresDB)
	case "mysql_user":
		return stringValue(model.MySQLUser)
	case "mysql_password":
		return stringValue(model.MySQLPassword)
	case "mysql_database":
		return stringValue(model.MySQLDatabase)
	case "mysql_root_password":
		return stringValue(model.MySQLRootPassword)
	case "redis_password":
		return stringValue(model.RedisPassword)
	case "mongo_initdb_root_username":
		return stringValue(model.MongoInitdbRootUsername)
	case "mongo_initdb_root_password":
		return stringValue(model.MongoInitdbRootPassword)
	case "mongo_initdb_database":
		return stringValue(model.MongoInitdbDatabase)
	case "application_uuid":
		return stringValue(model.ApplicationUUID)
	case "key":
		return stringValue(model.Key)
	case "value":
		return stringValue(model.Value)
	case "is_preview":
		return boolValue(model.IsPreview)
	case "is_literal":
		return boolValue(model.IsLiteral)
	case "is_multiline":
		return boolValue(model.IsMultiline)
	case "is_shown_once":
		return boolValue(model.IsShownOnce)
	default:
		return nil, false
	}
}

func setModelField(model *resourceModel, name string, value any) {
	switch name {
	case "name":
		model.Name = types.StringValue(stringFromAny(value))
	case "description":
		model.Description = types.StringValue(stringFromAny(value))
	case "ip":
		model.IP = types.StringValue(stringFromAny(value))
	case "port":
		model.Port = types.Int64Value(int64FromAny(value))
	case "user":
		model.User = types.StringValue(stringFromAny(value))
	case "private_key":
		model.PrivateKey = types.StringValue(stringFromAny(value))
	case "private_key_uuid":
		model.PrivateKeyUUID = types.StringValue(stringFromAny(value))
	case "is_build_server":
		model.IsBuildServer = types.BoolValue(boolFromAny(value))
	case "instant_validate":
		model.InstantValidate = types.BoolValue(boolFromAny(value))
	case "proxy_type":
		model.ProxyType = types.StringValue(stringFromAny(value))
	case "type":
		model.Type = types.StringValue(stringFromAny(value))
	case "project_uuid":
		model.ProjectUUID = types.StringValue(stringFromAny(value))
	case "server_uuid":
		model.ServerUUID = types.StringValue(stringFromAny(value))
	case "environment_name":
		model.EnvironmentName = types.StringValue(stringFromAny(value))
	case "ports_exposes":
		model.PortsExposes = types.StringValue(stringFromAny(value))
	case "git_repository":
		model.GitRepository = types.StringValue(stringFromAny(value))
	case "git_branch":
		model.GitBranch = types.StringValue(stringFromAny(value))
	case "build_pack":
		model.BuildPack = types.StringValue(stringFromAny(value))
	case "dockerfile":
		model.Dockerfile = types.StringValue(stringFromAny(value))
	case "docker_registry_image_name":
		model.DockerRegistryImageName = types.StringValue(stringFromAny(value))
	case "github_app_uuid":
		model.GithubAppUUID = types.StringValue(stringFromAny(value))
	case "domains":
		model.Domains = types.StringValue(stringFromAny(value))
	case "install_command":
		model.InstallCommand = types.StringValue(stringFromAny(value))
	case "build_command":
		model.BuildCommand = types.StringValue(stringFromAny(value))
	case "start_command":
		model.StartCommand = types.StringValue(stringFromAny(value))
	case "is_auto_deploy_enabled":
		model.IsAutoDeployEnabled = types.BoolValue(boolFromAny(value))
	case "is_force_https_enabled":
		model.IsForceHTTPSEnabled = types.BoolValue(boolFromAny(value))
	case "instant_deploy":
		model.InstantDeploy = types.BoolValue(boolFromAny(value))
	case "docker_compose_raw":
		model.DockerComposeRaw = types.StringValue(stringFromAny(value))
	case "postgres_user":
		model.PostgresUser = types.StringValue(stringFromAny(value))
	case "postgres_password":
		model.PostgresPassword = types.StringValue(stringFromAny(value))
	case "postgres_db":
		model.PostgresDB = types.StringValue(stringFromAny(value))
	case "mysql_user":
		model.MySQLUser = types.StringValue(stringFromAny(value))
	case "mysql_password":
		model.MySQLPassword = types.StringValue(stringFromAny(value))
	case "mysql_database":
		model.MySQLDatabase = types.StringValue(stringFromAny(value))
	case "mysql_root_password":
		model.MySQLRootPassword = types.StringValue(stringFromAny(value))
	case "redis_password":
		model.RedisPassword = types.StringValue(stringFromAny(value))
	case "mongo_initdb_root_username":
		model.MongoInitdbRootUsername = types.StringValue(stringFromAny(value))
	case "mongo_initdb_root_password":
		model.MongoInitdbRootPassword = types.StringValue(stringFromAny(value))
	case "mongo_initdb_database":
		model.MongoInitdbDatabase = types.StringValue(stringFromAny(value))
	case "application_uuid":
		model.ApplicationUUID = types.StringValue(stringFromAny(value))
	case "key":
		model.Key = types.StringValue(stringFromAny(value))
	case "value":
		model.Value = types.StringValue(stringFromAny(value))
	case "is_preview":
		model.IsPreview = types.BoolValue(boolFromAny(value))
	case "is_literal":
		model.IsLiteral = types.BoolValue(boolFromAny(value))
	case "is_multiline":
		model.IsMultiline = types.BoolValue(boolFromAny(value))
	case "is_shown_once":
		model.IsShownOnce = types.BoolValue(boolFromAny(value))
	}
}

func stringValue(value types.String) (any, bool) {
	if value.IsNull() || value.IsUnknown() {
		return nil, false
	}
	return value.ValueString(), true
}

func boolValue(value types.Bool) (any, bool) {
	if value.IsNull() || value.IsUnknown() {
		return nil, false
	}
	return value.ValueBool(), true
}

func int64Value(value types.Int64) (any, bool) {
	if value.IsNull() || value.IsUnknown() {
		return nil, false
	}
	return value.ValueInt64(), true
}

func objectPayload(data map[string]any) map[string]any {
	if nested, ok := data["data"].(map[string]any); ok {
		return nested
	}
	return data
}

func firstString(data map[string]any, keys ...string) string {
	data = objectPayload(data)
	for _, key := range keys {
		if value, ok := data[key]; ok {
			if text := stringFromAny(value); text != "" {
				return text
			}
		}
	}
	return ""
}

func stringFromAny(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		return strconv.FormatBool(v)
	case []any, map[string]any:
		payload, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprint(v)
		}
		return string(payload)
	default:
		return fmt.Sprint(v)
	}
}

func boolFromAny(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		parsed, _ := strconv.ParseBool(v)
		return parsed
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return false
	}
}

func int64FromAny(value any) int64 {
	switch v := value.(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	case string:
		parsed, _ := strconv.ParseInt(v, 10, 64)
		return parsed
	default:
		return 0
	}
}

func isNotFound(diags diag.Diagnostics) bool {
	for _, item := range diags {
		if strings.Contains(item.Detail(), "HTTP 404") {
			return true
		}
	}
	return false
}
