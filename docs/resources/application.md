---
page_title: "coolify_application Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Deploy and manage an application on a Coolify server. Supports Docker images, Dockerfiles, public/private Git repositories, and GitHub App-connected repos.
---

# coolify_application (Resource)

Deploys and manages an application on a Coolify server. Coolify supports several application types — pick the one that matches your source:

| `type` value | Source |
|---|---|
| `dockerimage` | Pull a pre-built image from a registry |
| `dockerfile` | Build from a Dockerfile in a Git repo |
| `public` | Public Git repository (no auth) |
| `private-gh-app` | Private repo via a GitHub App integration |
| `private-deploy-key` | Private repo via an SSH deploy key |

## Example Usage

### Docker Image

```hcl
resource "coolify_application" "web" {
  type                       = "dockerimage"
  project_uuid               = coolify_project.myapp.id
  server_uuid                = var.server_uuid
  environment_name           = "production"
  name                       = "web"
  docker_registry_image_name = "nginx:latest"
  ports_exposes              = "80"
  instant_deploy             = true
}
```

### Public Git Repository

```hcl
resource "coolify_application" "api" {
  type             = "public"
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "api"
  git_repository   = "https://github.com/my-org/my-api"
  git_branch       = "main"
  build_pack       = "nixpacks"
  ports_exposes    = "3000"
  domains          = "api.example.com"
}
```

### Private GitHub App Repo

```hcl
resource "coolify_github_app" "gh" {
  name             = "my-github-app"
  api_url          = "https://api.github.com"
  html_url         = "https://github.com"
  app_id           = var.github_app_id
  installation_id  = var.github_installation_id
  client_id        = var.github_client_id
  client_secret    = var.github_client_secret
  webhook_secret   = var.github_webhook_secret
  private_key_uuid = coolify_private_key.deploy.id
}

resource "coolify_application" "private_app" {
  type             = "private-gh-app"
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "private-app"
  git_repository   = "https://github.com/my-org/private-repo"
  git_branch       = "main"
  build_pack       = "dockerfile"
  github_app_uuid  = coolify_github_app.gh.id
  ports_exposes    = "8080"
}
```

### With Environment Variables and Auto-Deploy

```hcl
resource "coolify_application" "app" {
  type                       = "dockerimage"
  project_uuid               = coolify_project.myapp.id
  server_uuid                = var.server_uuid
  environment_name           = "production"
  name                       = "my-app"
  docker_registry_image_name = "my-registry.example.com/my-app:latest"
  ports_exposes              = "3000"
  domains                    = "app.example.com"
  is_auto_deploy_enabled     = true
  is_force_https_enabled     = true
}

resource "coolify_envs_bulk" "app_vars" {
  resource_type = "application"
  resource_uuid = coolify_application.app.id
  variables = {
    DATABASE_URL = "postgresql://user:pass@localhost/db"
    REDIS_URL    = "redis://localhost:6379"
    NODE_ENV     = "production"
  }
}
```

## Schema

### Required

- `environment_name` (String) Name of the environment to deploy into (e.g. `production`, `staging`). Must match an existing environment in the project.
- `project_uuid` (String) UUID of the project this application belongs to.
- `server_uuid` (String) UUID of the server to deploy on.
- `type` (String) Application type. One of: `dockerimage`, `dockerfile`, `public`, `private-gh-app`, `private-deploy-key`.

### Optional

- `build_command` (String) Command to build the application (overrides the default for the build pack).
- `build_pack` (String) Build system to use when building from source. Common values: `nixpacks`, `dockerfile`, `static`.
- `description` (String) Optional description shown in the Coolify dashboard.
- `docker_registry_image_name` (String) Docker image name including tag, e.g. `nginx:latest` or `my-registry.example.com/app:v1.2`. Required when `type = "dockerimage"`.
- `dockerfile` (String) Inline Dockerfile content. Use this to supply the Dockerfile without a separate file in the repo.
- `domains` (String) Comma-separated list of domains to expose the application on, e.g. `app.example.com,www.example.com`.
- `git_branch` (String) Branch to deploy from (e.g. `main`, `production`).
- `git_repository` (String) Git repository URL (HTTPS or SSH).
- `github_app_uuid` (String) UUID of the `coolify_github_app` to use for private repository access.
- `install_command` (String) Command to install dependencies (e.g. `npm install`).
- `instant_deploy` (Boolean) When `true`, Coolify immediately starts a deployment after the resource is created. Defaults to `false`.
- `is_auto_deploy_enabled` (Boolean) Automatically deploy when the source branch receives a new commit.
- `is_force_https_enabled` (Boolean) Redirect all HTTP traffic to HTTPS.
- `name` (String) Display name for the application in the Coolify dashboard.
- `ports_exposes` (String) Comma-separated port numbers the container exposes, e.g. `80` or `80,443`.
- `private_key_uuid` (String) UUID of the `coolify_private_key` to use as a Git deploy key.
- `start_command` (String) Command to start the application (overrides the process file or default).

### Read-Only

- `id` (String) Coolify resource UUID. Use this to reference the application in other resources.

## Import

Import an existing application using its UUID:

```shell
terraform import coolify_application.web <application-uuid>
```
