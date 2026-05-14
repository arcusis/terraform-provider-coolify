---
page_title: "coolify_service Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Deploys a Coolify one-click service stack such as Ghost, WordPress, Plausible, Umami, or any other supported service type.
---

# coolify_service (Resource)

Deploys a Coolify one-click service stack. Services are pre-configured multi-container application stacks (e.g. Ghost, WordPress, Plausible Analytics) that Coolify manages as a single unit. The `type` field determines which stack is deployed.

## Example Usage

### Ghost Blog

```hcl
resource "coolify_service" "blog" {
  type             = "ghost"
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "blog"
  instant_deploy   = true
}
```

### Plausible Analytics

```hcl
resource "coolify_service" "analytics" {
  type             = "plausible"
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "analytics"
  description      = "Self-hosted Plausible Analytics"
}
```

### Custom Docker Compose

```hcl
resource "coolify_service" "custom" {
  type             = "compose"
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "custom-stack"
  docker_compose_raw = file("${path.module}/docker-compose.yml")
  instant_deploy   = true
}
```

## Schema

### Required

- `environment_name` (String) Name of the environment to deploy into (e.g. `production`, `staging`).
- `project_uuid` (String) UUID of the project this service belongs to.
- `server_uuid` (String) UUID of the server to deploy on.
- `type` (String) One-click service type to deploy (e.g. `ghost`, `wordpress`, `plausible`, `umami`, `meilisearch`, `compose`). Refer to the Coolify documentation for the full list of supported service types.

### Optional

- `description` (String) Optional description shown in the Coolify dashboard.
- `docker_compose_raw` (String) Raw Docker Compose YAML content. Used when `type` is `compose` or to override the default compose file for a service type.
- `instant_deploy` (Boolean) When `true`, Coolify immediately starts the service after creation. Defaults to `false`.
- `name` (String) Display name for the service. If omitted, Coolify generates a name from the service type.

### Read-Only

- `id` (String) Coolify resource UUID. Use this to reference the service in other resources.

## Import

Import an existing service using its UUID:

```shell
terraform import coolify_service.blog <service-uuid>
```
