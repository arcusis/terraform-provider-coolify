package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewMongoDBDatabaseResource() resource.Resource {
	return newGenericResource("database_mongodb", "MongoDB database", "/api/v1/databases/mongodb", "/api/v1/databases/%s", []resourceField{
		{Name: "server_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true},
		{Name: "project_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true},
		{Name: "environment_name", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true},
		stringField("name", false, true, false),
		stringField("mongo_initdb_root_username", false, true, false),
		stringField("mongo_initdb_root_password", false, true, true),
		stringField("mongo_initdb_database", false, true, false),
		boolField("instant_deploy", false, true, false),
	})
}
