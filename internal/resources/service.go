package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewServiceResource() resource.Resource {
	return newGenericResource("service", "service", "/api/v1/services", "/api/v1/services/%s", []resourceField{
		stringField("type", true, false, false),
		stringField("project_uuid", true, false, false),
		stringField("server_uuid", true, false, false),
		stringField("environment_name", true, false, false),
		stringField("name", false, true, false),
		stringField("description", false, true, false),
		boolField("instant_deploy", false, true, false),
		stringField("docker_compose_raw", false, true, false),
	})
}
