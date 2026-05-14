---
page_title: "coolify_cloud_token Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Stores a cloud provider API token in Coolify. Supported providers are Hetzner and DigitalOcean.
---

# coolify_cloud_token (Resource)

Stores a cloud provider API token in Coolify. The token is used to list available server types and datacenter locations when provisioning new servers through the Coolify UI or provider data sources. The token value is write-only and is not returned by the API after creation.

## Example Usage

```hcl
resource "coolify_cloud_token" "hetzner" {
  name           = "hetzner-token"
  cloud_provider = "hetzner"
  token          = var.hetzner_api_token
}

data "coolify_hetzner_locations" "all" {
  cloud_token_uuid = coolify_cloud_token.hetzner.id
}

data "coolify_hetzner_server_types" "all" {
  cloud_token_uuid = coolify_cloud_token.hetzner.id
}
```

### DigitalOcean

```hcl
resource "coolify_cloud_token" "digitalocean" {
  name           = "do-token"
  cloud_provider = "digitalocean"
  token          = var.digitalocean_api_token
}
```

## Schema

### Required

- `cloud_provider` (String) Cloud provider for this token. Accepted values: `hetzner`, `digitalocean`.
- `name` (String) Display name for the token in the Coolify dashboard.
- `token` (String, Sensitive) API token for the cloud provider. This value is write-only — it is sent to the API on create/update but is not returned on read.

### Read-Only

- `id` (String) Coolify resource UUID. Use this to reference the token in data sources.

## Import

Import an existing cloud token using its UUID:

```shell
terraform import coolify_cloud_token.hetzner <token-uuid>
```
