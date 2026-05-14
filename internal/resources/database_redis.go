package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewRedisDatabaseResource() resource.Resource {
	return newGenericResource("database_redis", "Redis database", "/api/v1/databases/redis", "/api/v1/databases/%s", []resourceField{
		stringField("server_uuid", true, false, false),
		stringField("project_uuid", true, false, false),
		stringField("environment_name", true, false, false),
		stringField("redis_password", false, true, true),
	})
}
