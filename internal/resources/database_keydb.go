package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewKeyDBDatabaseResource() resource.Resource {
	return newGenericResource("database_keydb", "KeyDB database", "/api/v1/databases/keydb", "/api/v1/databases/%s", []resourceField{
		{Name: "server_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true},
		{Name: "project_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true},
		{Name: "environment_name", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true},
		stringField("name", false, true, false),
		stringField("keydb_password", false, true, true),
		boolField("instant_deploy", false, true, false),
	})
}
