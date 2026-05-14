package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewApplicationResource() resource.Resource {
	return &genericResource{
		typeName:    "application",
		displayName: "application",
		createPath: func(vals map[string]string) string {
			appType := vals["type"]
			if appType == "" {
				appType = "public"
			}
			return fmt.Sprintf("/api/v1/applications/%s", appType)
		},
		readPath:   func(id string) string { return fmt.Sprintf("/api/v1/applications/%s", id) },
		updatePath: func(id string) string { return fmt.Sprintf("/api/v1/applications/%s", id) },
		deletePath: func(id string) string { return fmt.Sprintf("/api/v1/applications/%s", id) },
		fields: []resourceField{
			{Name: "type", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true, Description: "Application type: dockerimage, dockerfile, public, private-gh-app, private-deploy-key, or dockercompose."},
			{Name: "project_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true, Description: "UUID of the project to deploy into."},
			{Name: "server_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true, Description: "UUID of the server to deploy on."},
			{Name: "environment_name", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true, Description: "Name of the environment (e.g. production, staging)."},
			stringField("ports_exposes", false, true, false),
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
			boolField("instant_deploy", false, true, false),
		},
	}
}
