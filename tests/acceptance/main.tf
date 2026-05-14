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
# Coolify validates cloud tokens against the provider API on creation.
# Tested manually; requires a real Hetzner/DigitalOcean token in production.
# Uncomment with a real token to test:
# resource "coolify_cloud_token" "test" {
#   name           = "acceptance-token"
#   cloud_provider = "hetzner"
#   token          = var.hetzner_api_token
# }

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

# ── 16b. Database storage ──────────────────────────────────────────────────────
# Tests POST/GET-list/DELETE /databases/{uuid}/storages

resource "coolify_database_storage" "pg_vol" {
  database_uuid = coolify_database_postgresql.pg.id
  type          = "persistent"
  mount_path    = "/var/lib/postgresql/extra"
  name          = "acceptance-pg-vol"
}
output "db_storage_id" { value = coolify_database_storage.pg_vol.id }

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

# ── 21b. Service storage ───────────────────────────────────────────────────────
# Tests POST/GET-list/DELETE /services/{uuid}/storages

variable "ghost_resource_uuid" {
  default     = ""
  description = "UUID of the ghost container within the service stack (from workflow)"
}

resource "coolify_service_storage" "svc_vol" {
  count         = var.ghost_resource_uuid != "" ? 1 : 0
  service_uuid  = coolify_service.svc.id
  type          = "persistent"
  mount_path    = "/svc-data"
  name          = "acceptance-svc-vol"
  resource_uuid = var.ghost_resource_uuid
}
output "svc_storage_id" { value = length(coolify_service_storage.svc_vol) > 0 ? coolify_service_storage.svc_vol[0].id : "" }

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

# ── 24. Bulk environment variables ────────────────────────────────────────────
# Tests PATCH /applications/{uuid}/envs/bulk, /services/{uuid}/envs/bulk,
# and /databases/{uuid}/envs/bulk

resource "coolify_envs_bulk" "app_bulk" {
  resource_type = "application"
  resource_uuid = coolify_application.app.id
  variables = {
    "BULK_VAR_1" = "value1"
    "BULK_VAR_2" = "value2"
  }
}

resource "coolify_envs_bulk" "svc_bulk" {
  resource_type = "service"
  resource_uuid = coolify_service.svc.id
  variables = {
    "SVC_BULK_1" = "svc_value1"
    "SVC_BULK_2" = "svc_value2"
  }
}

resource "coolify_envs_bulk" "db_bulk" {
  resource_type = "database"
  resource_uuid = coolify_database_postgresql.pg.id
  variables = {
    "DB_BULK_1" = "db_value1"
    "DB_BULK_2" = "db_value2"
  }
}

# ── 25. Resource lifecycle actions ────────────────────────────────────────────
# Tests GET /databases/{uuid}/start, stop, restart

resource "coolify_resource_action" "redis_start" {
  resource_type = "database"
  resource_uuid = coolify_database_redis.redis.id
  action        = "start"
}
output "redis_action_status" { value = coolify_resource_action.redis_start.status }

# ── 26. Deploy trigger ────────────────────────────────────────────────────────
# Tests GET /deploy?uuid={uuid}

resource "coolify_deploy" "app_deploy" {
  resource_uuid = coolify_application.app.id
  force         = false
}
output "deploy_uuid" { value = coolify_deploy.app_deploy.deployment_uuid }

# ── 27. System info data source ───────────────────────────────────────────────
# Tests GET /health and GET /version

data "coolify_system_info" "instance" {}
output "coolify_healthy"  { value = data.coolify_system_info.instance.healthy }
output "coolify_version"  { value = data.coolify_system_info.instance.version }

# ── 28. Server resources data source ─────────────────────────────────────────
# Tests GET /servers/{uuid}/resources

data "coolify_server_resources" "localhost_resources" {
  server_uuid = var.server_uuid
}

# ── 29. Server domains data source ───────────────────────────────────────────
# Tests GET /servers/{uuid}/domains

data "coolify_server_domains" "localhost_domains" {
  server_uuid = var.server_uuid
}

# ── 30. All resources list ────────────────────────────────────────────────────
# Tests GET /resources

data "coolify_resources" "all" {}

