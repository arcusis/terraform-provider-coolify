---
page_title: "coolify_server_resources Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all resources (applications, services, databases) deployed on a specific Coolify server.
---

# coolify_server_resources (Data Source)

Lists all resources (applications, services, and databases) deployed on a specific Coolify server. Useful for inventory, auditing, or building conditional logic based on what is running on a server.

## Example Usage

```hcl
data "coolify_server_resources" "web" {
  server_uuid = var.server_uuid
}

output "running_apps" {
  value = [for r in data.coolify_server_resources.web.resources : r.name if r.status == "running"]
}
```

## Schema

### Required

- `server_uuid` (String) UUID of the server to list resources for.

### Read-Only

- `resources` (Attributes List) List of all resources deployed on the server. (see [below for nested schema](#nestedatt--resources))

<a id="nestedatt--resources"></a>
### Nested Schema for `resources`

Read-Only:

- `name` (String) Display name of the resource.
- `project_uuid` (String) UUID of the project the resource belongs to.
- `status` (String) Current status of the resource (e.g. `running`, `stopped`, `exited`).
- `type` (String) Resource type (e.g. `application`, `service`, `database`).
- `uuid` (String) UUID of the resource.
