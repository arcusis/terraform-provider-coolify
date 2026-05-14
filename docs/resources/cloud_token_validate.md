---
page_title: "coolify_cloud_token_validate Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Validates a Coolify cloud token against the cloud provider API. Use as a precondition before provisioning cloud servers.
---

# coolify_cloud_token_validate (Resource)

Calls `POST /api/v1/cloud-tokens/{uuid}/validate` to verify a cloud token against the provider API (Hetzner, DigitalOcean). The result is surfaced in `status`.

## Example Usage

```hcl
resource "coolify_cloud_token" "hetzner" {
  name           = "hetzner"
  cloud_provider = "hetzner"
  token          = var.hetzner_api_token
}

resource "coolify_cloud_token_validate" "hetzner" {
  cloud_token_uuid = coolify_cloud_token.hetzner.id
}

output "token_valid" {
  value = coolify_cloud_token_validate.hetzner.status
}
```

## Schema

### Required

- `cloud_token_uuid` (String) UUID of the cloud token to validate.

### Read-Only

- `status` (String) Validation result from Coolify.
