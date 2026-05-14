---
page_title: "coolify_github_app Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Registers a GitHub App integration with Coolify, enabling deployment of private GitHub repositories without per-repo deploy keys.
---

# coolify_github_app (Resource)

Registers a GitHub App integration with Coolify. Once registered, the GitHub App can be used to deploy applications from private repositories. This is the recommended approach for organizations that manage many private repos, as it grants access at the installation level rather than per-repo.

## Example Usage

```hcl
resource "coolify_private_key" "gh_app_key" {
  name        = "github-app-private-key"
  private_key = var.github_app_private_key_pem
}

resource "coolify_github_app" "main" {
  name            = "my-github-app"
  api_url         = "https://api.github.com"
  html_url        = "https://github.com"
  app_id          = var.github_app_id
  installation_id = var.github_installation_id
  client_id       = var.github_client_id
  client_secret   = var.github_client_secret
  webhook_secret  = var.github_webhook_secret
  private_key_uuid = coolify_private_key.gh_app_key.id
}

resource "coolify_application" "api" {
  type             = "private-gh-app"
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "api"
  git_repository   = "https://github.com/my-org/private-api"
  git_branch       = "main"
  build_pack       = "nixpacks"
  github_app_uuid  = coolify_github_app.main.id
  ports_exposes    = "3000"
}
```

### Organization-scoped GitHub App

```hcl
resource "coolify_github_app" "org_app" {
  name             = "my-org-github-app"
  api_url          = "https://api.github.com"
  html_url         = "https://github.com"
  app_id           = var.github_app_id
  installation_id  = var.github_installation_id
  client_id        = var.github_client_id
  client_secret    = var.github_client_secret
  webhook_secret   = var.github_webhook_secret
  private_key_uuid = coolify_private_key.gh_app_key.id
  organization     = "my-org"
  is_system_wide   = true
}
```

## Schema

### Required

- `api_url` (String) GitHub API base URL. Use `https://api.github.com` for GitHub.com, or your GitHub Enterprise Server URL.
- `app_id` (Number) GitHub App ID found on the App settings page.
- `client_id` (String) OAuth client ID of the GitHub App.
- `client_secret` (String, Sensitive) OAuth client secret of the GitHub App.
- `html_url` (String) GitHub web base URL. Use `https://github.com` for GitHub.com, or your GitHub Enterprise Server URL.
- `installation_id` (Number) Installation ID of the GitHub App on the target account or organization.
- `name` (String) Display name for this GitHub App integration in the Coolify dashboard.
- `private_key_uuid` (String) UUID of the `coolify_private_key` holding the GitHub App's RSA private key (downloaded from the App settings page).
- `webhook_secret` (String, Sensitive) Webhook secret configured on the GitHub App for verifying incoming webhook payloads.

### Optional

- `is_system_wide` (Boolean) When `true`, this GitHub App is available to all users on the Coolify instance. Defaults to `false`.
- `organization` (String) GitHub organization name if the App is installed on an organization rather than a personal account.

### Read-Only

- `id` (String) Coolify resource UUID. Use this as `github_app_uuid` in `coolify_application` resources.

## Import

Import an existing GitHub App integration using its UUID:

```shell
terraform import coolify_github_app.main <github-app-uuid>
```
