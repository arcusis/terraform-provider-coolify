---
page_title: "coolify_hetzner_server_types Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists available Hetzner Cloud server types (instance sizes) using the provided cloud token.
---

# coolify_hetzner_server_types (Data Source)

Lists all available Hetzner Cloud server types (instance sizes) using a Coolify cloud token. Use this data source to discover valid server type identifiers and their resource specifications before provisioning a Hetzner server.

## Example Usage

```hcl
resource "coolify_cloud_token" "hetzner" {
  name           = "hetzner"
  cloud_provider = "hetzner"
  token          = var.hetzner_api_token
}

data "coolify_hetzner_server_types" "all" {
  cloud_token_uuid = coolify_cloud_token.hetzner.id
}

output "server_types" {
  value = [for t in data.coolify_hetzner_server_types.all.server_types : "${t.name}: ${t.cores} vCPU, ${t.memory} GB RAM"]
}
```

## Schema

### Required

- `cloud_token_uuid` (String) UUID of the `coolify_cloud_token` with Hetzner credentials to use for the API call.

### Read-Only

- `server_types` (Attributes List) List of available Hetzner server types. (see [below for nested schema](#nestedatt--server_types))

<a id="nestedatt--server_types"></a>
### Nested Schema for `server_types`

Read-Only:

- `cores` (Number) Number of vCPU cores for this server type.
- `description` (String) Human-readable description of the server type.
- `memory` (Number) Amount of RAM in GB for this server type.
- `name` (String) Server type identifier (e.g. `cx22`, `cx32`, `cpx41`). Use this value when selecting a server type for provisioning.
