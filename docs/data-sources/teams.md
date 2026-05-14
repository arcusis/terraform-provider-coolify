---
page_title: "coolify_teams Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all teams accessible to the current API token on this Coolify instance.
---

# coolify_teams (Data Source)

Returns all teams the current API token has access to. Most Coolify instances have a single "Root Team", but enterprise setups may have multiple.

## Example Usage

```hcl
data "coolify_teams" "all" {}

output "team_names" {
  value = [for t in data.coolify_teams.all.teams : t.name]
}
```

## Schema

### Read-Only

- `teams` (List of Object) All teams accessible to the current API token.
  - `id` (String) Team ID.
  - `name` (String) Team name.
  - `description` (String) Team description.
