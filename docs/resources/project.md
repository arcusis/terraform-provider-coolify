---
page_title: "coolify_project Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Creates and manages a Coolify project. Projects are the top-level grouping for applications, services, and databases.
---

# coolify_project (Resource)

Creates and manages a Coolify project. Projects are the top-level grouping for all resources (applications, services, databases) and contain one or more environments.

## Example Usage

```hcl
resource "coolify_project" "myapp" {
  name        = "my-application"
  description = "Production infrastructure for my-application"
}

resource "coolify_environment" "production" {
  project_uuid = coolify_project.myapp.id
  name         = "production"
}
```

## Schema

### Required

- `name` (String) Display name for the project in the Coolify dashboard.

### Optional

- `description` (String) Optional description for the project.

### Read-Only

- `id` (String) Coolify resource UUID. Use this to reference the project in other resources.

## Import

Import an existing project using its UUID:

```shell
terraform import coolify_project.myapp <project-uuid>
```
