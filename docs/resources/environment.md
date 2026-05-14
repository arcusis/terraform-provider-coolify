---
page_title: "coolify_environment Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Creates and manages an environment inside a Coolify project. Environments allow you to group resources by deployment stage, such as production, staging, or development.
---

# coolify_environment (Resource)

Creates and manages an environment inside a Coolify project. Environments are named slots within a project that hold applications, services, and databases for a particular deployment stage.

## Example Usage

```hcl
resource "coolify_project" "myapp" {
  name = "my-application"
}

resource "coolify_environment" "production" {
  project_uuid = coolify_project.myapp.id
  name         = "production"
}

resource "coolify_environment" "staging" {
  project_uuid = coolify_project.myapp.id
  name         = "staging"
}
```

## Schema

### Required

- `name` (String) Name of the environment (e.g. `production`, `staging`, `development`).
- `project_uuid` (String) UUID of the project this environment belongs to.

### Read-Only

- `id` (String) Composite ID in the format `project_uuid/environment_name`.

## Import

Import an existing environment using its composite ID:

```shell
terraform import coolify_environment.production <project-uuid>/<environment-name>
```
