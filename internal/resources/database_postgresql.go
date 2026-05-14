package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewPostgreSQLDatabaseResource() resource.Resource {
	return newGenericResource("database_postgresql", "PostgreSQL database", "/api/v1/databases/postgresql", "/api/v1/databases/%s", []resourceField{
		stringField("server_uuid", true, false, false),
		stringField("project_uuid", true, false, false),
		stringField("environment_name", true, false, false),
		stringField("name", false, true, false),
		stringField("postgres_user", false, true, false),
		stringField("postgres_password", false, true, true),
		stringField("postgres_db", false, true, false),
		boolField("instant_deploy", false, true, false),
	})
}
