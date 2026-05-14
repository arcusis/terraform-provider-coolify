---
page_title: "coolify_deployment Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Looks up a single Coolify deployment by UUID and returns its status, commit information, and associated application.
---

# coolify_deployment (Data Source)

Looks up a single Coolify deployment by UUID and returns its status and commit information. Useful for checking the outcome of a specific deployment triggered by `coolify_deploy`.

## Example Usage

```hcl
resource "coolify_deploy" "api" {
  resource_uuid = coolify_application.api.id
}

data "coolify_deployment" "latest" {
  uuid = coolify_deploy.api.deployment_uuid
}

output "deploy_status" {
  value = data.coolify_deployment.latest.status
}
```

## Schema

### Required

- `uuid` (String) UUID of the deployment to look up.

### Read-Only

- `application_uuid` (String) UUID of the application this deployment belongs to.
- `commit_hash` (String) Git commit SHA that was deployed.
- `commit_message` (String) Git commit message for the deployed commit.
- `id` (String) Terraform data source ID (same as `uuid`).
- `status` (String) Current deployment status (e.g. `queued`, `in_progress`, `finished`, `failed`, `cancelled`).
