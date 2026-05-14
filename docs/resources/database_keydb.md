---
page_title: "coolify_database_keydb Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Provisions a KeyDB container on a Coolify server within a project environment. KeyDB is a high-performance, Redis-compatible key-value store.
---

# coolify_database_keydb (Resource)

Provisions a KeyDB container on a Coolify-managed server. KeyDB is a Redis-compatible key-value store with multi-threading support. It can be used as a drop-in replacement for Redis. The instance is created within the specified project and environment.

## Example Usage

```hcl
resource "coolify_database_keydb" "cache" {
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "cache"
  keydb_password   = var.keydb_password
  instant_deploy   = true
}
```

## Schema

### Required

- `environment_name` (String) Name of the environment to provision the database in (e.g. `production`, `staging`).
- `project_uuid` (String) UUID of the project this database belongs to.
- `server_uuid` (String) UUID of the server to provision the database on.

### Optional

- `instant_deploy` (Boolean) When `true`, Coolify immediately starts the database container after creation. Defaults to `false`.
- `keydb_password` (String, Sensitive) Password for KeyDB authentication. If omitted, the instance will start without authentication.
- `name` (String) Display name for the database. If omitted, Coolify generates a name.

### Read-Only

- `id` (String) Coolify resource UUID.

## Import

Import an existing KeyDB database using its UUID:

```shell
terraform import coolify_database_keydb.cache <database-uuid>
```
