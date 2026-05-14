---
page_title: "coolify_server_domains Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all domains configured on a specific Coolify server across all its hosted resources.
---

# coolify_server_domains (Data Source)

Lists all domains configured on a specific Coolify server, aggregated across all applications and services hosted on it. Useful for auditing which domains are in use or for detecting conflicts before adding new ones.

## Example Usage

```hcl
data "coolify_server_domains" "web" {
  server_uuid = var.server_uuid
}

output "all_domains" {
  value = [for d in data.coolify_server_domains.web.domains : d.domain]
}
```

## Schema

### Required

- `server_uuid` (String) UUID of the server to list domains for.

### Read-Only

- `domains` (Attributes List) List of domains configured on the server. (see [below for nested schema](#nestedatt--domains))

<a id="nestedatt--domains"></a>
### Nested Schema for `domains`

Read-Only:

- `domain` (String) Domain name (e.g. `app.example.com`).
