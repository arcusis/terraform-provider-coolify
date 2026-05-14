package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewKeyDBDatabaseResource() resource.Resource {
	return newGenericResource("database_keydb", "KeyDB database", "/api/v1/databases/keydb", "/api/v1/databases/%s", []resourceField{
		stringField("server_uuid", true, false, false),
		stringField("project_uuid", true, false, false),
		stringField("environment_name", true, false, false),
		stringField("name", false, true, false),
		stringField("keydb_password", false, true, true),
	})
}
