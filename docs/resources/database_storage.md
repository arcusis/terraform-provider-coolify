---
page_title: "coolify_database_storage Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Mounts a persistent volume on a Coolify-managed database container, ensuring data survives container restarts.
---

# coolify_database_storage (Resource)

Mounts a persistent volume on a Coolify-managed database container. Use this resource to ensure that database data directories are persisted across container restarts and redeployments.

## Example Usage

### Named Docker volume

```hcl
resource "coolify_database_storage" "data" {
  database_uuid = coolify_database_postgresql.main.id
  type          = "volume"
  mount_path    = "/var/lib/postgresql/data"
  name          = "pg-data"
}
```

### Bind mount from host path

```hcl
resource "coolify_database_storage" "data" {
  database_uuid = coolify_database_mysql.app_db.id
  type          = "bind"
  mount_path    = "/var/lib/mysql"
  host_path     = "/data/mysql"
  is_readonly   = false
}
```

## Schema

### Required

- `database_uuid` (String) UUID of the database to attach the volume to.
- `mount_path` (String) Absolute path inside the container where the volume is mounted (e.g. `/var/lib/postgresql/data`).
- `type` (String) Storage type. Use `volume` for a named Docker volume or `bind` for a host bind mount.

### Optional

- `host_path` (String) Absolute path on the host to bind-mount into the container. Only used when `type` is `bind`.
- `is_readonly` (Boolean) When `true`, the volume is mounted read-only inside the container. Defaults to `false`.
- `name` (String) Name for the Docker volume. If omitted, Coolify generates a name.

### Read-Only

- `id` (String) Composite ID in the format `database_uuid/storage_uuid`.

## Import

Import an existing database storage volume using its composite ID:

```shell
terraform import coolify_database_storage.data <database-uuid>/<storage-uuid>
```