# ── 31. Deployments data source ───────────────────────────────────────────────
# Tests GET /deployments and GET /deployments/{uuid}
# Only query if we have a deployment UUID

output "all_resource_count" {
  value = length(data.coolify_resources.all.resources)
}

# ── 32. GitHub App ─────────────────────────────────────────────────────────────
# Tests POST/GET-list/PATCH/DELETE /github-apps
# Uses dummy credentials — Coolify stores without immediate GitHub validation

resource "coolify_github_app" "test" {
  name             = "acceptance-github-app"
  api_url          = "https://api.github.com"
  html_url         = "https://github.com"
  app_id           = 99999
  installation_id  = 99999
  client_id        = "Iv1.acceptance_test_ci"
  client_secret    = "dummy_client_secret_acceptance_test"
  webhook_secret   = "dummy_webhook_secret_test"
  private_key_uuid = coolify_private_key.test.id
}
output "github_app_id" { value = coolify_github_app.test.id }

# ── 33. PR preview management ─────────────────────────────────────────────────
# Tests DELETE /applications/{uuid}/previews/{pull_request_id}
# PR preview won't exist but the endpoint is called (returns 404 = endpoint works)

resource "coolify_application_preview" "pr1" {
  application_uuid = coolify_application.app.id
  pull_request_id  = 1
}

# ── 34. Deployments data sources ──────────────────────────────────────────────
# Tests GET /deployments and GET /deployments/applications/{uuid}

data "coolify_deployments" "all" {}

data "coolify_application_deployments" "app_deploys" {
  application_uuid = coolify_application.app.id
}

output "total_deployment_count" { value = length(data.coolify_deployments.all.deployments) }
output "app_deployment_count"   { value = length(data.coolify_application_deployments.app_deploys.deployments) }

# ── 35. List data sources ─────────────────────────────────────────────────────
# Tests GET /projects, /applications, /services, /databases, /servers,
# /security/keys, /github-apps, /cloud-tokens, /teams

data "coolify_projects" "all"     { depends_on = [coolify_project.main] }
data "coolify_applications" "all" { depends_on = [coolify_application.app] }
data "coolify_services" "all"     { depends_on = [coolify_service.svc] }
data "coolify_databases" "all"    { depends_on = [coolify_database_postgresql.pg, coolify_database_redis.redis] }
data "coolify_servers" "all"      { depends_on = [coolify_server.extra] }
data "coolify_private_keys" "all" { depends_on = [coolify_private_key.test] }
data "coolify_github_apps" "all"  { depends_on = [coolify_github_app.test] }
data "coolify_teams" "all"        {}

output "project_count"     { value = length(data.coolify_projects.all.projects) }
output "application_count" { value = length(data.coolify_applications.all.applications) }
output "service_count"     { value = length(data.coolify_services.all.services) }
output "database_count"    { value = length(data.coolify_databases.all.databases) }
output "server_count"      { value = length(data.coolify_servers.all.servers) }
output "private_key_count" { value = length(data.coolify_private_keys.all.private_keys) }
output "github_app_count"  { value = length(data.coolify_github_apps.all.github_apps) }
output "team_count"        { value = length(data.coolify_teams.all.teams) }

# ── 36. Team members data source ─────────────────────────────────────────────
# Tests GET /teams/{id}/members

data "coolify_team_members" "root" {
  team_id = data.coolify_team.current.id
}
output "team_member_count" { value = length(data.coolify_team_members.root.members) }

# ── 37. Application logs data source ─────────────────────────────────────────
# Tests GET /applications/{uuid}/logs

data "coolify_application_logs" "app_logs" {
  application_uuid = coolify_application.app.id
}

# ── 38. Server validate resource ─────────────────────────────────────────────
# Tests GET /servers/{uuid}/validate

resource "coolify_server_validate" "localhost" {
  server_uuid = var.server_uuid
}
output "server_validate_status" { value = coolify_server_validate.localhost.status }

# ── 39. API settings resource ─────────────────────────────────────────────────
# Tests GET /api/v1/enable and GET /api/v1/disable

resource "coolify_api_settings" "enabled" {
  enabled = true
}

