package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ResourceStorageMaintenanceSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	// When true, invoke the appliance stop-format action. Optional on update;
	// if false/omitted during update, the provider will not call the action.
	"should_stop_all_volume_format": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "When true, invoke the appliance-wide stop-format action. If unset/false on update, the action is skipped.",
	},
}
