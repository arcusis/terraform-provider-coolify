---
page_title: "coolify_backup_execution Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Tracks a database backup execution record. Destroying this resource purges the log entry via the Coolify API.
---

# coolify_backup_execution (Resource)

Tracks a single database backup execution record. Creating this resource stores the execution UUID in Terraform state. Destroying it calls `DELETE /api/v1/databases/{database_uuid}/backups/{scheduled_backup_uuid}/executions/{execution_uuid}` to remove the log entry.

A 404 on delete is treated as success.

## Example Usage

```hcl
resource "coolify_backup_execution" "old_run" {
  database_uuid        = coolify_database_postgresql.app_db.id
  scheduled_backup_uuid = coolify_database_backup.daily.id
  execution_uuid       = "abc12345-0000-0000-0000-000000000001"
}
```

## Schema

### Required

- `database_uuid` (String) UUID of the database that owns the backup schedule.
- `execution_uuid` (String) UUID of the specific execution record to manage.
- `scheduled_backup_uuid` (String) UUID of the scheduled backup configuration (from `coolify_database_backup`).
