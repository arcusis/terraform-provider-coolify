---
page_title: "coolify_service Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Looks up an existing Coolify service by UUID and returns its type, name, and current status.
---

# coolify_service (Data Source)

Looks up an existing Coolify service stack by UUID and returns its configuration and current status. Use this to reference services created outside of Terraform.

## Example Usage

```hcl
data "coolify_service" "blog" {
  uuid = var.service_uuid
}

output "service_status" {
  value = data.coolify_service.blog.status
}
```

## Schema

### Required

- `uuid` (String) UUID of the service to look up.

### Read-Only

- `description` (String) Optional description of the service.
- `id` (String) Terraform data source ID (same as `uuid`).
- `name` (String) Display name of the service.
- `service_type` (String) The type of one-click service stack (e.g. `ghost`, `wordpress`, `plausible`).
- `status` (String) Current status of the service (e.g. `running`, `stopped`).
