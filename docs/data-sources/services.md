---
page_title: "coolify_services Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all service stacks (Ghost, WordPress, Plausible, etc.) deployed on this Coolify instance.
---

# coolify_services (Data Source)

Lists all one-click service stacks deployed across all projects and environments.

## Example Usage

```hcl
data "coolify_services" "all" {}

output "service_uuids" {
  value = { for s in data.coolify_services.all.services : s.name => s.uuid }
}
```

## Schema

### Read-Only

- `services` (List of Object) All services on this Coolify instance.
  - `uuid` (String) Service UUID.
  - `name` (String) Service name.
  - `status` (String) Current status.
  - `project_uuid` (String) UUID of the project.
  - `environment_name` (String) Name of the environment.
  - `server_uuid` (String) UUID of the server.
