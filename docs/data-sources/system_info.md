---
page_title: "coolify_system_info Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Reads Coolify instance health and version information. Useful for verifying provider connectivity.
---

# coolify_system_info (Data Source)

Reads the health status and version of the Coolify instance. This data source takes no inputs and is useful for verifying provider connectivity or gating resources on instance availability.

## Example Usage

```hcl
data "coolify_system_info" "info" {}

output "coolify_version" {
  value = data.coolify_system_info.info.version
}

output "coolify_healthy" {
  value = data.coolify_system_info.info.healthy
}
```

## Schema

### Read-Only

- `healthy` (Boolean) Whether the Coolify instance is healthy and fully operational.
- `version` (String) Current version of the Coolify instance (e.g. `v4.0.0-beta.123`).
