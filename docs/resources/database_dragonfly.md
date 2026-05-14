---
page_title: "coolify_database_dragonfly Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Provisions a DragonflyDB container on a Coolify server within a project environment. DragonflyDB is a high-throughput, Redis-compatible in-memory data store.
---

# coolify_database_dragonfly (Resource)

Provisions a DragonflyDB container on a Coolify-managed server. DragonflyDB is a Redis- and Memcached-compatible in-memory data store designed for high throughput on modern hardware. It can be used as a drop-in replacement for Redis. The instance is created within the specified project and environment.

## Example Usage

```hcl
resource "coolify_database_dragonfly" "cache" {
  project_uuid       = coolify_project.myapp.id
  server_uuid        = var.server_uuid
  environment_name   = "production"
  name               = "cache"
  dragonfly_password = var.dragonfly_password
  instant_deploy     = true
}
```

## Schema

### Required

- `environment_name` (String) Name of the environment to provision the database in (e.g. `production`, `staging`).
- `project_uuid` (String) UUID of the project this database belongs to.
- `server_uuid` (String) UUID of the server to provision the database on.

### Optional

- `dragonfly_password` (String, Sensitive) Password for DragonflyDB authentication. If omitted, the instance will start without authentication.
- `instant_deploy` (Boolean) When `true`, Coolify immediately starts the database container after creation. Defaults to `false`.
- `name` (String) Display name for the database. If omitted, Coolify generates a name.

### Read-Only

- `id` (String) Coolify resource UUID.

## Import

Import an existing DragonflyDB database using its UUID:

```shell
terraform import coolify_database_dragonfly.cache <database-uuid>
```
