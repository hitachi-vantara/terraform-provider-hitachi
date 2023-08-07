package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var SSInfoSchema = map[string]*schema.Schema{
	"storage_device_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a storage device id",
	},
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a storage serial number",
	},
	"storage_device_model": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a storage device model",
	},
	"dkc_micro_code_version": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a dkc micro code version of storage",
	},
	"management_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a storage management ip",
	},
	"svp_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a storage svp ip",
	},
	"controller1_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a storage controller1 ip",
	},
	"controller2_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a storage controller2 ip",
	},
	"free_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a free capacity of storage in MB",
	},
	"used_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a used capacity of storage in MB",
	},
	"total_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a total capacity of storage in MB",
	},
}

// first 4 are inputs
var StorageSystemSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage device is required",
	},
	// output
	"storage_system": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: SSInfoSchema,
		},
	},
}
