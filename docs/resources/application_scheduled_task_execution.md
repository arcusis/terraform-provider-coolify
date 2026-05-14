---
page_title: "coolify_application_scheduled_task_execution Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Tracks and manages a single execution record for an application scheduled task. Destroying this resource deletes the execution log entry.
---

# coolify_application_scheduled_task_execution (Resource)

Tracks a scheduled task execution record for an application. Creating this resource records the execution UUID in Terraform state. Destroying it calls `DELETE /api/v1/applications/{parent_uuid}/scheduled-tasks/{task_uuid}/executions/{execution_uuid}` to remove the log entry from Coolify.

A 404 response on delete is treated as success — the record is already gone.

## Example Usage

```hcl
# Clean up a specific execution record
resource "coolify_application_scheduled_task_execution" "old_run" {
  parent_uuid    = coolify_application.app.id
  task_uuid      = coolify_application_scheduled_task.cleanup.id
  execution_uuid = "abc12345-0000-0000-0000-000000000001"
}
```

## Schema

### Required

- `execution_uuid` (String) UUID of the specific execution record to manage.
- `parent_uuid` (String) UUID of the application that owns the scheduled task.
- `task_uuid` (String) UUID of the scheduled task.
