---
page_title: "coolify_database Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Looks up any Coolify-managed database by UUID regardless of type, returning its type, name, and current status.
---

# coolify_database (Data Source)

Looks up any Coolify-managed database by UUID regardless of its type (PostgreSQL, MySQL, Redis, MongoDB, etc.). Use this data source when you know a database UUID but do not need type-specific configuration fields.

## Example Usage

```hcl
data "coolify_database" "main" {
  uuid = var.database_uuid
}

output "db_status" {
  value = data.coolify_database.main.status
}

output "db_type" {
  value = data.coolify_database.main.db_type
}
```

## Schema

### Required

- `uuid` (String) UUID of the database to look up.

### Read-Only

- `db_type` (String) Database engine type (e.g. `postgresql`, `mysql`, `mariadb`, `redis`, `mongodb`, `clickhouse`, `keydb`, `dragonfly`).
- `description` (String) Optional description of the database.
- `id` (String) Terraform data source ID (same as `uuid`).
- `name` (String) Display name of the database.
- `status` (String) Current status of the database container (e.g. `running`, `stopped`).
