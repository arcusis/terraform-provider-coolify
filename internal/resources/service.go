package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewServiceResource() resource.Resource {
	return newGenericResource("service", "service", "/api/v1/services", "/api/v1/services/%s", []resourceField{
		{Name: "type", Kind: kindString, Required: true, Send: true, ForceNew: true, Description: "Service stack type (e.g. ghost, wordpress, plausible)."},
		{Name: "project_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, Description: "UUID of the project."},
		{Name: "server_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, Description: "UUID of the server."},
		{Name: "environment_name", Kind: kindString, Required: true, Send: true, ForceNew: true, Description: "Name of the environment."},
		stringField("name", false, true, false),
		stringField("description", false, true, false),
		boolField("instant_deploy", false, true, false),
		stringField("docker_compose_raw", false, true, false),
	})
}
