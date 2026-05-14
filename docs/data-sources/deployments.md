---
page_title: "coolify_deployments Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all deployments across all applications on the Coolify instance.
---

# coolify_deployments (Data Source)

Lists all deployments across all applications on the Coolify instance. Use `coolify_application_deployments` instead if you need deployments for a specific application.

## Example Usage

```hcl
data "coolify_deployments" "all" {}

output "deployment_count" {
  value = length(data.coolify_deployments.all.deployments)
}
```

## Schema

### Read-Only

- `deployments` (Attributes List) List of all deployments. (see [below for nested schema](#nestedatt--deployments))

<a id="nestedatt--deployments"></a>
### Nested Schema for `deployments`

Read-Only:

- `application_uuid` (String) UUID of the application this deployment belongs to.
- `commit_message` (String) Git commit message for the deployed commit.
- `status` (String) Current deployment status (e.g. `queued`, `in_progress`, `finished`, `failed`).
- `uuid` (String) UUID of the deployment.
