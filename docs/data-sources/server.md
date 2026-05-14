---
page_title: "coolify_server Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Looks up an existing Coolify server by UUID and returns its connection details and status.
---

# coolify_server (Data Source)

Looks up an existing Coolify server by UUID and returns its connection details and current status. Use this to reference servers registered outside of Terraform.

## Example Usage

```hcl
data "coolify_server" "web" {
  uuid = var.server_uuid
}

output "server_ip" {
  value = data.coolify_server.web.ip
}
```

## Schema

### Required

- `uuid` (String) UUID of the server to look up.

### Read-Only

- `description` (String) Optional description of the server.
- `id` (String) Terraform data source ID (same as `uuid`).
- `ip` (String) IP address or hostname of the server.
- `is_reachable` (Boolean) Whether Coolify can reach the server via SSH.
- `is_usable` (Boolean) Whether the server is validated and ready to host workloads.
- `name` (String) Display name of the server.
- `port` (Number) SSH port number.
- `proxy_type` (String) Reverse proxy type configured on the server (e.g. `traefik`, `nginx`).
- `user` (String) SSH username used to connect to the server.
