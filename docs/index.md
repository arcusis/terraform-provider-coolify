---
page_title: "coolify Provider"
description: |-
  Manage Coolify projects, servers, databases, applications, and services using Terraform.
---

# Coolify Provider

The **Coolify** provider lets you manage every resource on a [Coolify](https://coolify.io) instance through Terraform. Use it to provision projects, environments, applications, databases, services, servers, and more — all from Infrastructure as Code.

## Authentication

The provider requires a Coolify API bearer token. You can generate one from **Settings → API Tokens** in the Coolify dashboard.

```hcl
provider "coolify" {
  endpoint = "https://coolify.example.com"
  token    = var.coolify_token
}
```

Both values can also be supplied through environment variables:

```shell
export COOLIFY_ENDPOINT=https://coolify.example.com
export COOLIFY_TOKEN=your-token-here
```

## Example Usage

```hcl
terraform {
  required_providers {
    coolify = {
      source  = "arcusis/coolify"
      version = "~> 1.0"
    }
  }
}

provider "coolify" {
  endpoint = var.coolify_endpoint
  token    = var.coolify_token
}

# Create a project
resource "coolify_project" "myapp" {
  name        = "my-application"
  description = "Production workloads"
}

# Create an environment inside the project
resource "coolify_environment" "production" {
  project_uuid = coolify_project.myapp.id
  name         = "production"
}

# Deploy a Docker image application
resource "coolify_application" "web" {
  type                       = "dockerimage"
  project_uuid               = coolify_project.myapp.id
  server_uuid                = var.server_uuid
  environment_name           = coolify_environment.production.name
  name                       = "web"
  docker_registry_image_name = "nginx:latest"
  ports_exposes              = "80"
}

# Attach a PostgreSQL database
resource "coolify_database_postgresql" "db" {
  project_uuid     = coolify_project.myapp.id
  server_uuid      = var.server_uuid
  environment_name = coolify_environment.production.name
  name             = "app-db"
  postgres_user    = "appuser"
  postgres_db      = "app"
}
```

## Resource Overview

### Projects & Environments

| Resource | Description |
|---|---|
| `coolify_project` | Top-level grouping for all workloads |
| `coolify_environment` | Named environment inside a project (e.g. production, staging) |

### Servers & Keys

| Resource | Description |
|---|---|
| `coolify_server` | Register a server with Coolify over SSH |
| `coolify_server_validate` | Verify Coolify can connect to a server |
| `coolify_private_key` | Store an SSH private key for server authentication |
| `coolify_cloud_token` | Store a cloud provider API token (Hetzner, DigitalOcean) |

### Applications

| Resource | Description |
|---|---|
| `coolify_application` | Any application type: docker image, Dockerfile, Git repo |
| `coolify_application_storage` | Persistent volume mount for an application |
| `coolify_application_scheduled_task` | Cron-style scheduled command inside an application |

### Services (One-Click Stacks)

| Resource | Description |
|---|---|
| `coolify_service` | Deploy a pre-built service stack (Ghost, WordPress, Plausible, etc.) |
| `coolify_service_scheduled_task` | Scheduled task inside a service |

### Databases

| Resource | Description |
|---|---|
| `coolify_database_postgresql` | PostgreSQL database |
| `coolify_database_mysql` | MySQL database |
| `coolify_database_mariadb` | MariaDB database |
| `coolify_database_redis` | Redis |
| `coolify_database_mongodb` | MongoDB |
| `coolify_database_keydb` | KeyDB |
| `coolify_database_dragonfly` | DragonflyDB |
| `coolify_database_clickhouse` | Clickhouse |
| `coolify_database_backup` | Scheduled backup for a database |

### Environment Variables

| Resource | Description |
|---|---|
| `coolify_environment_variable` | Single env var for an application |
| `coolify_service_environment_variable` | Single env var for a service |
| `coolify_database_environment_variable` | Single env var for a database |
| `coolify_envs_bulk` | Set multiple env vars in one operation |

### Operations & Lifecycle

| Resource | Description |
|---|---|
| `coolify_deploy` | Trigger a deployment |
| `coolify_resource_action` | Start, stop, or restart a resource |
| `coolify_application_preview` | Manage PR preview environments |
| `coolify_api_settings` | Enable or disable the Coolify API |

### GitHub Integration

| Resource | Description |
|---|---|
| `coolify_github_app` | Register a GitHub App for source control integration |

## Schema

### Optional

- `endpoint` (String) Base URL of your Coolify instance, e.g. `https://coolify.example.com`. Can also be set with the `COOLIFY_ENDPOINT` environment variable.
- `token` (String, Sensitive) Coolify API bearer token. Generate one from **Settings → API Tokens**. Can also be set with `COOLIFY_TOKEN`.
