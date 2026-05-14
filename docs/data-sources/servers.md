---
page_title: "coolify_servers Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all servers registered with this Coolify instance. Use this to look up server UUIDs without hardcoding them.
---

# coolify_servers (Data Source)

Lists all servers registered with Coolify. Useful for discovering server UUIDs when you can't look them up from the dashboard.

## Example Usage

```hcl
data "coolify_servers" "all" {}

# Find a server by name
locals {
  prod_server = one([
    for s in data.coolify_servers.all.servers :
    s if s.name == "prod-1"
  ])
}

output "prod_server_uuid" {
  value = local.prod_server.uuid
}
```

## Schema

### Read-Only

- `servers` (List of Object) All servers registered with this Coolify instance.
  - `uuid` (String) Server UUID.
  - `name` (String) Server name.
  - `ip` (String) Server IP address or hostname.
  - `status` (String) Connection status (`reachable`, `unreachable`).
  - `description` (String) Optional description.
