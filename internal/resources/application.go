package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewApplicationResource() resource.Resource {
	app := &genericResource{
		typeName:    "application",
		displayName: "application",
		createPath: func(model resourceModel) string {
			return fmt.Sprintf("/api/v1/applications/%s", model.Type.ValueString())
		},
		readPath:   func(id string) string { return fmt.Sprintf("/api/v1/applications/%s", id) },
		updatePath: func(id string) string { return fmt.Sprintf("/api/v1/applications/%s", id) },
		deletePath: func(id string) string { return fmt.Sprintf("/api/v1/applications/%s", id) },
		fields: normalizeFields([]resourceField{
			noSendStringField("type", true),
			stringField("project_uuid", true, false, false),
			stringField("server_uuid", true, false, false),
			stringField("environment_name", true, false, false),
			stringField("ports_exposes", true, false, false),
			stringField("git_repository", false, true, false),
			stringField("git_branch", false, true, false),
			stringField("build_pack", false, true, false),
			stringField("dockerfile", false, true, false),
			stringField("docker_registry_image_name", false, true, false),
			stringField("github_app_uuid", false, true, false),
			stringField("private_key_uuid", false, true, false),
			stringField("name", false, true, false),
			stringField("description", false, true, false),
			stringField("domains", false, true, false),
			stringField("install_command", false, true, false),
			stringField("build_command", false, true, false),
			stringField("start_command", false, true, false),
			boolField("is_auto_deploy_enabled", false, true, false),
			boolField("is_force_https_enabled", false, true, false),
		}),
	}
	return app
}
