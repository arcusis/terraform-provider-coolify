---
page_title: "coolify_project Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Looks up an existing Coolify project by UUID and returns its name, description, and environments.
---

# coolify_project (Data Source)

Looks up an existing Coolify project by UUID. Use this to reference projects that were created outside of Terraform.

## Example Usage

```hcl
data "coolify_project" "existing" {
  uuid = "abc123de-f456-7890-abcd-ef1234567890"
}

output "project_name" {
  value = data.coolify_project.existing.name
}
```

## Schema

### Required

- `uuid` (String) UUID of the project to look up.

### Read-Only

- `description` (String) Description of the project.
- `id` (String) Terraform data source ID (same as `uuid`).
- `name` (String) Display name of the project.
