---
page_title: "coolify_application_preview Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Manages a pull request preview deployment for a Coolify application. Creating the resource starts the preview; destroying it triggers preview cleanup.
---

# coolify_application_preview (Resource)

Manages a pull request preview deployment for a Coolify application. When created, Coolify spins up a preview environment for the specified pull request. Destroying the resource calls the Coolify preview cleanup endpoint to tear down the preview.

This resource is useful for automating PR preview lifecycle management from a CI/CD pipeline that also runs Terraform.

## Example Usage

```hcl
variable "pull_request_id" {
  type        = number
  description = "GitHub pull request number to create a preview for"
}

resource "coolify_application_preview" "pr_preview" {
  application_uuid = coolify_application.api.id
  pull_request_id  = var.pull_request_id
}
```

### Conditional preview (only for PRs targeting main)

```hcl
resource "coolify_application_preview" "pr_preview" {
  count            = var.base_branch == "main" ? 1 : 0
  application_uuid = coolify_application.api.id
  pull_request_id  = var.pull_request_id
}
```

## Schema

### Required

- `application_uuid` (String) UUID of the application to create the preview deployment for.
- `pull_request_id` (Number) GitHub pull request number whose preview deployment to manage.

### Read-Only

- `id` (String) Coolify resource UUID for the preview deployment.

## Import

Import an existing PR preview using its UUID:

```shell
terraform import coolify_application_preview.pr_preview <preview-uuid>
```
