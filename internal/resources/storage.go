package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewApplicationStorageResource() resource.Resource {
	return storageResource("application", "applications", false)
}

func NewServiceStorageResource() resource.Resource {
	return storageResource("service", "services", true)
}

func NewDatabaseStorageResource() resource.Resource {
	return storageResource("database", "databases", false)
}

func storageResource(parentType, parentPath string, requiresResourceUUID bool) resource.Resource {
	fields := []resourceField{
		stringField(parentType+"_uuid", true, false, false),
		stringField("type", true, false, false),
		stringField("mount_path", true, false, false),
		stringField("name", false, true, false),
		stringField("host_path", false, true, false),
		boolField("is_readonly", false, true, false),
	}
	if requiresResourceUUID {
		// Optional — Coolify assigns to the first sub-service when omitted
		fields = append(fields, stringField("resource_uuid", false, true, false))
	}

	return &compositeResource{
		typeName:    parentType + "_storage",
		displayName: parentType + " storage",
		parentField: parentType + "_uuid",
		listPath: func(p string) string {
			return fmt.Sprintf("/api/v1/%s/%s/storages", parentPath, p)
		},
		createPath: func(p string) string {
			return fmt.Sprintf("/api/v1/%s/%s/storages", parentPath, p)
		},
		updatePath: func(p, child string) string {
			return fmt.Sprintf("/api/v1/%s/%s/storages/%s", parentPath, p, child)
		},
		deletePath: func(p, child string) string {
			return fmt.Sprintf("/api/v1/%s/%s/storages/%s", parentPath, p, child)
		},
		fields: fields,
	}
}
