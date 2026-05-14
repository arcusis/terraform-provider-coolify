---
page_title: "coolify_database_backup Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Configures a scheduled backup for a Coolify-managed database. Supports local storage and S3-compatible object storage destinations.
---

# coolify_database_backup (Resource)

Configures a scheduled backup for a Coolify-managed database. Backups run on the schedule defined by `frequency` and can be stored locally on the server or uploaded to an S3-compatible storage bucket.

## Example Usage

### Daily backup with S3 upload

```hcl
resource "coolify_database_backup" "daily" {
  database_uuid  = coolify_database_postgresql.main.id
  frequency      = "daily"
  enabled        = true
  save_s3        = true
  s3_storage_uuid = var.s3_storage_uuid
}
```

### Custom cron schedule

```hcl
resource "coolify_database_backup" "nightly" {
  database_uuid = coolify_database_postgresql.main.id
  frequency     = "0 2 * * *"
  enabled       = true
  dump_all      = false
  timeout       = 300
}
```

## Schema

### Required

- `database_uuid` (String) UUID of the database to back up.
- `frequency` (String) Backup schedule. Accepts preset values (`daily`, `weekly`, `monthly`) or a standard cron expression (e.g. `0 2 * * *`).

### Optional

- `dump_all` (Boolean) When `true`, dumps all databases on the server in a single backup operation. Defaults to `false`.
- `enabled` (Boolean) Whether the backup schedule is active. Set to `false` to pause backups without deleting the configuration. Defaults to `true`.
- `s3_storage_uuid` (String) UUID of the S3 storage destination to upload backups to. Required when `save_s3` is `true`.
- `save_s3` (Boolean) When `true`, backup files are uploaded to the configured S3 storage destination. Defaults to `false`.
- `timeout` (Number) Maximum time in seconds to wait for the backup to complete before aborting.

### Read-Only

- `id` (String) Composite ID in the format `database_uuid/backup_uuid`.

## Import

Import an existing database backup configuration using its composite ID:

```shell
terraform import coolify_database_backup.daily <database-uuid>/<backup-uuid>
```
