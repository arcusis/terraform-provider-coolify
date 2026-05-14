package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewApplicationScheduledTaskResource() resource.Resource {
	return scheduledTaskResource("application", "applications")
}

func NewServiceScheduledTaskResource() resource.Resource {
	return scheduledTaskResource("service", "services")
}

func scheduledTaskResource(parentType, parentPath string) resource.Resource {
	return &compositeResource{
		typeName:    parentType + "_scheduled_task",
		displayName: parentType + " scheduled task",
		parentField: parentType + "_uuid",
		listPath: func(p string) string {
			return fmt.Sprintf("/api/v1/%s/%s/scheduled-tasks", parentPath, p)
		},
		createPath: func(p string) string {
			return fmt.Sprintf("/api/v1/%s/%s/scheduled-tasks", parentPath, p)
		},
		updatePath: func(p, child string) string {
			return fmt.Sprintf("/api/v1/%s/%s/scheduled-tasks/%s", parentPath, p, child)
		},
		deletePath: func(p, child string) string {
			return fmt.Sprintf("/api/v1/%s/%s/scheduled-tasks/%s", parentPath, p, child)
		},
		fields: []resourceField{
			stringField(parentType+"_uuid", true, false, false),
			stringField("name", true, false, false),
			stringField("command", true, false, false),
			stringField("frequency", true, false, false),
			stringField("container", false, true, false),
			int64Field("timeout", false, true, false),
			boolField("enabled", false, true, false),
		},
	}
}