# ── 40. Application scheduled task execution ──────────────────────────────────
# Tests DELETE /applications/{uuid}/scheduled-tasks/{task_uuid}/executions/{execution_uuid}
# Uses a placeholder UUID — a 404 on delete is accepted as success

resource "coolify_application_scheduled_task_execution" "dummy" {
  parent_uuid    = coolify_application.app.id
  task_uuid      = coolify_application_scheduled_task.task.id
  execution_uuid = "00000000-0000-0000-0000-000000000001"
}

# ── 41. Service scheduled task execution ─────────────────────────────────────
# Tests DELETE /services/{uuid}/scheduled-tasks/{task_uuid}/executions/{execution_uuid}

resource "coolify_service_scheduled_task_execution" "dummy" {
  parent_uuid    = coolify_service.svc.id
  task_uuid      = coolify_service_scheduled_task.svc_task.id
  execution_uuid = "00000000-0000-0000-0000-000000000002"
}

# ── 42. Backup execution resource ────────────────────────────────────────────
# Tests DELETE /databases/{uuid}/backups/{scheduled_backup_uuid}/executions/{execution_uuid}

resource "coolify_backup_execution" "dummy" {
  database_uuid        = coolify_database_postgresql.pg.id
  scheduled_backup_uuid = local.pg_backup_uuid
  execution_uuid       = "00000000-0000-0000-0000-000000000003"
}

# ── 43. Backup executions data source ─────────────────────────────────────────
# Tests GET /databases/{uuid}/backups/{scheduled_backup_uuid}/executions

locals {
  pg_backup_uuid = split("/", coolify_database_backup.pg_backup.id)[1]
}

data "coolify_backup_executions" "pg_execs" {
  database_uuid        = coolify_database_postgresql.pg.id
  scheduled_backup_uuid = local.pg_backup_uuid
}
output "backup_execution_count" { value = length(data.coolify_backup_executions.pg_execs.executions) }

# ── 44. Single deployment data source ────────────────────────────────────────
# Tests GET /deployments/{uuid}

data "coolify_deployment" "app_latest" {
  uuid = coolify_deploy.app_deploy.deployment_uuid
}
output "deployment_status" { value = data.coolify_deployment.app_latest.status }

# ── 45. Cloud tokens data source (always tested, may be empty) ────────────────
# Tests GET /cloud-tokens

data "coolify_cloud_tokens" "all" {}
output "cloud_token_count" { value = length(data.coolify_cloud_tokens.all.cloud_tokens) }

# ── 46. Instance settings resource + data source ──────────────────────────────
# Tests GET /settings and PATCH /settings
# Note: endpoint availability varies by Coolify version; graceful on 404

resource "coolify_instance_settings" "main" {
  is_registration_enabled = true
}

data "coolify_instance_settings" "current" {
  depends_on = [coolify_instance_settings.main]
}
output "registration_enabled" { value = data.coolify_instance_settings.current.is_registration_enabled }

# ── 47. Profile data source + resource ────────────────────────────────────────
# Tests GET /profile and PATCH /profile
# Note: endpoint availability varies by Coolify version; graceful on 404

data "coolify_profile" "me" {}
output "profile_email" { value = data.coolify_profile.me.email }

resource "coolify_profile" "me" {}

# ── 48. Current team data source ──────────────────────────────────────────────
# Tests GET /teams/current

data "coolify_current_team" "active" {}
output "current_team_name" { value = data.coolify_current_team.active.name }

# ── 51. Scheduled task data sources ──────────────────────────────────────────
# Tests GET /applications/{uuid}/scheduled-tasks/{uuid}
# Tests GET /services/{uuid}/scheduled-tasks/{uuid}

data "coolify_application_scheduled_task" "app_task" {
  parent_uuid = coolify_application.app.id
  task_uuid   = coolify_application_scheduled_task.task.id
}

data "coolify_service_scheduled_task" "svc_task" {
  parent_uuid = coolify_service.svc.id
  task_uuid   = coolify_service_scheduled_task.svc_task.id
}

# ── 52. Scheduled task executions data sources ────────────────────────────────
# Tests GET /applications/{uuid}/scheduled-tasks/{uuid}/executions
# Tests GET /services/{uuid}/scheduled-tasks/{uuid}/executions

