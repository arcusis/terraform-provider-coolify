---
page_title: "coolify_databases Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all databases (PostgreSQL, MySQL, Redis, MongoDB, etc.) managed by this Coolify instance.
---

# coolify_databases (Data Source)

Returns all databases across all projects and environments on the Coolify instance.

## Example Usage

```hcl
data "coolify_databases" "all" {}

# Find all PostgreSQL databases
locals {
  postgres_dbs = [
    for db in data.coolify_databases.all.databases :
    db if db.type == "postgresql"
  ]
}
```

## Schema

### Read-Only

- `databases` (List of Object) All databases on this Coolify instance.
  - `uuid` (String) Database UUID.
  - `name` (String) Database name.
  - `status` (String) Current status.
  - `type` (String) Database engine: `postgresql`, `mysql`, `mariadb`, `redis`, `mongodb`, `keydb`, `dragonfly`, or `clickhouse`.
  - `project_uuid` (String) UUID of the project.
  - `environment_name` (String) Name of the environment.
  - `server_uuid` (String) UUID of the server.
