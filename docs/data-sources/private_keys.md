---
page_title: "coolify_private_keys Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all SSH private keys stored in Coolify. Note: the actual key material is never returned by the API.
---

# coolify_private_keys (Data Source)

Lists all SSH private keys stored in Coolify. This is useful for finding a key UUID to reference when creating servers or applications. The private key material itself is never returned by the API.

## Example Usage

```hcl
data "coolify_private_keys" "all" {}

# Find a specific key by name
locals {
  deploy_key = one([
    for k in data.coolify_private_keys.all.private_keys :
    k if k.name == "deploy-key"
  ])
}
```

## Schema

### Read-Only

- `private_keys` (List of Object) All private keys stored in this Coolify instance.
  - `uuid` (String) Private key UUID.
  - `name` (String) Key name.
  - `description` (String) Key description.
