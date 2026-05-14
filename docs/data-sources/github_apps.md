---
page_title: "coolify_github_apps Data Source - terraform-provider-coolify"
subcategory: ""
description: |-
  Lists all GitHub App integrations configured in Coolify.
---

# coolify_github_apps (Data Source)

Lists all GitHub App integrations registered with Coolify. Use this to look up a GitHub App UUID for use in `coolify_application` when deploying from a private repository.

## Example Usage

```hcl
data "coolify_github_apps" "all" {}

locals {
  my_app = one([
    for ga in data.coolify_github_apps.all.github_apps :
    ga if ga.name == "my-org-app"
  ])
}

resource "coolify_application" "private_repo" {
  type             = "private-gh-app"
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  git_repository   = "https://github.com/my-org/private-repo"
  git_branch       = "main"
  github_app_uuid  = local.my_app.uuid
  ports_exposes    = "3000"
}
```

## Schema

### Read-Only

- `github_apps` (List of Object) All GitHub App integrations.
  - `uuid` (String) GitHub App UUID.
  - `name` (String) App name.
  - `app_id` (Number) GitHub App ID.
  - `is_system_wide` (Boolean) Whether this app is available to all teams.
