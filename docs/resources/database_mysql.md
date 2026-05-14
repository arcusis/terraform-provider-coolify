---
page_title: "coolify_database_mysql Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Provisions a MySQL container on a Coolify server within a project environment.
---

# coolify_database_mysql (Resource)

Provisions a MySQL container on a Coolify-managed server. The database is created within the specified project and environment.

## Example Usage

```hcl
resource "coolify_database_mysql" "app_db" {
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "app-db"
  mysql_database   = "appdb"
  mysql_user       = "appuser"
  mysql_password   = var.db_password
  mysql_root_password = var.db_root_password
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
- `mysql_database` (String) Name of the initial database to create inside MySQL.
- `mysql_password` (String, Sensitive) Password for `mysql_user`. If omitted, Coolify generates a random password.
- `mysql_root_password` (String, Sensitive) Password for the MySQL root user. If omitted, Coolify generates a random password.
- `mysql_user` (String) Non-root MySQL username to create. If omitted, Coolify generates a username.
- `name` (String) Display name for the database. If omitted, Coolify generates a name.

### Read-Only

- `id` (String) Coolify resource UUID.

## Import

Import an existing MySQL database using its UUID:

```shell
terraform import coolify_database_mysql.app_db <database-uuid>
```
