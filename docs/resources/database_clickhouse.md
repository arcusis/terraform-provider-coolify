---
page_title: "coolify_database_clickhouse Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Provisions a ClickHouse container on a Coolify server within a project environment.
---

# coolify_database_clickhouse (Resource)

Provisions a ClickHouse container on a Coolify-managed server. ClickHouse is a high-performance columnar database suited for analytics workloads. The instance is created within the specified project and environment.

## Example Usage

```hcl
resource "coolify_database_clickhouse" "analytics" {
  project_uuid             = coolify_project.myapp.id
  server_uuid              = var.server_uuid
  environment_name         = "production"
  name                     = "analytics"
  clickhouse_admin_user    = "admin"
  clickhouse_admin_password = var.clickhouse_password
  instant_deploy           = true
}
```

## Schema

### Required

- `environment_name` (String) Name of the environment to provision the database in (e.g. `production`, `staging`).
- `project_uuid` (String) UUID of the project this database belongs to.
- `server_uuid` (String) UUID of the server to provision the database on.

### Optional

- `clickhouse_admin_password` (String, Sensitive) Password for the ClickHouse admin user. If omitted, Coolify generates a random password.
- `clickhouse_admin_user` (String) Admin username for ClickHouse. If omitted, Coolify generates a username.
- `instant_deploy` (Boolean) When `true`, Coolify immediately starts the database container after creation. Defaults to `false`.
- `name` (String) Display name for the database. If omitted, Coolify generates a name.

### Read-Only

- `id` (String) Coolify resource UUID.

## Import

Import an existing ClickHouse database using its UUID:

```shell
terraform import coolify_database_clickhouse.analytics <database-uuid>
```
