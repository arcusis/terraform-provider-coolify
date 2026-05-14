---
page_title: "coolify_deploy Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Triggers a deployment for a Coolify application or service. The deployment is fired once on resource creation; destroy sends a cancellation request.
---

# coolify_deploy (Resource)

Triggers a deployment for a Coolify application or service. The deployment fires once when the resource is created. Destroying this resource sends a cancellation request for any in-progress deployment.

> **Note:** This resource is imperative. It triggers a deployment on `terraform apply` but does not track whether the application is currently deployed. To retrigger a deployment, you can use a `triggers` meta-argument or replace the resource.

## Example Usage

### Deploy on apply

```hcl
resource "coolify_deploy" "api" {
  resource_uuid = coolify_application.api.id
}

output "deployment_uuid" {
  value = coolify_deploy.api.deployment_uuid
}
```

### Force rebuild (bypass cache)

```hcl
resource "coolify_deploy" "api_rebuild" {
  resource_uuid = coolify_application.api.id
  force         = true
}
```

### Redeploy when config changes

```hcl
resource "coolify_deploy" "api" {
  resource_uuid = coolify_application.api.id

  triggers = {
    vars_hash = sha256(jsonencode(coolify_envs_bulk.app_vars.variables))
  }
}
```

## Schema

### Required

- `resource_uuid` (String) UUID of the application or service to deploy.

### Optional

- `force` (Boolean) When `true`, forces a rebuild even if nothing has changed, bypassing the build cache. Defaults to `false`.

### Read-Only

- `deployment_uuid` (String) UUID of the triggered deployment. Use with the `coolify_deployment` data source to track its status.
- `id` (String) Coolify resource UUID.
- `status` (String) Status of the deployment at the time of the last read (e.g. `queued`, `in_progress`, `finished`, `failed`).

## Import

Import an existing deploy resource using its deployment UUID:

```shell
terraform import coolify_deploy.api <deployment-uuid>
```
