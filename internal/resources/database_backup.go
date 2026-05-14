package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewDatabaseBackupResource() resource.Resource {
	return &compositeResource{
		typeName:    "database_backup",
		displayName: "database backup",
		parentField: "database_uuid",
		listPath: func(p string) string {
			return fmt.Sprintf("/api/v1/databases/%s/backups", p)
		},
		createPath: func(p string) string {
			return fmt.Sprintf("/api/v1/databases/%s/backups", p)
		},
		updatePath: func(p, child string) string {
			return fmt.Sprintf("/api/v1/databases/%s/backups/%s", p, child)
		},
		deletePath: func(p, child string) string {
			return fmt.Sprintf("/api/v1/databases/%s/backups/%s", p, child)
		},
		fields: []resourceField{
			stringField("database_uuid", true, false, false),
			stringField("frequency", true, false, false),
			boolField("enabled", false, true, false),
			boolField("save_s3", false, true, false),
			stringField("s3_storage_uuid", false, true, false),
			boolField("dump_all", false, true, false),
			int64Field("timeout", false, true, false),
		},
	}
}
