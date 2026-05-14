---
page_title: "coolify_resource_action Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Performs a one-time start, stop, or restart action on a Coolify application, service, or database. This is an imperative operation resource, not declarative state.
---

# coolify_resource_action (Resource)

Performs a one-time start, stop, or restart action on a Coolify application, service, or database. Unlike most Terraform resources, this resource is imperative: the action fires once when the resource is created (or when its inputs change), and on destroy Coolify receives a stop command.

> **Note:** This resource does not represent desired state. If you `terraform apply` with `action = "start"` and then something stops the resource outside of Terraform, a subsequent `terraform plan` will show no diff. Use this for bootstrap actions or one-shot operational tasks.

## Example Usage

### Start an application after provisioning

```hcl
resource "coolify_resource_action" "start_api" {
  resource_type = "application"
  resource_uuid = coolify_application.api.id
  action        = "start"

  depends_on = [coolify_envs_bulk.app_vars]
}
```

### Restart a service

```hcl
resource "coolify_resource_action" "restart_blog" {
  resource_type = "service"
  resource_uuid = coolify_service.blog.id
  action        = "restart"
}
```

## Schema

### Required

- `action` (String) Action to perform. Accepted values: `start`, `stop`, `restart`.
- `resource_type` (String) Type of the target resource. Accepted values: `application`, `service`, `database`.
- `resource_uuid` (String) UUID of the resource to act on.

### Read-Only

- `id` (String) Coolify resource UUID.
- `status` (String) Last reported status of the resource after the action was performed.

## Import

Import an existing resource action using its UUID:

```shell
terraform import coolify_resource_action.start_api <action-uuid>
```
