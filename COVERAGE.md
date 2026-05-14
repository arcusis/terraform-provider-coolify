# Coolify API Coverage Audit

Source: `https://raw.githubusercontent.com/coollabsio/coolify/v4.x/openapi.json`
Total API endpoints: **135**
Provider coverage: **135/135 (100%)**

## Coverage Table

| # | Endpoint | Terraform Resource/Data Source | Acceptance Test |
|---|----------|-------------------------------|-----------------|
| 1 | GET /applications | `coolify_applications` (DS) | ✓ |
| 2 | POST /applications/dockercompose | `coolify_application` (R, type=dockercompose) | via type param |
| 3 | POST /applications/dockerfile | `coolify_application` (R, type=dockerfile) | via type param |
| 4 | POST /applications/dockerimage | `coolify_application` (R, type=dockerimage) | ✓ explicit |
| 5 | POST /applications/private-deploy-key | `coolify_application` (R, type=private-deploy-key) | via type param |
| 6 | POST /applications/private-github-app | `coolify_application` (R, type=private-gh-app) | via type param |
| 7 | POST /applications/public | `coolify_application` (R, type=public) | via type param |
| 8 | GET /applications/{uuid} | `coolify_application` (DS) | ✓ |
| 9 | DELETE /applications/{uuid} | `coolify_application` (R destroy) | ✓ |
| 10 | PATCH /applications/{uuid} | `coolify_application` (R update) | ✓ |
| 11 | GET /applications/{uuid}/envs | `coolify_environment_variable` (R) | ✓ |
| 12 | POST /applications/{uuid}/envs | `coolify_environment_variable` (R create) | ✓ |
| 13 | PATCH /applications/{uuid}/envs | `coolify_environment_variable` (R update) | ✓ |
| 14 | PATCH /applications/{uuid}/envs/bulk | `coolify_envs_bulk` (R) | ✓ |
| 15 | DELETE /applications/{uuid}/envs/{env_uuid} | `coolify_environment_variable` (R destroy) | ✓ |
| 16 | GET /applications/{uuid}/logs | `coolify_application_logs` (DS) | ✓ |
| 17 | DELETE /applications/{uuid}/previews/{pr_id} | `coolify_application_preview` (R destroy) | ✓ |
| 18 | GET /applications/{uuid}/restart | `coolify_resource_action` (R, action=restart) | ✓ |
| 19 | GET /applications/{uuid}/scheduled-tasks | `coolify_application_scheduled_tasks` (DS) | ✓ |
| 20 | POST /applications/{uuid}/scheduled-tasks | `coolify_application_scheduled_task` (R create) | ✓ |
| 21 | DELETE /applications/{uuid}/scheduled-tasks/{uuid} | `coolify_application_scheduled_task` (R destroy) | ✓ |
| 22 | PATCH /applications/{uuid}/scheduled-tasks/{uuid} | `coolify_application_scheduled_task` (R update) | ✓ |
| 23 | GET /applications/{uuid}/scheduled-tasks/{uuid}/executions | `coolify_application_scheduled_task_executions` (DS) | ✓ |
| 24 | GET /applications/{uuid}/start | `coolify_resource_action` (R, action=start) | ✓ |
| 25 | GET /applications/{uuid}/stop | `coolify_resource_action` (R, action=stop) | ✓ |
| 26 | GET /applications/{uuid}/storages | `coolify_application_storage` (R, list in read) | ✓ |
| 27 | POST /applications/{uuid}/storages | `coolify_application_storage` (R create) | ✓ |
| 28 | PATCH /applications/{uuid}/storages | `coolify_application_storage` (R update) | ✓ |
| 29 | DELETE /applications/{uuid}/storages/{uuid} | `coolify_application_storage` (R destroy) | ✓ |
| 30 | GET /cloud-tokens | `coolify_cloud_tokens` (DS) | ✓ |
| 31 | POST /cloud-tokens | `coolify_cloud_token` (R create) | ✓ cond. |
| 32 | GET /cloud-tokens/{uuid} | `coolify_cloud_token` (R read) | ✓ cond. |
| 33 | DELETE /cloud-tokens/{uuid} | `coolify_cloud_token` (R destroy) | ✓ cond. |
| 34 | PATCH /cloud-tokens/{uuid} | `coolify_cloud_token` (R update) | ✓ cond. |
| 35 | POST /cloud-tokens/{uuid}/validate | `coolify_cloud_token_validate` (R) | ✓ cond. |
| 36 | GET /databases | `coolify_databases` (DS) | ✓ |
| 37 | POST /databases/clickhouse | `coolify_database_clickhouse` (R) | ✓ |
| 38 | POST /databases/dragonfly | `coolify_database_dragonfly` (R) | ✓ |
| 39 | POST /databases/keydb | `coolify_database_keydb` (R) | ✓ |
| 40 | POST /databases/mariadb | `coolify_database_mariadb` (R) | ✓ |
| 41 | POST /databases/mongodb | `coolify_database_mongodb` (R) | ✓ |
| 42 | POST /databases/mysql | `coolify_database_mysql` (R) | ✓ |
| 43 | POST /databases/postgresql | `coolify_database_postgresql` (R) | ✓ |
| 44 | POST /databases/redis | `coolify_database_redis` (R) | ✓ |
| 45 | GET /databases/{uuid} | `coolify_database` (DS) | ✓ |
| 46 | DELETE /databases/{uuid} | `coolify_database_*` (R destroy) | ✓ |
| 47 | PATCH /databases/{uuid} | `coolify_database_*` (R update) | ✓ |
| 48 | GET /databases/{uuid}/backups | `coolify_database_backups` (DS) | ✓ |
| 49 | POST /databases/{uuid}/backups | `coolify_database_backup` (R create) | ✓ |
| 50 | DELETE /databases/{uuid}/backups/{uuid} | `coolify_database_backup` (R destroy) | ✓ |
| 51 | PATCH /databases/{uuid}/backups/{uuid} | `coolify_database_backup` (R update) | ✓ |
| 52 | GET /databases/{uuid}/backups/{uuid}/executions | `coolify_backup_executions` (DS) | ✓ |
| 53 | DELETE /databases/{uuid}/backups/{uuid}/executions/{uuid} | `coolify_backup_execution` (R destroy) | ✓ |
| 54 | GET /databases/{uuid}/envs | `coolify_database_environment_variable` (R) | ✓ |
| 55 | POST /databases/{uuid}/envs | `coolify_database_environment_variable` (R create) | ✓ |
| 56 | PATCH /databases/{uuid}/envs | `coolify_database_environment_variable` (R update) | ✓ |
| 57 | PATCH /databases/{uuid}/envs/bulk | `coolify_envs_bulk` (R, type=database) | ✓ |
| 58 | DELETE /databases/{uuid}/envs/{uuid} | `coolify_database_environment_variable` (R destroy) | ✓ |
| 59 | GET /databases/{uuid}/restart | `coolify_resource_action` (R, action=restart) | ✓ |
| 60 | GET /databases/{uuid}/start | `coolify_resource_action` (R, action=start) | ✓ |
| 61 | GET /databases/{uuid}/stop | `coolify_resource_action` (R, action=stop) | ✓ |
| 62 | GET /databases/{uuid}/storages | `coolify_database_storage` (R, list in read) | ✓ |
| 63 | POST /databases/{uuid}/storages | `coolify_database_storage` (R create) | ✓ |
| 64 | PATCH /databases/{uuid}/storages | `coolify_database_storage` (R update) | ✓ |
| 65 | DELETE /databases/{uuid}/storages/{uuid} | `coolify_database_storage` (R destroy) | ✓ |
| 66 | GET /deploy | `coolify_deploy` (R create) | ✓ |
| 67 | GET /deployments | `coolify_deployments` (DS) | ✓ |
| 68 | GET /deployments/applications/{uuid} | `coolify_application_deployments` (DS) | ✓ |
| 69 | GET /deployments/{uuid} | `coolify_deployment` (DS) | ✓ |
| 70 | POST /deployments/{uuid}/cancel | `coolify_deploy` (R destroy) | ✓ |
| 71 | GET /disable | `coolify_api_settings` (R, enabled=false) | ✓ |
| 72 | GET /enable | `coolify_api_settings` (R, enabled=true) | ✓ |
| 73 | GET /github-apps | `coolify_github_apps` (DS) | ✓ |
| 74 | POST /github-apps | `coolify_github_app` (R create) | ✓ |
| 75 | DELETE /github-apps/{uuid} | `coolify_github_app` (R destroy) | ✓ |
| 76 | PATCH /github-apps/{uuid} | `coolify_github_app` (R update) | ✓ |
| 77 | GET /github-apps/{uuid}/repositories | `coolify_github_app_repositories` (DS) | ✓ |
| 78 | GET /github-apps/{uuid}/repositories/{owner}/{repo}/branches | `coolify_github_app_branches` (DS) | via workflow |
| 79 | GET /health | `coolify_system_info` (DS) | ✓ |
| 80 | GET /hetzner/images | `coolify_hetzner_images` (DS) | ✓ cond. |
| 81 | GET /hetzner/locations | `coolify_hetzner_locations` (DS) | ✓ cond. |
| 82 | GET /hetzner/server-types | `coolify_hetzner_server_types` (DS) | ✓ cond. |
| 83 | GET /hetzner/ssh-keys | `coolify_hetzner_ssh_keys` (DS) | ✓ cond. |
| 84 | GET /projects | `coolify_projects` (DS) | ✓ |
| 85 | POST /projects | `coolify_project` (R create) | ✓ |
| 86 | GET /projects/{uuid} | `coolify_project` (DS) | ✓ |
| 87 | DELETE /projects/{uuid} | `coolify_project` (R destroy) | ✓ |
| 88 | PATCH /projects/{uuid} | `coolify_project` (R update) | ✓ |
| 89 | GET /projects/{uuid}/environments | `coolify_project_environments` (DS) | ✓ |
| 90 | POST /projects/{uuid}/environments | `coolify_environment` (R create) | ✓ |
| 91 | DELETE /projects/{uuid}/environments/{name} | `coolify_environment` (R destroy) | ✓ |
| 92 | GET /projects/{uuid}/{env_name} | `coolify_environment` (R read) | ✓ |
| 93 | GET /resources | `coolify_resources` (DS) | ✓ |
| 94 | GET /security/keys | `coolify_private_keys` (DS) | ✓ |
| 95 | POST /security/keys | `coolify_private_key` (R create) | ✓ |
| 96 | PATCH /security/keys | `coolify_private_key` (R update, uuid in body) | ✓ |
| 97 | GET /security/keys/{uuid} | `coolify_private_key` (DS) | ✓ |
| 98 | DELETE /security/keys/{uuid} | `coolify_private_key` (R destroy) | ✓ |
| 99 | GET /servers | `coolify_servers` (DS) | ✓ |
| 100 | POST /servers | `coolify_server` (R create) | ✓ |
| 101 | POST /servers/hetzner | `coolify_server_hetzner` (R create) | ✓ cond. |
| 102 | GET /servers/{uuid} | `coolify_server` (DS) | ✓ |
| 103 | DELETE /servers/{uuid} | `coolify_server` (R destroy) | ✓ |
| 104 | PATCH /servers/{uuid} | `coolify_server` (R update) | ✓ |
| 105 | GET /servers/{uuid}/domains | `coolify_server_domains` (DS) | ✓ |
| 106 | GET /servers/{uuid}/resources | `coolify_server_resources` (DS) | ✓ |
| 107 | GET /servers/{uuid}/validate | `coolify_server_validate` (R) | ✓ |
| 108 | GET /services | `coolify_services` (DS) | ✓ |
| 109 | POST /services | `coolify_service` (R create) | ✓ |
| 110 | GET /services/{uuid} | `coolify_service` (DS) | ✓ |
| 111 | DELETE /services/{uuid} | `coolify_service` (R destroy) | ✓ |
| 112 | PATCH /services/{uuid} | `coolify_service` (R update) | ✓ |
| 113 | GET /services/{uuid}/envs | `coolify_service_environment_variable` (R) | ✓ |
| 114 | POST /services/{uuid}/envs | `coolify_service_environment_variable` (R create) | ✓ |
| 115 | PATCH /services/{uuid}/envs | `coolify_service_environment_variable` (R update) | ✓ |
| 116 | PATCH /services/{uuid}/envs/bulk | `coolify_envs_bulk` (R, type=service) | ✓ |
| 117 | DELETE /services/{uuid}/envs/{uuid} | `coolify_service_environment_variable` (R destroy) | ✓ |
| 118 | GET /services/{uuid}/restart | `coolify_resource_action` (R, action=restart) | ✓ |
| 119 | GET /services/{uuid}/scheduled-tasks | `coolify_service_scheduled_tasks` (DS) | ✓ |
| 120 | POST /services/{uuid}/scheduled-tasks | `coolify_service_scheduled_task` (R create) | ✓ |
| 121 | DELETE /services/{uuid}/scheduled-tasks/{uuid} | `coolify_service_scheduled_task` (R destroy) | ✓ |
| 122 | PATCH /services/{uuid}/scheduled-tasks/{uuid} | `coolify_service_scheduled_task` (R update) | ✓ |
| 123 | GET /services/{uuid}/scheduled-tasks/{uuid}/executions | `coolify_service_scheduled_task_executions` (DS) | ✓ |
| 124 | GET /services/{uuid}/start | `coolify_resource_action` (R, action=start) | ✓ |
| 125 | GET /services/{uuid}/stop | `coolify_resource_action` (R, action=stop) | ✓ |
| 126 | GET /services/{uuid}/storages | `coolify_service_storage` (R, list in read) | ✓ |
| 127 | POST /services/{uuid}/storages | `coolify_service_storage` (R create) | ✓ |
| 128 | PATCH /services/{uuid}/storages | `coolify_service_storage` (R update) | ✓ |
| 129 | DELETE /services/{uuid}/storages/{uuid} | `coolify_service_storage` (R destroy) | ✓ |
| 130 | GET /teams | `coolify_teams` (DS) | ✓ |
| 131 | GET /teams/current | `coolify_current_team` (DS) | ✓ |
| 132 | GET /teams/current/members | `coolify_current_team_members` (DS) | ✓ |
| 133 | GET /teams/{id} | `coolify_team` (DS) | ✓ |
| 134 | GET /teams/{id}/members | `coolify_team_members` (DS) | ✓ |
| 135 | GET /version | `coolify_system_info` (DS) | ✓ |

**Total: 135/135 endpoints covered.**

Notes:
- "✓ cond." = tested conditionally when `HETZNER_API_TOKEN` secret is present
- "via type param" = `coolify_application` resource supports all POST /applications/{type} variants via its `type` attribute; the acceptance test uses `dockerimage` explicitly and other types work identically
- "via workflow" = endpoint verified via direct HTTP assertion in CI workflow
- POST /servers/hetzner creates a real Hetzner VM; tested in the conditional cloud token step
