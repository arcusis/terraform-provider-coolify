---
page_title: "coolify_profile Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Manages the display name and email of the user associated with the current API token.
---

# coolify_profile (Resource)

Updates the profile of the API token's user via `PATCH /api/v1/profile`.

## Example Usage

```hcl
data "coolify_profile" "me" {}

resource "coolify_profile" "me" {
  name = "Infrastructure Bot"
}
```

## Schema

### Optional

- `email` (String) User email address.
- `name` (String) User display name.

### Read-Only

- `id` (String) User ID.
