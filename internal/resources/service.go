package resources

import (
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewServiceResource() resource.Resource {
	r := newGenericResource("service", "service", "/api/v1/services", "/api/v1/services/%s", []resourceField{
		// type and docker_compose_raw are mutually exclusive on create — the
		// createBodyTransform below drops type when docker_compose_raw is set.
		{Name: "type", Kind: kindString, Optional: true, Send: true, ForceNew: true, SkipAPIRead: true, SkipUpdate: true, Description: "One-click service type (e.g. infisical, ghost). Omit when using docker_compose_raw."},
		{Name: "project_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true, SkipUpdate: true, Description: "UUID of the project."},
		{Name: "server_uuid", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true, SkipUpdate: true, Description: "UUID of the server."},
		{Name: "environment_name", Kind: kindString, Required: true, Send: true, ForceNew: true, SkipAPIRead: true, SkipUpdate: true, Description: "Name of the environment."},
		stringField("name", false, true, false),
		stringField("description", false, true, false),
		boolField("instant_deploy", false, true, false),
		{Name: "docker_compose_raw", Kind: kindString, Optional: true, Send: true, SkipAPIRead: true, Description: "Raw Docker Compose YAML. Coolify normalises the returned YAML, so this field is write-only in state to avoid perpetual drift."},
		// fqdn cannot be set on creation — Coolify only accepts it via PATCH.
		{Name: "fqdn", Kind: kindString, Optional: true, Send: true, SkipCreate: true, Description: "Public FQDN for the service (e.g. https://secrets.arcusis.com). Set after creation via update."},
	})
	gr := r.(*genericResource)
	encodeCompose := func(body map[string]any) map[string]any {
		if v, ok := body["docker_compose_raw"].(string); ok && v != "" {
			body["docker_compose_raw"] = base64.StdEncoding.EncodeToString([]byte(v))
		}
		return body
	}
	gr.createBodyTransform = func(body map[string]any) map[string]any {
		if _, hasCompose := body["docker_compose_raw"]; hasCompose {
			delete(body, "type")
		}
		return encodeCompose(body)
	}
	gr.updateBodyTransform = func(_ string, body map[string]any) map[string]any {
		return encodeCompose(body)
	}
	return gr
}
