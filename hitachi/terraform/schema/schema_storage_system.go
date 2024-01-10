package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var SSInfoSchema = map[string]*schema.Schema{
	"storage_device_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage device ID",
	},
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
	},
	"storage_device_model": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage device model",
	},
	"dkc_micro_code_version": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "DKC micro code version of storage",
	},
	"management_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage management ip",
	},
	"svp_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage svp IP",
	},
	"controller1_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage controller1 IP",
	},
	"controller2_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage controller2 IP",
	},
	"free_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free capacity in MB",
	},
	"used_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity in MB",
	},
	"total_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity in MB",
	},
}

// first 4 are inputs
var StorageSystemSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage",
	},
	// output
	"storage_system": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is storage system output",
		Elem: &schema.Resource{
			Schema: SSInfoSchema,
		},
	},
}
