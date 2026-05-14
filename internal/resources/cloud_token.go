package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewCloudTokenResource() resource.Resource {
	return newGenericResource("cloud_token", "cloud token", "/api/v1/cloud-tokens", "/api/v1/cloud-tokens/%s", []resourceField{
		stringField("name", true, false, false),
		stringField("provider", true, false, false),
		{Name: "token", Kind: kindString, Required: true, Sensitive: true, Send: true, SkipAPIRead: true},
	})
}
