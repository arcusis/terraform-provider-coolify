---
page_title: "coolify_server Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Registers an SSH server with Coolify so it can be used as a deployment target for applications, services, and databases.
---

# coolify_server (Resource)

Registers a remote server with Coolify over SSH. Once registered, the server can be used as a deployment target for applications, services, and databases.

## Example Usage

```hcl
resource "coolify_private_key" "server_key" {
  name        = "my-server-key"
  private_key = var.ssh_private_key
}

resource "coolify_server" "web" {
  name             = "web-01"
  ip               = "203.0.113.10"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.server_key.id
  description      = "Primary web server in EU"
  instant_validate = true
}
```

### Build Server

```hcl
resource "coolify_server" "builder" {
  name             = "builder-01"
  ip               = "203.0.113.20"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.server_key.id
  is_build_server  = true
}
```

## Schema

### Required

- `ip` (String) IP address or hostname of the server.
- `name` (String) Display name for the server in the Coolify dashboard.
- `port` (Number) SSH port number (typically `22`).
- `private_key_uuid` (String) UUID of the `coolify_private_key` used to authenticate via SSH.
- `user` (String) SSH username (e.g. `root`).

### Optional

- `description` (String) Optional description shown in the Coolify dashboard.
- `instant_validate` (Boolean) When `true`, Coolify immediately tests the SSH connection after the server is registered. Defaults to `false`.
- `is_build_server` (Boolean) When `true`, this server is designated as a build-only server and will not run deployed workloads directly. Defaults to `false`.
- `proxy_type` (String) Reverse proxy type to use on the server. Common values: `traefik`, `nginx`, `caddy`, `none`.

### Read-Only

- `id` (String) Coolify resource UUID. Use this to reference the server in other resources.

## Import

Import an existing server using its UUID:

```shell
terraform import coolify_server.web <server-uuid>
```
