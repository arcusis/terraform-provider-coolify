package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewClickhouseDatabaseResource() resource.Resource {
	return newGenericResource("database_clickhouse", "Clickhouse database", "/api/v1/databases/clickhouse", "/api/v1/databases/%s", []resourceField{
		stringField("server_uuid", true, false, false),
		stringField("project_uuid", true, false, false),
		stringField("environment_name", true, false, false),
		stringField("name", false, true, false),
		stringField("clickhouse_admin_user", false, true, false),
		stringField("clickhouse_admin_password", false, true, true),
		boolField("instant_deploy", false, true, false),
	})
}
