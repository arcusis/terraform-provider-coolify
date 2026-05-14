# Terraform Provider for Coolify

A Terraform provider for managing resources on a [Coolify](https://coolify.io) self-hosted instance. Covers the full Coolify API surface.

[![Acceptance Tests](https://github.com/arcusis/terraform-provider-coolify/actions/workflows/acceptance.yml/badge.svg)](https://github.com/arcusis/terraform-provider-coolify/actions/workflows/acceptance.yml)

## Requirements

- Terraform >= 1.10
- Go >= 1.21 (to build from source)
- A running Coolify instance with an API token

## Installation

```hcl
terraform {
  required_providers {
    coolify = {
      source  = "arcusis/coolify"
      version = "~> 1.0"
    }
  }
}
```

## Authentication

```hcl
provider "coolify" {
  endpoint = "https://coolify.example.com"
  token    = var.coolify_token
}
```

Environment variables are also supported: `COOLIFY_ENDPOINT` and `COOLIFY_TOKEN`.

## Resources

| Resource | Description |
|---|---|
| `coolify_project` | Project |
| `coolify_environment` | Project environment (staging, production, etc.) |
| `coolify_server` | Server |
| `coolify_private_key` | SSH private key |
| `coolify_cloud_token` | Cloud provider API token (Hetzner, DigitalOcean) |
| `coolify_github_app` | GitHub App integration |
| `coolify_application` | Application (public git, private git, Dockerfile, Docker image, Docker Compose) |
| `coolify_application_storage` | Persistent volume or file mount for an application |
| `coolify_application_scheduled_task` | Cron task for an application |
| `coolify_environment_variable` | Environment variable for an application |
| `coolify_service` | One-click service (Ghost, WordPress, Plausible, etc.) |
| `coolify_service_storage` | Persistent volume or file mount for a service |
| `coolify_service_scheduled_task` | Cron task for a service |
| `coolify_service_environment_variable` | Environment variable for a service |
| `coolify_database_postgresql` | PostgreSQL database |
| `coolify_database_mysql` | MySQL database |
| `coolify_database_mariadb` | MariaDB database |
| `coolify_database_redis` | Redis database |
| `coolify_database_mongodb` | MongoDB database |
| `coolify_database_keydb` | KeyDB database |
| `coolify_database_dragonfly` | Dragonfly database |
| `coolify_database_clickhouse` | Clickhouse database |
| `coolify_database_backup` | Scheduled backup configuration for a database |
| `coolify_database_storage` | Persistent volume or file mount for a database |
| `coolify_database_environment_variable` | Environment variable for a database |

## Data Sources

| Data Source | Description |
|---|---|
| `coolify_project` | Look up a project by UUID |
| `coolify_server` | Look up a server by UUID |
| `coolify_private_key` | Look up a private key by UUID |
| `coolify_application` | Look up an application by UUID |
| `coolify_service` | Look up a service by UUID |
| `coolify_database` | Look up a database by UUID |
| `coolify_team` | Look up a team by ID or name |
| `coolify_deployment` | Look up a deployment by UUID |
| `coolify_hetzner_locations` | List Hetzner locations for a cloud token |
| `coolify_hetzner_server_types` | List Hetzner server types for a cloud token |

## Example

```hcl
resource "coolify_project" "my_app" {
  name        = "my-app"
  description = "Production application"
}

resource "coolify_private_key" "deploy_key" {
  name        = "deploy-key"
  private_key = var.ssh_private_key
}

resource "coolify_server" "web" {
  name             = "web-01"
  ip               = "203.0.113.10"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.deploy_key.id
}

resource "coolify_application" "api" {
  type             = "public"
  project_uuid     = coolify_project.my_app.id
  server_uuid      = coolify_server.web.id
  environment_name = "production"
  git_repository   = "https://github.com/my-org/my-api"
  git_branch       = "main"
  build_pack       = "nixpacks"
  name             = "api"
  ports_exposes    = "3000"
  instant_deploy   = true
}

resource "coolify_environment_variable" "db_url" {
  application_uuid = coolify_application.api.id
  key              = "DATABASE_URL"
  value            = "postgresql://..."
}
```

## Development

```bash
# Build
go build -o terraform-provider-coolify .

# Configure dev override
cat > ~/.terraformrc <<EOF
provider_installation {
  dev_overrides { "arcusis/coolify" = "$PWD" }
  direct {}
}
EOF

# Run acceptance tests (requires a running Coolify instance)
cd tests/acceptance
terraform plan
terraform apply -auto-approve
terraform destroy -auto-approve
```

## License

MIT
