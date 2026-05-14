---
page_title: "coolify_hetzner_ssh_keys Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists SSH keys uploaded to your Hetzner Cloud account. Use the key IDs when provisioning Hetzner servers.
---

# coolify_hetzner_ssh_keys (Data Source)

Lists SSH keys available in your Hetzner Cloud account. Requires a `coolify_cloud_token` with Hetzner credentials.

## Example Usage

```hcl
resource "coolify_cloud_token" "hetzner" {
  name           = "hetzner-prod"
  cloud_provider = "hetzner"
  token          = var.hetzner_api_token
}

data "coolify_hetzner_ssh_keys" "available" {
  cloud_token_uuid = coolify_cloud_token.hetzner.id
}

output "ssh_key_names" {
  value = [for k in data.coolify_hetzner_ssh_keys.available.ssh_keys : k.name]
}
```

## Schema

### Required

- `cloud_token_uuid` (String) UUID of the Hetzner cloud token stored in Coolify.

### Read-Only

- `ssh_keys` (List of Object) SSH keys available in the Hetzner account.
  - `id` (Number) Hetzner SSH key ID.
  - `name` (String) Key name.
  - `fingerprint` (String) MD5 fingerprint of the public key.
