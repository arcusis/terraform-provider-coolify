---
page_title: "coolify_team Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Creates and manages a Coolify team via POST, PATCH, and DELETE /api/v1/teams.
---

# coolify_team (Resource)

Creates a Coolify team. Teams group users and projects together. Use `coolify_team_member` to add users to the team after creation.

## Example Usage

```hcl
resource "coolify_team" "engineering" {
  name        = "engineering"
  description = "Engineering team"
}

resource "coolify_team_member" "alice" {
  team_id = coolify_team.engineering.id
  email   = "alice@example.com"
  role    = "member"
}
```

## Schema

### Required

- `name` (String) Team name.

### Optional

- `description` (String) Team description.

### Read-Only

- `id` (String) Team ID.

## Import

```shell
terraform import coolify_team.engineering <team-id>
```
