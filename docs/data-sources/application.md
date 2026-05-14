---
page_title: "coolify_application Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Looks up an existing Coolify application by UUID and returns its configuration and current status.
---

# coolify_application (Data Source)

Looks up an existing Coolify application by UUID and returns its configuration and current deployment status. Use this to reference applications created outside of Terraform or to inspect their current state.

## Example Usage

```hcl
data "coolify_application" "api" {
  uuid = var.application_uuid
}

output "app_status" {
  value = data.coolify_application.api.status
}

output "app_domains" {
  value = data.coolify_application.api.domains
}
```

## Schema

### Required

- `uuid` (String) UUID of the application to look up.

### Read-Only

- `build_pack` (String) Build system in use (e.g. `nixpacks`, `dockerfile`, `static`).
- `description` (String) Optional description of the application.
- `domains` (String) Comma-separated list of domains the application is exposed on.
- `git_branch` (String) Git branch the application is deployed from.
- `git_repository` (String) Git repository URL.
- `id` (String) Terraform data source ID (same as `uuid`).
- `name` (String) Display name of the application.
- `status` (String) Current status of the application (e.g. `running`, `stopped`, `exited`).
