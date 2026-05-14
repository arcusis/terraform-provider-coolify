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

variable "coolify_endpoint" {}
variable "coolify_token"    { sensitive = true }

# ── Project CRUD ──────────────────────────────────────────────────────────────
resource "coolify_project" "test" {
  name        = "acceptance-test-project"
  description = "Created by provider acceptance test"
}

resource "coolify_project" "updated" {
  name        = "acceptance-test-project-renamed"
  description = "Updated description"
}

output "project_uuid"         { value = coolify_project.test.id }
output "project_uuid_updated" { value = coolify_project.updated.id }

# ── Data source: read back project ───────────────────────────────────────────
data "coolify_project" "readback" {
  uuid = coolify_project.test.id
}
output "project_name_readback" { value = data.coolify_project.readback.name }
