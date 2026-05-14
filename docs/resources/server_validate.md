---
page_title: "coolify_server_validate Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Validates that Coolify can reach a server over SSH. Use as a dependency gate for resources that require a healthy server.
---

# coolify_server_validate (Resource)

Triggers the Coolify server validation check (`GET /api/v1/servers/{uuid}/validate`) and records the result. Use this as a dependency to ensure a server is reachable before creating applications or databases on it.

A validation failure does not cause `terraform apply` to fail — the result is surfaced in the `status` output so you can decide how to proceed.

## Example Usage

```hcl
resource "coolify_server" "prod" {
  name             = "prod-1"
  ip               = "10.0.0.10"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.deploy.id
}

resource "coolify_server_validate" "prod" {
  server_uuid = coolify_server.prod.id
}

# Application only deploys if the server is reachable
resource "coolify_application" "web" {
  depends_on       = [coolify_server_validate.prod]
  project_uuid     = coolify_project.myapp.id
  server_uuid      = coolify_server.prod.id
  environment_name = "production"
  type             = "dockerimage"
  docker_registry_image_name = "nginx:latest"
  ports_exposes    = "80"
}

output "server_status" {
  value = coolify_server_validate.prod.status
}
```

## Schema

### Required

- `server_uuid` (String) UUID of the server to validate.

### Read-Only

- `status` (String) Validation result returned by Coolify. Contains `ok` when the server is reachable, `unreachable` on connection failure.
