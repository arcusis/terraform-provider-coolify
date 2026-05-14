---
page_title: "coolify_backup_executions Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists execution history for a scheduled database backup. Each record shows whether the backup succeeded and when it ran.
---

# coolify_backup_executions (Data Source)

Returns the execution history for a specific scheduled database backup schedule. Use this to audit backup runs or detect failures.

## Example Usage

```hcl
resource "coolify_database_backup" "daily" {
  database_uuid = coolify_database_postgresql.app_db.id
  frequency     = "daily"
  enabled       = true
}

data "coolify_backup_executions" "recent" {
  database_uuid        = coolify_database_postgresql.app_db.id
  scheduled_backup_uuid = coolify_database_backup.daily.id
}

output "last_backup_status" {
  value = length(data.coolify_backup_executions.recent.executions) > 0 ? (
    data.coolify_backup_executions.recent.executions[0].status
  ) : "no runs yet"
}
```

## Schema

### Required

- `database_uuid` (String) UUID of the database that owns the backup schedule.
- `scheduled_backup_uuid` (String) UUID of the scheduled backup configuration (from `coolify_database_backup`).

### Read-Only

- `executions` (List of Object) Backup run records.
  - `uuid` (String) Execution UUID.
  - `status` (String) Run status (`success`, `failed`, `running`).
  - `started_at` (String) ISO 8601 start timestamp.
  - `ended_at` (String) ISO 8601 end timestamp.
