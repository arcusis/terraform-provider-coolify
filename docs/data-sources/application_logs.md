---
page_title: "coolify_application_logs Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Fetches recent log output from a Coolify application container. Useful for inspecting startup errors after a deployment.
---

# coolify_application_logs (Data Source)

Reads the most recent log output from an application container via `GET /api/v1/applications/{uuid}/logs`. Useful for checking startup errors or last-run output without leaving Terraform.

## Example Usage

```hcl
resource "coolify_application" "web" {
  type                       = "dockerimage"
  project_uuid               = coolify_project.myapp.id
  server_uuid                = var.server_uuid
  environment_name           = "production"
  docker_registry_image_name = "nginx:latest"
  ports_exposes              = "80"
  instant_deploy             = true
}

data "coolify_application_logs" "web" {
  application_uuid = coolify_application.web.id
  depends_on       = [coolify_application.web]
}

output "web_logs" {
  value = data.coolify_application_logs.web.logs
}
```

## Schema

### Required

- `application_uuid` (String) UUID of the application to fetch logs for.

### Read-Only

- `logs` (String) Raw log output from the application container. Empty string if the application has never started.
