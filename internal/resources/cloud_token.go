package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewCloudTokenResource() resource.Resource {
	return newGenericResource("cloud_token", "cloud token", "/api/v1/cloud-tokens", "/api/v1/cloud-tokens/%s", []resourceField{
		stringField("name", true, false, false),
		// "provider" is a reserved Terraform meta-arg; expose as cloud_provider, send as "provider" to API
		{Name: "cloud_provider", APIName: "provider", Kind: kindString, Required: true, Send: true},
		{Name: "token", Kind: kindString, Required: true, Sensitive: true, Send: true, SkipAPIRead: true},
	})
}