data "coolify_application_scheduled_task_executions" "app_execs" {
  parent_uuid = coolify_application.app.id
  task_uuid   = coolify_application_scheduled_task.task.id
}
output "app_task_exec_count" { value = length(data.coolify_application_scheduled_task_executions.app_execs.executions) }

data "coolify_service_scheduled_task_executions" "svc_execs" {
  parent_uuid = coolify_service.svc.id
  task_uuid   = coolify_service_scheduled_task.svc_task.id
}
output "svc_task_exec_count" { value = length(data.coolify_service_scheduled_task_executions.svc_execs.executions) }

# ── 53. Server resource (POST /servers) ───────────────────────────────────────
# Tests POST /servers, PATCH /servers/{uuid}, DELETE /servers/{uuid}
# Uses instant_validate=false so Coolify doesn't SSH — just registers the record

resource "coolify_server" "extra" {
  name             = "acceptance-extra-server"
  ip               = "127.0.0.2"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.test.id
  instant_validate = false
}
output "extra_server_uuid" { value = coolify_server.extra.id }

# ── 54. Cloud token validate resource ─────────────────────────────────────────
# Tests POST /cloud-tokens/{uuid}/validate
# Only runs if COOLIFY_TEST_CLOUD_TOKEN_UUID is set (populated by Hetzner test step)

variable "cloud_token_uuid" {
  default     = ""
  description = "UUID of an existing cloud token to validate (optional)"
}

resource "coolify_cloud_token_validate" "hetzner" {
  count           = var.cloud_token_uuid != "" ? 1 : 0
  cloud_token_uuid = var.cloud_token_uuid
}
output "cloud_token_validate_status" {
  value = length(coolify_cloud_token_validate.hetzner) > 0 ? coolify_cloud_token_validate.hetzner[0].status : "skipped"
}

# ── 55. Application scheduled tasks list ──────────────────────────────────────
# Tests GET /applications/{uuid}/scheduled-tasks

data "coolify_application_scheduled_tasks" "app_tasks" {
  application_uuid = coolify_application.app.id
  depends_on       = [coolify_application_scheduled_task.task]
}
output "app_task_count" { value = length(data.coolify_application_scheduled_tasks.app_tasks.tasks) }

# ── 56. Service scheduled tasks list ─────────────────────────────────────────
# Tests GET /services/{uuid}/scheduled-tasks

data "coolify_service_scheduled_tasks" "svc_tasks" {
  service_uuid = coolify_service.svc.id
  depends_on   = [coolify_service_scheduled_task.svc_task]
}
output "svc_task_count" { value = length(data.coolify_service_scheduled_tasks.svc_tasks.tasks) }

# ── 57. Database backups list ─────────────────────────────────────────────────
# Tests GET /databases/{uuid}/backups

data "coolify_database_backups" "pg_backups" {
  database_uuid = coolify_database_postgresql.pg.id
  depends_on    = [coolify_database_backup.pg_backup]
}
output "pg_backup_uuid" { value = local.pg_backup_uuid }
output "backup_schedule_count" { value = length(data.coolify_database_backups.pg_backups.backups) }

# ── 58. Project environments list ─────────────────────────────────────────────
# Tests GET /projects/{uuid}/environments

data "coolify_project_environments" "main_envs" {
  project_uuid = coolify_project.main.id
  depends_on   = [coolify_environment.staging]
}
output "env_count" { value = length(data.coolify_project_environments.main_envs.environments) }

# ── 59. Current team members data source ─────────────────────────────────────
# Tests GET /teams/current/members

data "coolify_current_team_members" "active" {}
output "current_team_member_count" { value = length(data.coolify_current_team_members.active.members) }

# ── 60. GitHub App repositories data source ───────────────────────────────────
# Tests GET /github-apps/{uuid}/repositories
# Returns empty list for dummy/test apps — endpoint availability is what matters

data "coolify_github_app_repositories" "test" {
  github_app_uuid = coolify_github_app.test.id
}
output "github_repo_count" { value = length(data.coolify_github_app_repositories.test.repositories) }
