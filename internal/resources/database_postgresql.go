package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewPostgreSQLDatabaseResource() resource.Resource {
	return newGenericResource("database_postgresql", "PostgreSQL database", "/api/v1/databases/postgresql", "/api/v1/databases/%s", []resourceField{
		{Name: "server_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true},
		{Name: "project_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true},
		{Name: "environment_name", Kind: kindString, Required: true, Send: true, ForceNew: true},
		stringField("name", false, true, false),
		stringField("postgres_user", false, true, false),
		stringField("postgres_password", false, true, true),
		stringField("postgres_db", false, true, false),
		boolField("instant_deploy", false, true, false),
	})
}
