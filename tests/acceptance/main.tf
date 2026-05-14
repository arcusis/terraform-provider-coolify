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

# ── Input variables ────────────────────────────────────────────────────────────

variable "coolify_endpoint" {}
variable "coolify_token" { sensitive = true }
variable "server_uuid" { description = "UUID of the Coolify localhost server" }
variable "test_private_key" {
  sensitive   = true
  description = "RSA private key in PEM format for testing"
}

# ── 1. Project + data source ───────────────────────────────────────────────────

resource "coolify_project" "main" {
  name        = "acceptance-project"
  description = "Full acceptance test project"
}
output "project_uuid" { value = coolify_project.main.id }

data "coolify_project" "readback" { uuid = coolify_project.main.id }
output "project_name_readback" { value = data.coolify_project.readback.name }

# ── 2. Project environment ─────────────────────────────────────────────────────

resource "coolify_environment" "staging" {
  project_uuid = coolify_project.main.id
  name         = "staging"
}
output "environment_id" { value = coolify_environment.staging.id }

# ── 3. Private key + data source ───────────────────────────────────────────────

resource "coolify_private_key" "test" {
  name        = "acceptance-key"
  description = "Test SSH key for acceptance suite"
  private_key = var.test_private_key
}
output "private_key_uuid" { value = coolify_private_key.test.id }

data "coolify_private_key" "readback" { uuid = coolify_private_key.test.id }
output "private_key_name_readback" { value = data.coolify_private_key.readback.name }

# ── 4. Cloud token ─────────────────────────────────────────────────────────────

resource "coolify_cloud_token" "test" {
  name           = "acceptance-token"
  cloud_provider = "hetzner"
  token          = "dummy-ci-token-not-for-real-use"
}
output "cloud_token_uuid" { value = coolify_cloud_token.test.id }

# ── 5. Server data source ──────────────────────────────────────────────────────

data "coolify_server" "localhost" { uuid = var.server_uuid }
output "server_name" { value = data.coolify_server.localhost.name }

# ── 6. Team data source ────────────────────────────────────────────────────────

data "coolify_team" "current" { name = "Root Team" }
output "team_id" { value = data.coolify_team.current.id }

# ── 7. PostgreSQL database ─────────────────────────────────────────────────────

resource "coolify_database_postgresql" "pg" {
  server_uuid      = var.server_uuid
  project_uuid     = coolify_project.main.id
  environment_name = "production"
  name             = "acceptance-pg"
  postgres_user    = "pguser"
  postgres_db      = "acceptance"
  instant_deploy   = false
}
output "pg_uuid" { value = coolify_database_postgresql.pg.id }

data "coolify_database" "pg_readback" { uuid = coolify_database_postgresql.pg.id }

# ── 8. Redis ───────────────────────────────────────────────────────────────────

resource "coolify_database_redis" "redis" {
  server_uuid      = var.server_uuid
  project_uuid     = coolify_project.main.id
  environment_name = "production"
  name             = "acceptance-redis"
  instant_deploy   = false
}
output "redis_uuid" { value = coolify_database_redis.redis.id }

# ── 9. MySQL ───────────────────────────────────────────────────────────────────

resource "coolify_database_mysql" "mysql" {
  server_uuid      = var.server_uuid
  project_uuid     = coolify_project.main.id
  environment_name = "production"
  name             = "acceptance-mysql"
  mysql_database   = "acceptance"
  instant_deploy   = false
}
output "mysql_uuid" { value = coolify_database_mysql.mysql.id }

# ── 10. MariaDB ────────────────────────────────────────────────────────────────

resource "coolify_database_mariadb" "mariadb" {
  server_uuid      = var.server_uuid
  project_uuid     = coolify_project.main.id
  environment_name = "production"
  name             = "acceptance-mariadb"
  instant_deploy   = false
}
output "mariadb_uuid" { value = coolify_database_mariadb.mariadb.id }

# ── 11. MongoDB ────────────────────────────────────────────────────────────────

resource "coolify_database_mongodb" "mongo" {
  server_uuid      = var.server_uuid
  project_uuid     = coolify_project.main.id
  environment_name = "production"
  name             = "acceptance-mongo"
  instant_deploy   = false
}
output "mongo_uuid" { value = coolify_database_mongodb.mongo.id }

