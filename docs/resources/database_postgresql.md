---
page_title: "coolify_database_postgresql Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Provisions and manages a PostgreSQL database container on a Coolify server.
---

# coolify_database_postgresql (Resource)

Provisions a PostgreSQL database container on a Coolify server. Coolify handles the container lifecycle, and you can optionally set up scheduled backups with `coolify_database_backup`.

## Example Usage

### Basic PostgreSQL Database

```hcl
resource "coolify_database_postgresql" "app_db" {
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "app-db"
  postgres_user    = "appuser"
  postgres_db      = "app"
  instant_deploy   = true
}
```

### With Daily Backups

```hcl
resource "coolify_database_postgresql" "app_db" {
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "app-db"
  postgres_user    = "appuser"
  postgres_db      = "app"
  instant_deploy   = true
}

resource "coolify_database_backup" "daily" {
  database_uuid = coolify_database_postgresql.app_db.id
  frequency     = "daily"
  enabled       = true
}
```

### Environment Variable on the Database Container

```hcl
resource "coolify_database_environment_variable" "timezone" {
  database_uuid = coolify_database_postgresql.app_db.id
  key           = "PGTZ"
  value         = "UTC"
}
```

## Schema

### Required

- `environment_name` (String) Name of the environment to deploy into (e.g. `production`, `staging`).
- `project_uuid` (String) UUID of the project.
- `server_uuid` (String) UUID of the server to run the database on.

### Optional

- `instant_deploy` (Boolean) Start the database container immediately after creation. Defaults to `false`.
- `name` (String) Display name shown in the Coolify dashboard.
- `postgres_db` (String) Name of the default database to create inside PostgreSQL.
- `postgres_password` (String, Sensitive) Password for the PostgreSQL superuser. Coolify generates a random password if omitted.
- `postgres_user` (String) Username for the PostgreSQL superuser.

### Read-Only

- `id` (String) Coolify resource UUID. Use this with `coolify_database_backup`, `coolify_database_environment_variable`, and `coolify_resource_action`.

## Import

```shell
terraform import coolify_database_postgresql.app_db <database-uuid>
```
