package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewGitHubAppResource() resource.Resource {
	return &genericResource{
		typeName:    "github_app",
		displayName: "GitHub App",
		createPath:  func(map[string]string) string { return "/api/v1/github-apps" },
		readPath:    func(id string) string { return "/api/v1/github-apps" }, // list and filter
		updatePath:  func(id string) string { return fmt.Sprintf("/api/v1/github-apps/%s", id) },
		deletePath:  func(id string) string { return fmt.Sprintf("/api/v1/github-apps/%s", id) },
		fields: []resourceField{
			stringField("name", true, false, false),
			stringField("api_url", true, false, false),
			stringField("html_url", true, false, false),
			int64Field("app_id", true, false, false),
			int64Field("installation_id", true, false, false),
			stringField("client_id", true, false, false),
			{Name: "client_secret", Kind: kindString, Required: true, Sensitive: true, Send: true, SkipAPIRead: true},
			stringField("private_key_uuid", true, false, false),
			stringField("organization", false, true, false),
			stringField("webhook_secret", false, true, true),
			boolField("is_system_wide", false, true, false),
		},
	}
}
