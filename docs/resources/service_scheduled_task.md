---
page_title: "coolify_service_scheduled_task Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Creates a scheduled cron task that runs a command inside a service container on a defined schedule.
---

# coolify_service_scheduled_task (Resource)

Creates a scheduled cron task that runs a command inside a Coolify service container. Tasks are useful for periodic maintenance jobs, data processing, or any recurring work that should run within the service stack.

## Example Usage

### Nightly Ghost backup

```hcl
resource "coolify_service_scheduled_task" "backup" {
  service_uuid = coolify_service.blog.id
  name         = "nightly-backup"
  command      = "ghost backup"
  frequency    = "0 3 * * *"
  enabled      = true
}
```

### Weekly maintenance task

```hcl
resource "coolify_service_scheduled_task" "maintenance" {
  service_uuid = coolify_service.blog.id
  name         = "weekly-maintenance"
  command      = "/scripts/maintenance.sh"
  frequency    = "0 4 * * 0"
  container    = "ghost"
  enabled      = true
  timeout      = 300
}
```

## Schema

### Required

- `command` (String) Shell command to execute inside the container.
- `frequency` (String) Cron expression defining the schedule (e.g. `0 3 * * *` for 3 AM daily).
- `name` (String) Display name for the scheduled task.
- `service_uuid` (String) UUID of the service this scheduled task belongs to.

### Optional

- `container` (String) Name of the specific container within the service stack to run the command in. If omitted, the command runs in the primary container.
- `enabled` (Boolean) Whether the scheduled task is active. Set to `false` to pause without deleting. Defaults to `true`.
- `timeout` (Number) Maximum time in seconds to wait for the task to complete before aborting.

### Read-Only

- `id` (String) Composite ID in the format `service_uuid/task_uuid`.

## Import

Import an existing service scheduled task using its composite ID:

```shell
terraform import coolify_service_scheduled_task.backup <service-uuid>/<task-uuid>
```
