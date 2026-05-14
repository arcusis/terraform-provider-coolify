---
page_title: "coolify_environment_variable Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Sets a single environment variable on a Coolify application. Use coolify_envs_bulk to manage multiple variables at once.
---

# coolify_environment_variable (Resource)

Sets a single environment variable on a Coolify application. For managing multiple variables in a single resource, use `coolify_envs_bulk`.

## Example Usage

```hcl
resource "coolify_environment_variable" "db_url" {
  application_uuid = coolify_application.api.id
  key              = "DATABASE_URL"
  value            = "postgresql://user:${var.db_password}@db:5432/appdb"
}

resource "coolify_environment_variable" "node_env" {
  application_uuid = coolify_application.api.id
  key              = "NODE_ENV"
  value            = "production"
}
```

### Preview environment variable

```hcl
resource "coolify_environment_variable" "preview_api_url" {
  application_uuid = coolify_application.api.id
  key              = "API_URL"
  value            = "https://preview.example.com"
  is_preview       = true
}
```

## Schema

### Required

- `application_uuid` (String) UUID of the application this environment variable belongs to.
- `key` (String) Name of the environment variable (e.g. `DATABASE_URL`).
- `value` (String, Sensitive) Value of the environment variable.

### Optional

- `is_literal` (Boolean) When `true`, the value is treated as a literal string and variable interpolation (`$VAR`) is disabled. Defaults to `false`.
- `is_multiline` (Boolean) When `true`, indicates the value spans multiple lines (e.g. a PEM certificate). Defaults to `false`.
- `is_preview` (Boolean) When `true`, this variable is only injected into PR preview deployments, not the main deployment. Defaults to `false`.
- `is_shown_once` (Boolean) When `true`, the value is hidden in the Coolify UI after being set. Defaults to `false`.

### Read-Only

- `id` (String) Coolify resource UUID.

## Import

Import an existing application environment variable using its UUID:

```shell
terraform import coolify_environment_variable.db_url <env-var-uuid>
```
