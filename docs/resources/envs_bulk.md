---
page_title: "coolify_envs_bulk Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Sets the complete set of environment variables for an application, service, or database in a single operation, replacing all existing variables on each apply.
---

# coolify_envs_bulk (Resource)

Sets the complete set of environment variables for an application, service, or database in a single resource. Each `apply` replaces all existing environment variables for the target resource with the provided `variables` map. This is the recommended approach when you need to manage many variables together.

> **Note:** Because this resource replaces all variables on each apply, any variables set outside of Terraform (e.g. via the Coolify UI) will be removed on the next `terraform apply`.

## Example Usage

### Application environment variables

```hcl
resource "coolify_envs_bulk" "app_vars" {
  resource_type = "application"
  resource_uuid = coolify_application.api.id
  variables = {
    DATABASE_URL = "postgresql://user:${var.db_password}@db:5432/appdb"
    REDIS_URL    = "redis://:${var.redis_password}@redis:6379"
    NODE_ENV     = "production"
    PORT         = "3000"
  }
}
```

### Service environment variables

```hcl
resource "coolify_envs_bulk" "blog_vars" {
  resource_type = "service"
  resource_uuid = coolify_service.blog.id
  variables = {
    GHOST_URL       = "https://blog.example.com"
    SMTP_HOST       = "smtp.example.com"
    SMTP_PASSWORD   = var.smtp_password
  }
}
```

### Database environment variables

```hcl
resource "coolify_envs_bulk" "db_vars" {
  resource_type = "database"
  resource_uuid = coolify_database_postgresql.main.id
  variables = {
    EXTRA_OPTION = "value"
  }
}
```

## Schema

### Required

- `resource_type` (String) Type of the target resource. Accepted values: `application`, `service`, `database`.
- `resource_uuid` (String) UUID of the target resource.
- `variables` (Map of String) Map of environment variable names to their values. This map replaces all existing environment variables on each apply.

### Read-Only

- `id` (String) Coolify resource UUID.

## Import

Import an existing bulk env configuration using its resource UUID:

```shell
terraform import coolify_envs_bulk.app_vars <resource-uuid>
```
