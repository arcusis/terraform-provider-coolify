---
page_title: "coolify_application_deployments Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all deployments for a specific Coolify application.
---

# coolify_application_deployments (Data Source)

Lists all deployments for a specific Coolify application. Returns deployment history including status and commit information for each deployment.

## Example Usage

```hcl
data "coolify_application_deployments" "api" {
  application_uuid = coolify_application.api.id
}

output "latest_deployment_status" {
  value = data.coolify_application_deployments.api.deployments[0].status
}
```

## Schema

### Required

- `application_uuid` (String) UUID of the application to list deployments for.

### Read-Only

- `deployments` (Attributes List) List of deployments for the application, ordered most-recent first. (see [below for nested schema](#nestedatt--deployments))

<a id="nestedatt--deployments"></a>
### Nested Schema for `deployments`

Read-Only:

- `application_uuid` (String) UUID of the application this deployment belongs to.
- `commit_message` (String) Git commit message for the deployed commit.
- `status` (String) Current deployment status (e.g. `queued`, `in_progress`, `finished`, `failed`).
- `uuid` (String) UUID of the deployment.
