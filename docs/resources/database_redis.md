---
page_title: "coolify_database_redis Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Provisions a Redis container on a Coolify server within a project environment.
---

# coolify_database_redis (Resource)

Provisions a Redis container on a Coolify-managed server. The database is created within the specified project and environment and can be referenced by applications running on the same server.

## Example Usage

```hcl
resource "coolify_database_redis" "cache" {
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "cache"
  redis_password   = var.redis_password
  instant_deploy   = true
}

resource "coolify_environment_variable" "redis_url" {
  application_uuid = coolify_application.api.id
  key              = "REDIS_URL"
  value            = "redis://:${var.redis_password}@${coolify_database_redis.cache.id}:6379"
}
```

## Schema

### Required

- `environment_name` (String) Name of the environment to provision the database in (e.g. `production`, `staging`).
- `project_uuid` (String) UUID of the project this database belongs to.
- `server_uuid` (String) UUID of the server to provision the database on.

### Optional

- `instant_deploy` (Boolean) When `true`, Coolify immediately starts the database container after creation. Defaults to `false`.
- `name` (String) Display name for the database. If omitted, Coolify generates a name.
- `redis_password` (String, Sensitive) Password for Redis authentication. If omitted, the Redis instance will start without authentication.

### Read-Only

- `id` (String) Coolify resource UUID.

## Import

Import an existing Redis database using its UUID:

```shell
terraform import coolify_database_redis.cache <database-uuid>
```
