---
page_title: "coolify_application_storage Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Mounts a persistent volume on a Coolify application container, ensuring data survives redeployments.
---

# coolify_application_storage (Resource)

Mounts a persistent volume on a Coolify application container. Use this resource to persist directories that should survive container restarts and redeployments, such as upload directories or local caches.

## Example Usage

### Named Docker volume

```hcl
resource "coolify_application_storage" "uploads" {
  application_uuid = coolify_application.api.id
  type             = "volume"
  mount_path       = "/app/uploads"
  name             = "api-uploads"
}
```

### Bind mount from host path

```hcl
resource "coolify_application_storage" "data" {
  application_uuid = coolify_application.api.id
  type             = "bind"
  mount_path       = "/app/data"
  host_path        = "/data/api"
  is_readonly      = false
}
```

## Schema

### Required

- `application_uuid` (String) UUID of the application to attach the volume to.
- `mount_path` (String) Absolute path inside the container where the volume is mounted (e.g. `/app/uploads`).
- `type` (String) Storage type. Use `volume` for a named Docker volume or `bind` for a host bind mount.

### Optional

- `host_path` (String) Absolute path on the host to bind-mount into the container. Only used when `type` is `bind`.
- `is_readonly` (Boolean) When `true`, the volume is mounted read-only inside the container. Defaults to `false`.
- `name` (String) Name for the Docker volume. If omitted, Coolify generates a name.

### Read-Only

- `id` (String) Composite ID in the format `application_uuid/storage_uuid`.

## Import

Import an existing application storage volume using its composite ID:

```shell
terraform import coolify_application_storage.uploads <application-uuid>/<storage-uuid>
```
