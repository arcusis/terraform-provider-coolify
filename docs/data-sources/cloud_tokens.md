---
page_title: "coolify_cloud_tokens Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all cloud provider API tokens (Hetzner, DigitalOcean, etc.) stored in Coolify. Use their UUIDs with the Hetzner data sources.
---

# coolify_cloud_tokens (Data Source)

Lists all cloud provider tokens stored in Coolify. Use the token UUIDs with `coolify_hetzner_locations`, `coolify_hetzner_server_types`, `coolify_hetzner_images`, and `coolify_hetzner_ssh_keys`.

## Example Usage

```hcl
data "coolify_cloud_tokens" "all" {}

locals {
  hetzner_token = one([
    for t in data.coolify_cloud_tokens.all.cloud_tokens :
    t if t.cloud_provider == "hetzner"
  ])
}

data "coolify_hetzner_locations" "regions" {
  cloud_token_uuid = local.hetzner_token.uuid
}
```

## Schema

### Read-Only

- `cloud_tokens` (List of Object) All cloud provider tokens.
  - `uuid` (String) Token UUID.
  - `name` (String) Token name.
  - `cloud_provider` (String) Provider name (`hetzner`, `digitalocean`).
