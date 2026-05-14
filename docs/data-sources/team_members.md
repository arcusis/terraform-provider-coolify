---
page_title: "coolify_team_members Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all members of a Coolify team, including their role (owner or member).
---

# coolify_team_members (Data Source)

Returns all members of a specific Coolify team. Use the team ID from `coolify_team` or `coolify_teams`.

## Example Usage

```hcl
data "coolify_team" "root" {
  name = "Root Team"
}

data "coolify_team_members" "root_members" {
  team_id = data.coolify_team.root.id
}

output "team_members" {
  value = [for m in data.coolify_team_members.root_members.members : m.email]
}
```

## Schema

### Required

- `team_id` (String) The team ID to list members for. Use the `id` output from `coolify_team`.

### Read-Only

- `members` (List of Object) All members of the team.
  - `id` (String) User ID.
  - `name` (String) User display name.
  - `email` (String) User email address.
  - `role` (String) Team role (`owner` or `member`).
