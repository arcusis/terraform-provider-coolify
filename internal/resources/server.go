package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewServerResource() resource.Resource {
	return newGenericResource("server", "server", "/api/v1/servers", "/api/v1/servers/%s", []resourceField{
		stringField("name", true, false, false),
		stringField("ip", true, false, false),
		int64Field("port", true, false, false),
		stringField("user", true, false, false),
		stringField("private_key_uuid", true, false, false),
		stringField("description", false, true, false),
		boolField("is_build_server", false, true, false),
		boolField("instant_validate", false, true, false),
		stringField("proxy_type", false, true, false),
	})
}
