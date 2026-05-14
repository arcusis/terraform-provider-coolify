---
page_title: "coolify_team_member Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Adds a user to a Coolify team. Destroying removes the user from the team.
---

# coolify_team_member (Resource)

Adds a user to a Coolify team via `POST /api/v1/teams/{id}/members`. Destroying calls `DELETE /api/v1/teams/{id}/members/{user_id}`.

## Example Usage

```hcl
resource "coolify_team_member" "alice" {
  team_id = coolify_team.engineering.id
  email   = "alice@example.com"
  role    = "member"
}
```

## Schema

### Required

- `team_id` (String) ID of the team.

### Optional

- `email` (String) Email of the user to add.
- `role` (String) Role within the team (`owner` or `member`).
- `user_id` (String) ID of the user. Computed if only `email` is provided.

### Read-Only

- `user_id` (String) The resolved user ID.
