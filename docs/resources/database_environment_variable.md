---
page_title: "coolify_database_environment_variable Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Sets a single environment variable on a Coolify-managed database.
---

# coolify_database_environment_variable (Resource)

Sets a single environment variable on a Coolify-managed database. Use `coolify_envs_bulk` instead if you need to manage multiple variables at once.

## Example Usage

```hcl
resource "coolify_database_environment_variable" "extra_config" {
  database_uuid = coolify_database_postgresql.main.id
  key           = "EXTRA_POSTGRES_OPTION"
  value         = "some-value"
}
```

### Multiline value

```hcl
resource "coolify_database_environment_variable" "ssl_cert" {
  database_uuid = coolify_database_postgresql.main.id
  key           = "SSL_CERT"
  value         = file("${path.module}/cert.pem")
  is_multiline  = true
}
```

## Schema

### Required

- `database_uuid` (String) UUID of the database this environment variable belongs to.
- `key` (String) Name of the environment variable (e.g. `MY_VAR`).
- `value` (String, Sensitive) Value of the environment variable.

### Optional

- `is_literal` (Boolean) When `true`, the value is treated as a literal string and variable interpolation (`$VAR`) is disabled. Defaults to `false`.
- `is_multiline` (Boolean) When `true`, indicates the value spans multiple lines. Defaults to `false`.
- `is_shown_once` (Boolean) When `true`, the value is hidden in the Coolify UI after being set. Defaults to `false`.

### Read-Only

- `id` (String) Coolify resource UUID.

## Import

Import an existing database environment variable using its UUID:

```shell
terraform import coolify_database_environment_variable.extra_config <env-var-uuid>
```
