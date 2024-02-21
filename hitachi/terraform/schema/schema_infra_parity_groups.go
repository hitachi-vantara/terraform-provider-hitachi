package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var InfraParityGroupSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
	},
	"resource_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource  ID",
	},
	"parity_group_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ID of Parity Group",
	},
	"free_capacity": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Free Capacity",
	},
	"resource_group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "ID of Resource Group",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Total Capacity",
	},
	"ldev_ids": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of Ldevs",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"raid_level": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "RAID Level",
	},
	"drive_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Drive type",
	},
	"copyback_mode": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether callback mode is enabled",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status",
	},
	"is_pool_array_group": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether pool array group is enabled",
	},
	"is_accelerated_compression": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether accelerated_compression is enabled",
	},
	"is_encryption_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether encryption is enabled",
	},
}

var DataInfraParityGroupsSchema = map[string]*schema.Schema{

	"storage_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Unique ID of the storage device",
	},

	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Serial Number of the storage device",
	},

	"parity_group_ids": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "List of parity group IDs to fetch",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// output
	"parity_groups": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is ports output",
		Elem: &schema.Resource{
			Schema: InfraParityGroupSchema,
		},
	},
}
