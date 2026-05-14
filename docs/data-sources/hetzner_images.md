---
page_title: "coolify_hetzner_images Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists available OS images in your Hetzner Cloud account. Requires a Hetzner cloud token stored in Coolify.
---

# coolify_hetzner_images (Data Source)

Lists available OS images from your Hetzner Cloud account. Requires a `coolify_cloud_token` with Hetzner credentials.

## Example Usage

```hcl
resource "coolify_cloud_token" "hetzner" {
  name           = "hetzner-prod"
  cloud_provider = "hetzner"
  token          = var.hetzner_api_token
}

data "coolify_hetzner_images" "available" {
  cloud_token_uuid = coolify_cloud_token.hetzner.id
}

# Find Ubuntu 22.04
locals {
  ubuntu_22 = one([
    for img in data.coolify_hetzner_images.available.images :
    img if img.name == "ubuntu-22.04"
  ])
}
```

## Schema

### Required

- `cloud_token_uuid` (String) UUID of the Hetzner cloud token stored in Coolify (from `coolify_cloud_token`).

### Read-Only

- `images` (List of Object) Available Hetzner OS images.
  - `id` (Number) Hetzner image ID.
  - `name` (String) Image name (e.g. `ubuntu-22.04`).
  - `description` (String) Human-readable description.
  - `os_family` (String) OS family (`ubuntu`, `debian`, `centos`, etc.).
  - `os_version` (String) OS version string.
