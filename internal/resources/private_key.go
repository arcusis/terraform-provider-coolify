package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewPrivateKeyResource() resource.Resource {
	return newGenericResource("private_key", "private key", "/api/v1/security/keys", "/api/v1/security/keys/%s", []resourceField{
		stringField("name", true, false, false),
		stringField("description", true, false, false),
		stringField("private_key", true, false, true),
	})
}
