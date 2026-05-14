---
page_title: "coolify_hetzner_locations Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists available Hetzner Cloud datacenter locations using the provided cloud token.
---

# coolify_hetzner_locations (Data Source)

Lists all available Hetzner Cloud datacenter locations using a Coolify cloud token. Use this data source to discover valid location identifiers when provisioning Hetzner servers.

## Example Usage

```hcl
resource "coolify_cloud_token" "hetzner" {
  name           = "hetzner"
  cloud_provider = "hetzner"
  token          = var.hetzner_api_token
}

data "coolify_hetzner_locations" "all" {
  cloud_token_uuid = coolify_cloud_token.hetzner.id
}

output "locations" {
  value = [for l in data.coolify_hetzner_locations.all.locations : "${l.name} (${l.city}, ${l.country})"]
}
```

## Schema

### Required

- `cloud_token_uuid` (String) UUID of the `coolify_cloud_token` with Hetzner credentials to use for the API call.

### Read-Only

- `locations` (Attributes List) List of available Hetzner datacenter locations. (see [below for nested schema](#nestedatt--locations))

<a id="nestedatt--locations"></a>
### Nested Schema for `locations`

Read-Only:

- `city` (String) City where the datacenter is located (e.g. `Nuremberg`).
- `country` (String) Country code where the datacenter is located (e.g. `DE`).
- `description` (String) Human-readable description of the location.
- `name` (String) Location identifier (e.g. `nbg1`, `fsn1`, `hel1`). Use this value when selecting a location for server provisioning.
