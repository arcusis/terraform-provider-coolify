---
page_title: "coolify_application_scheduled_task Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Creates a scheduled cron task that runs a command inside an application container on a defined schedule.
---

# coolify_application_scheduled_task (Resource)

Creates a scheduled cron task that runs a command inside an application container. Tasks are useful for periodic jobs such as database migrations, cache clearing, report generation, or data imports.

## Example Usage

### Daily database cleanup

```hcl
resource "coolify_application_scheduled_task" "cleanup" {
  application_uuid = coolify_application.api.id
  name             = "cleanup"
  command          = "node scripts/cleanup.js"
  frequency        = "0 2 * * *"
  enabled          = true
}
```

### Weekly report with specific container and timeout

```hcl
resource "coolify_application_scheduled_task" "weekly_report" {
  application_uuid = coolify_application.api.id
  name             = "weekly-report"
  command          = "php artisan report:generate"
  frequency        = "0 9 * * 1"
  container        = "api"
  enabled          = true
  timeout          = 600
}
```

## Schema

### Required

- `application_uuid` (String) UUID of the application this scheduled task belongs to.
- `command` (String) Shell command to execute inside the container (e.g. `node scripts/cleanup.js`).
- `frequency` (String) Cron expression defining the schedule (e.g. `0 2 * * *` for 2 AM daily).
- `name` (String) Display name for the scheduled task.

### Optional

- `container` (String) Name of the specific container to run the command in, for multi-container applications. If omitted, the command runs in the primary container.
- `enabled` (Boolean) Whether the scheduled task is active. Set to `false` to pause without deleting. Defaults to `true`.
- `timeout` (Number) Maximum time in seconds to wait for the task to complete before aborting.

### Read-Only

- `id` (String) Composite ID in the format `application_uuid/task_uuid`.

## Import

Import an existing scheduled task using its composite ID:

```shell
terraform import coolify_application_scheduled_task.cleanup <application-uuid>/<task-uuid>
```
