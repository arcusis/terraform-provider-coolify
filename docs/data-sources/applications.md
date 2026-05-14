---
page_title: "coolify_applications Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all applications deployed on this Coolify instance, across all projects and environments.
---

# coolify_applications (Data Source)

Returns a list of all applications across all projects and environments on the Coolify instance.

## Example Usage

```hcl
data "coolify_applications" "all" {}

output "running_apps" {
  value = [
    for app in data.coolify_applications.all.applications :
    app.name if app.status == "running"
  ]
}
```

## Schema

### Read-Only

- `applications` (List of Object) All applications on this Coolify instance.
  - `uuid` (String) Application UUID.
  - `name` (String) Application name.
  - `status` (String) Current status (e.g. `running`, `stopped`, `exited`).
  - `project_uuid` (String) UUID of the project this application belongs to.
  - `environment_name` (String) Name of the environment.
  - `server_uuid` (String) UUID of the server hosting this application.
