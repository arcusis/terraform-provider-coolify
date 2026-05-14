---
page_title: "coolify_service_environment_variable Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Sets a single environment variable on a Coolify service (one-click stack). Use coolify_envs_bulk to manage multiple variables at once.
---

# coolify_service_environment_variable (Resource)

Sets a single environment variable on a Coolify service stack. For managing multiple variables in a single resource, use `coolify_envs_bulk`.

## Example Usage

```hcl
resource "coolify_service_environment_variable" "smtp_host" {
  service_uuid = coolify_service.blog.id
  key          = "SMTP_HOST"
  value        = "smtp.example.com"
}

resource "coolify_service_environment_variable" "smtp_password" {
  service_uuid = coolify_service.blog.id
  key          = "SMTP_PASSWORD"
  value        = var.smtp_password
}
```

## Schema

### Required

- `key` (String) Name of the environment variable (e.g. `SMTP_HOST`).
- `service_uuid` (String) UUID of the service this environment variable belongs to.
- `value` (String, Sensitive) Value of the environment variable.

### Optional

- `is_literal` (Boolean) When `true`, the value is treated as a literal string and variable interpolation (`$VAR`) is disabled. Defaults to `false`.
- `is_multiline` (Boolean) When `true`, indicates the value spans multiple lines. Defaults to `false`.
- `is_shown_once` (Boolean) When `true`, the value is hidden in the Coolify UI after being set. Defaults to `false`.

### Read-Only

- `id` (String) Coolify resource UUID.

## Import

Import an existing service environment variable using its UUID:

```shell
terraform import coolify_service_environment_variable.smtp_host <env-var-uuid>
```
