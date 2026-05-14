package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewCloudTokenResource() resource.Resource {
	return newGenericResource("cloud_token", "cloud token", "/api/v1/cloud-tokens", "/api/v1/cloud-tokens/%s", []resourceField{
		stringField("name", true, false, false),
		stringField("cloud_provider", true, false, false), // "provider" is a reserved Terraform meta-argument
		{Name: "token", Kind: kindString, Required: true, Sensitive: true, Send: true, SkipAPIRead: true},
	})
}
