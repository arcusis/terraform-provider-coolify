---
page_title: "coolify_private_key Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Stores an SSH private key in Coolify. The key can be used to authenticate to remote servers or as a Git deploy key for private repositories.
---

# coolify_private_key (Resource)

Stores an SSH private key in Coolify. Once stored, the key can be referenced when registering servers or configuring private Git repository access. The private key material is write-only and is not returned by the API after creation.

## Example Usage

```hcl
resource "coolify_private_key" "deploy" {
  name        = "deploy-key"
  description = "ED25519 key for server authentication"
  private_key = var.ssh_private_key
}

resource "coolify_server" "web" {
  name             = "web-01"
  ip               = "203.0.113.10"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.deploy.id
}
```

## Schema

### Required

- `name` (String) Display name for the key in the Coolify dashboard.
- `private_key` (String, Sensitive) PEM-encoded SSH private key. This value is write-only — it is sent to the API on create/update but is not returned on read. Store this value in a secret manager or Terraform variable.

### Optional

- `description` (String) Optional description for the key.

### Read-Only

- `id` (String) Coolify resource UUID. Use this to reference the key in other resources.

## Import

Import an existing private key using its UUID:

```shell
terraform import coolify_private_key.deploy <key-uuid>
```
