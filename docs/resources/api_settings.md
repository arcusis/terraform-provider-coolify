---
page_title: "coolify_api_settings Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Enables or disables the Coolify REST API. Setting enabled = false prevents all API access until re-enabled.
---

# coolify_api_settings (Resource)

Controls whether the Coolify REST API is enabled. Setting `enabled = false` calls `GET /api/v1/disable`, which blocks all future API requests until the API is re-enabled.

**Warning:** Disabling the API will prevent Terraform from making any further changes to Coolify until the API is re-enabled through the Coolify dashboard or by setting `enabled = true` and applying again.

This resource restores `enabled = true` on `terraform destroy` to avoid locking yourself out.

## Example Usage

```hcl
# Keep the API enabled (default operational state)
resource "coolify_api_settings" "main" {
  enabled = true
}
```

```hcl
# Temporarily disable the API (maintenance mode)
resource "coolify_api_settings" "maintenance" {
  enabled = false
}
```

## Schema

### Required

- `enabled` (Boolean) Whether the Coolify API is enabled. Set to `false` to call `/api/v1/disable`; set to `true` to call `/api/v1/enable`.
