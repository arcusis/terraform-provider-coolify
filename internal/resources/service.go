package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewServiceResource() resource.Resource {
	r := newGenericResource("service", "service", "/api/v1/services", "/api/v1/services/%s", []resourceField{
		// type and docker_compose_raw are mutually exclusive on create — the
		// createBodyTransform below drops type when docker_compose_raw is set.
		{Name: "type", Kind: kindString, Optional: true, Send: true, ForceNew: true, SkipAPIRead: true, Description: "One-click service type (e.g. infisical, ghost). Omit when using docker_compose_raw."},
		{Name: "project_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true, Description: "UUID of the project."},
		{Name: "server_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true, Description: "UUID of the server."},
		{Name: "environment_name", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true, Description: "Name of the environment."},
		stringField("name", false, true, false),
		stringField("description", false, true, false),
		boolField("instant_deploy", false, true, false),
		stringField("docker_compose_raw", false, true, false),
		// fqdn cannot be set on creation — Coolify only accepts it via PATCH.
		{Name: "fqdn", Kind: kindString, Optional: true, Send: true, SkipCreate: true, Description: "Public FQDN for the service (e.g. https://secrets.arcusis.com). Set after creation via update."},
	})
	gr := r.(*genericResource)
	gr.createBodyTransform = func(body map[string]any) map[string]any {
		if _, hasCompose := body["docker_compose_raw"]; hasCompose {
			delete(body, "type")
		}
		return body
	}
	return gr
}
