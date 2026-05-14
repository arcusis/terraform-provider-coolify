---
page_title: "coolify_projects Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all Coolify projects. Use this to enumerate projects and look up UUIDs for use in other resources.
---

# coolify_projects (Data Source)

Lists all projects available on the Coolify instance. Use this to enumerate projects without knowing their UUIDs ahead of time.

## Example Usage

```hcl
data "coolify_projects" "all" {}

output "project_names" {
  value = [for p in data.coolify_projects.all.projects : p.name]
}

# Reference a specific project by name
locals {
  prod_project = one([
    for p in data.coolify_projects.all.projects : p
    if p.name == "production"
  ])
}
```

## Schema

### Read-Only

- `projects` (List of Object) All projects on this Coolify instance.
  - `uuid` (String) Project UUID.
  - `name` (String) Project name.
  - `description` (String) Project description.
