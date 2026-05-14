package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewMariaDBDatabaseResource() resource.Resource {
	return newGenericResource("database_mariadb", "MariaDB database", "/api/v1/databases/mariadb", "/api/v1/databases/%s", []resourceField{
		stringField("server_uuid", true, false, false),
		stringField("project_uuid", true, false, false),
		stringField("environment_name", true, false, false),
		stringField("name", false, true, false),
		stringField("mysql_user", false, true, false),
		stringField("mysql_password", false, true, true),
		stringField("mysql_database", false, true, false),
		stringField("mysql_root_password", false, true, true),
		boolField("instant_deploy", false, true, false),
	})
}
