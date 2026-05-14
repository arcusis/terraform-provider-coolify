---
page_title: "coolify_instance_settings Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Manages Coolify instance-wide settings such as auto-update, user registration, and the instance FQDN.
---

# coolify_instance_settings (Resource)

Reads and updates instance-wide Coolify settings via `GET /api/v1/settings` and `PATCH /api/v1/settings`.

## Example Usage

```hcl
resource "coolify_instance_settings" "main" {
  is_registration_enabled   = false
  is_auto_update_enabled    = true
  is_usage_tracking_enabled = false
  fqdn                      = "https://coolify.example.com"
}
```

## Schema

### Optional

- `fqdn` (String) Fully qualified domain name of this Coolify instance.
- `is_auto_update_enabled` (Boolean) Whether Coolify automatically updates itself.
- `is_registration_enabled` (Boolean) Whether new users can self-register.
- `is_usage_tracking_enabled` (Boolean) Whether anonymous usage telemetry is sent to the Coolify team.

### Read-Only

- `fqdn` (String) Populated from the API when not specified.
