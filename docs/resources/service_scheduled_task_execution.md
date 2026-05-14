---
page_title: "coolify_service_scheduled_task_execution Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Tracks and manages a single execution record for a service scheduled task. Destroying this resource deletes the execution log entry.
---

# coolify_service_scheduled_task_execution (Resource)

Tracks a scheduled task execution record for a service. Creating this resource records the execution UUID in Terraform state. Destroying it calls `DELETE /api/v1/services/{parent_uuid}/scheduled-tasks/{task_uuid}/executions/{execution_uuid}` to remove the log entry from Coolify.

A 404 response on delete is treated as success.

## Example Usage

```hcl
resource "coolify_service_scheduled_task_execution" "old_run" {
  parent_uuid    = coolify_service.svc.id
  task_uuid      = coolify_service_scheduled_task.maintenance.id
  execution_uuid = "abc12345-0000-0000-0000-000000000001"
}
```

## Schema

### Required

- `execution_uuid` (String) UUID of the specific execution record to manage.
- `parent_uuid` (String) UUID of the service that owns the scheduled task.
- `task_uuid` (String) UUID of the scheduled task.
