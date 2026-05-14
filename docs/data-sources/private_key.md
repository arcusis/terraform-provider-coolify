---
page_title: "coolify_private_key Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Looks up an existing Coolify private key by UUID and returns its name and description. The key material is never returned.
---

# coolify_private_key (Data Source)

Looks up an existing Coolify private key by UUID. Returns the key's name and description. The private key material is never returned by the API and is not available as an output.

## Example Usage

```hcl
data "coolify_private_key" "deploy" {
  uuid = var.deploy_key_uuid
}

resource "coolify_server" "web" {
  name             = "web-01"
  ip               = "203.0.113.10"
  port             = 22
  user             = "root"
  private_key_uuid = data.coolify_private_key.deploy.id
}
```

## Schema

### Required

- `uuid` (String) UUID of the private key to look up.

### Read-Only

- `description` (String) Optional description of the key.
- `id` (String) Terraform data source ID (same as `uuid`).
- `name` (String) Display name of the key.