# ── 12. KeyDB ──────────────────────────────────────────────────────────────────

resource "coolify_database_keydb" "keydb" {
  server_uuid      = var.server_uuid
  project_uuid     = coolify_project.main.id
  environment_name = "production"
  name             = "acceptance-keydb"
  instant_deploy   = false
}
output "keydb_uuid" { value = coolify_database_keydb.keydb.id }

# ── 13. Dragonfly ──────────────────────────────────────────────────────────────

resource "coolify_database_dragonfly" "dragonfly" {
  server_uuid      = var.server_uuid
  project_uuid     = coolify_project.main.id
  environment_name = "production"
  name             = "acceptance-dragonfly"
  instant_deploy   = false
}
output "dragonfly_uuid" { value = coolify_database_dragonfly.dragonfly.id }

# ── 14. Clickhouse ─────────────────────────────────────────────────────────────

resource "coolify_database_clickhouse" "clickhouse" {
  server_uuid                = var.server_uuid
  project_uuid               = coolify_project.main.id
  environment_name           = "production"
  name                       = "acceptance-clickhouse"
  clickhouse_admin_user      = "admin"
  instant_deploy             = false
}
output "clickhouse_uuid" { value = coolify_database_clickhouse.clickhouse.id }

# ── 15. Database backup ────────────────────────────────────────────────────────

resource "coolify_database_backup" "pg_backup" {
  database_uuid = coolify_database_postgresql.pg.id
  frequency     = "daily"
  enabled       = true
}
output "backup_id" { value = coolify_database_backup.pg_backup.id }

# ── 16. Database environment variable ─────────────────────────────────────────

resource "coolify_database_environment_variable" "pg_var" {
  database_uuid = coolify_database_postgresql.pg.id
  key           = "DB_ACCEPTANCE_TEST"
  value         = "true"
}

# ── 17. Application (Docker image) ────────────────────────────────────────────

resource "coolify_application" "app" {
  type                       = "dockerimage"
  project_uuid               = coolify_project.main.id
  server_uuid                = var.server_uuid
  environment_name           = "production"
  docker_registry_image_name = "nginx"
  name                       = "acceptance-app"
  ports_exposes              = "80"
  instant_deploy             = false
}
output "app_uuid" { value = coolify_application.app.id }

data "coolify_application" "app_readback" { uuid = coolify_application.app.id }

# ── 18. Application storage ────────────────────────────────────────────────────

resource "coolify_application_storage" "vol" {
  application_uuid = coolify_application.app.id
  type             = "persistent"
  mount_path       = "/data"
  name             = "acceptance-vol"
}
output "app_storage_id" { value = coolify_application_storage.vol.id }

# ── 19. Application scheduled task ────────────────────────────────────────────

resource "coolify_application_scheduled_task" "task" {
  application_uuid = coolify_application.app.id
  name             = "acceptance-task"
  command          = "echo hello"
  frequency        = "0 * * * *"
  enabled          = true
}
output "app_task_id" { value = coolify_application_scheduled_task.task.id }

# ── 20. Application environment variable ──────────────────────────────────────

resource "coolify_environment_variable" "app_var" {
  application_uuid = coolify_application.app.id
  key              = "ACCEPTANCE_TEST"
  value            = "true"
  is_literal       = false
}

# ── 21. Service ────────────────────────────────────────────────────────────────

resource "coolify_service" "svc" {
  type             = "ghost"
  project_uuid     = coolify_project.main.id
  server_uuid      = var.server_uuid
  environment_name = "production"
  name             = "acceptance-service"
  instant_deploy   = false
}
output "service_uuid" { value = coolify_service.svc.id }

data "coolify_service" "svc_readback" { uuid = coolify_service.svc.id }

# ── 22. Service scheduled task ────────────────────────────────────────────────

resource "coolify_service_scheduled_task" "svc_task" {
  service_uuid = coolify_service.svc.id
  name         = "acceptance-svc-task"
  command      = "echo svc"
  frequency    = "0 * * * *"
  enabled      = true
}
output "svc_task_id" { value = coolify_service_scheduled_task.svc_task.id }

# ── 23. Service environment variable ──────────────────────────────────────────

resource "coolify_service_environment_variable" "svc_var" {
  service_uuid = coolify_service.svc.id
  key          = "SVC_ACCEPTANCE_TEST"
  value        = "true"
}
