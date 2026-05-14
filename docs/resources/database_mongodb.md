---
page_title: "coolify_database_mongodb Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Provisions a MongoDB container on a Coolify server within a project environment.
---

# coolify_database_mongodb (Resource)

Provisions a MongoDB container on a Coolify-managed server. The instance is created within the specified project and environment.

## Example Usage

```hcl
resource "coolify_database_mongodb" "app_db" {
  project_uuid                = coolify_project.myapp.id
  server_uuid                 = var.server_uuid
  environment_name            = "production"
  name                        = "app-db"
  mongo_initdb_root_username  = "admin"
  mongo_initdb_root_password  = var.mongo_root_password
  mongo_initdb_database       = "appdb"
  instant_deploy              = true
}
```

## Schema

### Required

- `environment_name` (String) Name of the environment to provision the database in (e.g. `production`, `staging`).
- `project_uuid` (String) UUID of the project this database belongs to.
- `server_uuid` (String) UUID of the server to provision the database on.

### Optional

- `instant_deploy` (Boolean) When `true`, Coolify immediately starts the database container after creation. Defaults to `false`.
- `mongo_initdb_database` (String) Name of the initial database to create inside MongoDB.
- `mongo_initdb_root_password` (String, Sensitive) Password for the MongoDB root user. If omitted, Coolify generates a random password.
- `mongo_initdb_root_username` (String) Username for the MongoDB root user. If omitted, Coolify generates a username.
- `name` (String) Display name for the database. If omitted, Coolify generates a name.

### Read-Only

- `id` (String) Coolify resource UUID.

## Import

Import an existing MongoDB database using its UUID:

```shell
terraform import coolify_database_mongodb.app_db <database-uuid>
```
