terraform {
  required_providers {
    coolify = {
      source  = "arcusis/coolify"
      version = "~> 0.1"
    }
  }
}

provider "coolify" {
  endpoint = var.coolify_endpoint
  token    = var.coolify_token
}

variable "coolify_endpoint" {
  type = string
}

variable "coolify_token" {
  type      = string
  sensitive = true
}

resource "coolify_project" "example" {
  name        = "example"
  description = "Managed by Terraform"
}

resource "coolify_private_key" "example" {
  name        = "example"
  description = "Managed by Terraform"
  private_key = var.server_private_key
}

variable "server_private_key" {
  type      = string
  sensitive = true
}

resource "coolify_server" "example" {
  name             = "example"
  ip               = "203.0.113.10"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.example.id
}

resource "coolify_application" "example" {
  type             = "public"
  project_uuid     = coolify_project.example.id
  server_uuid      = coolify_server.example.id
  environment_name = "production"
  ports_exposes    = "3000"
  git_repository   = "https://github.com/example/app"
  git_branch       = "main"
  build_pack       = "nixpacks"
  name             = "example-app"
}

resource "coolify_environment_variable" "example" {
  application_uuid = coolify_application.example.id
  key              = "EXAMPLE"
  value            = "true"
}
