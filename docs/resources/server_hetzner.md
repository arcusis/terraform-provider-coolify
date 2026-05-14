---
page_title: "coolify_server_hetzner Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Provisions a new Hetzner Cloud server and registers it with Coolify in a single step. Requires a cloud token with Hetzner credentials.
---

# coolify_server_hetzner (Resource)

Provisions a Hetzner Cloud server and registers it with Coolify via `POST /api/v1/servers/hetzner`. This is a combined operation — Coolify calls the Hetzner API to create the VM, installs its agent, and adds the server to your Coolify instance.

**Note:** Destroying this resource removes the server registration from Coolify but does **not** terminate the Hetzner VM. You must delete the VM separately in the Hetzner Cloud console or via the Hetzner provider.

## Example Usage

```hcl
resource "coolify_cloud_token" "hetzner" {
  name           = "hetzner-prod"
  cloud_provider = "hetzner"
  token          = var.hetzner_api_token
}

resource "coolify_private_key" "deploy" {
  name        = "deploy-key"
  private_key = var.ssh_private_key
}

# Look up available server types and locations
data "coolify_hetzner_server_types" "types" {
  cloud_token_uuid = coolify_cloud_token.hetzner.id
}

data "coolify_hetzner_locations" "locations" {
  cloud_token_uuid = coolify_cloud_token.hetzner.id
}

# Provision a Hetzner server through Coolify
resource "coolify_server_hetzner" "prod" {
  name             = "prod-1"
  cloud_token_uuid = coolify_cloud_token.hetzner.id
  server_type      = "cpx21"
  location         = "nbg1"
  image            = "ubuntu-22.04"
  private_key_uuid = coolify_private_key.deploy.id
  description      = "Production application server"
}

output "prod_server_ip" {
  value = coolify_server_hetzner.prod.ip_address
}
```

## Schema

### Required

- `cloud_token_uuid` (String) UUID of the `coolify_cloud_token` containing Hetzner API credentials.
- `image` (String) Hetzner OS image name to install (e.g. `ubuntu-22.04`). Use `coolify_hetzner_images` to list available images.
- `location` (String) Hetzner datacenter location (e.g. `nbg1`, `fsn1`, `hel1`). Use `coolify_hetzner_locations` to list available locations.
- `name` (String) Name for the server in both Coolify and Hetzner.
- `private_key_uuid` (String) UUID of the `coolify_private_key` to install on the new server for SSH access.
- `server_type` (String) Hetzner server type (e.g. `cpx11`, `cpx21`, `cx22`). Use `coolify_hetzner_server_types` to list available types.

### Optional

- `description` (String) Optional description shown in the Coolify dashboard.

### Read-Only

- `id` (String) Coolify server UUID.
- `ip_address` (String) Public IP address of the provisioned Hetzner server.

## Import

```shell
terraform import coolify_server_hetzner.prod <server-uuid>
```
