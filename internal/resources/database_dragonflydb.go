package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewDragonflyDatabaseResource() resource.Resource {
	return newGenericResource("database_dragonfly", "Dragonfly database", "/api/v1/databases/dragonfly", "/api/v1/databases/%s", []resourceField{
		stringField("server_uuid", true, false, false),
		stringField("project_uuid", true, false, false),
		stringField("environment_name", true, false, false),
		stringField("name", false, true, false),
		stringField("dragonfly_password", false, true, true),
		boolField("instant_deploy", false, true, false),
	})
}
