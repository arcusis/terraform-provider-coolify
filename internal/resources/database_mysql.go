package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewMySQLDatabaseResource() resource.Resource {
	return newGenericResource("database_mysql", "MySQL database", "/api/v1/databases/mysql", "/api/v1/databases/%s", []resourceField{
		{Name: "server_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true},
		{Name: "project_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true},
		{Name: "environment_name", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true},
		stringField("name", false, true, false),
		stringField("mysql_user", false, true, false),
		stringField("mysql_password", false, true, true),
		stringField("mysql_database", false, true, false),
		stringField("mysql_root_password", false, true, true),
		boolField("instant_deploy", false, true, false),
	})
}
