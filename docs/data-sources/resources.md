---
page_title: "coolify_resources Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all resources (applications, services, databases) across all projects on the Coolify instance.
---

# coolify_resources (Data Source)

Lists all resources (applications, services, and databases) across all projects on the Coolify instance. Use `coolify_server_resources` instead if you need resources scoped to a specific server.

## Example Usage

```hcl
data "coolify_resources" "all" {}

output "resource_count" {
  value = length(data.coolify_resources.all.resources)
}

output "stopped_resources" {
  value = [for r in data.coolify_resources.all.resources : r.name if r.status == "stopped"]
}
```

## Schema

### Read-Only

- `resources` (Attributes List) List of all resources across all projects. (see [below for nested schema](#nestedatt--resources))

<a id="nestedatt--resources"></a>
### Nested Schema for `resources`

Read-Only:

- `name` (String) Display name of the resource.
- `project_uuid` (String) UUID of the project the resource belongs to.
- `status` (String) Current status of the resource (e.g. `running`, `stopped`, `exited`).
- `type` (String) Resource type (e.g. `application`, `service`, `database`).
- `uuid` (String) UUID of the resource.
