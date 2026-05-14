package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewPrivateKeyResource() resource.Resource {
	return &genericResource{
		typeName:    "private_key",
		displayName: "private key",
		createPath:  func(map[string]string) string { return "/api/v1/security/keys" },
		readPath:    func(id string) string { return "/api/v1/security/keys/" + id },
		// Coolify PATCH /security/keys has no UUID in path — must inject UUID into body
		updatePath: func(_ string) string { return "/api/v1/security/keys" },
		deletePath: func(id string) string { return "/api/v1/security/keys/" + id },
		updateBodyTransform: func(id string, body map[string]any) map[string]any {
			body["uuid"] = id
			return body
		},
		fields: []resourceField{
			stringField("name", true, false, false),
			stringField("description", false, true, false),
			// Coolify never returns private_key in GET responses — keep plan value
			{Name: "private_key", Kind: kindString, Required: true, Sensitive: true, Send: true, SkipAPIRead: true},
		},
	}
}
