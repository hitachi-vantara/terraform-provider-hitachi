package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var SSInfoSchema = map[string]*schema.Schema{
	"storage_device_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage system ID",
	},
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of the storage system",
	},
	"storage_device_model": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage system model",
	},
	"dkc_micro_code_version": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "DKC micro code version of the storage system",
	},
	"management_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage system management IP address",
	},
	"svp_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage system SVP IP address",
	},
	"controller1_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage system controller1 IP address",
	},
	"controller2_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage system controller2 IP address",
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
		Description: "Serial number of the storage system",
	},
	// output
	"storage_system": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Storage system output",
		Elem: &schema.Resource{
			Schema: SSInfoSchema,
		},
	},
}
