---
page_title: "coolify_service_storage Resource - terraform-provider-coolify"
subcategory: ""
description: |-
  Mounts a persistent volume on a container within a Coolify service stack, ensuring data survives redeployments.
---

# coolify_service_storage (Resource)

Mounts a persistent volume on a container within a Coolify service stack. Use this resource to persist directories such as media uploads, configuration files, or generated content that must survive service redeployments.

## Example Usage

### Named Docker volume for Ghost content

```hcl
resource "coolify_service_storage" "content" {
  service_uuid  = coolify_service.blog.id
  resource_uuid = coolify_service.blog.id
  type          = "volume"
  mount_path    = "/var/lib/ghost/content"
  name          = "ghost-content"
}
```

### Bind mount

```hcl
resource "coolify_service_storage" "uploads" {
  service_uuid  = coolify_service.blog.id
  resource_uuid = coolify_service.blog.id
  type          = "bind"
  mount_path    = "/app/uploads"
  host_path     = "/data/blog-uploads"
  is_readonly   = false
}
```

## Schema

### Required

- `mount_path` (String) Absolute path inside the container where the volume is mounted.
- `resource_uuid` (String) UUID of the specific container resource within the service to mount the volume on.
- `service_uuid` (String) UUID of the service this storage volume belongs to.
- `type` (String) Storage type. Use `volume` for a named Docker volume or `bind` for a host bind mount.

### Optional

- `host_path` (String) Absolute path on the host to bind-mount into the container. Only used when `type` is `bind`.
- `is_readonly` (Boolean) When `true`, the volume is mounted read-only inside the container. Defaults to `false`.
- `name` (String) Name for the Docker volume. If omitted, Coolify generates a name.

### Read-Only

- `id` (String) Composite ID in the format `service_uuid/storage_uuid`.

## Import

Import an existing service storage volume using its composite ID:

```shell
terraform import coolify_service_storage.content <service-uuid>/<storage-uuid>
```
