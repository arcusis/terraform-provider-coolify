---
page_title: "coolify_team Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Looks up a Coolify team by ID or name and returns its metadata.
---

# coolify_team (Data Source)

Looks up a Coolify team by ID or name. At least one of `id` or `name` must be provided. If both are given, `id` takes precedence.

## Example Usage

### Look up by name

```hcl
data "coolify_team" "default" {
  name = "default"
}

output "team_id" {
  value = data.coolify_team.default.id
}
```

### Look up by ID

```hcl
data "coolify_team" "infra" {
  id = "1"
}
```

## Schema

### Optional

- `id` (String) Numeric team ID. Used to look up the team directly.
- `name` (String) Team name. Used for lookup when `id` is omitted.

### Read-Only

- `description` (String) Description of the team.
