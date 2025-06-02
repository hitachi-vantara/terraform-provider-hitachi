package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ParityGroupsInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
	},
	"parity_group_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Parity group ID",
	},
	"num_of_ldevs": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of ldevs in parity group",
	},
	"used_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity rate of parity group",
	},
	"available_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available volume capacity of parity group",
	},
	"raid_level": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Raid level of parity group",
	},
	"raid_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Raid type of parity group",
	},
	"clpr_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Clpr ID of parity group",
	},
	"drive_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Drive type of parity group",
	},
	"drive_type_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Drive type name of parity group",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity of parity group",
	},
	"physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Physical capacity of parity group",
	},
	"available_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available physical capacity of parity group",
	},
	"is_accelerated_compression_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks if is acclerated compression enabled of parity group",
	},
	"available_volume_capacity_in_kb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available volume capacity of parity group in kb",
	},
}

var DataParityGroupSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage",
	},
	"parity_group_ids": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
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
		Description: "This is parity groups output",
		Elem: &schema.Resource{
			Schema: ParityGroupsInfoSchema,
		},
	},
}

var DataParityGroupsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage",
	},
	"parity_group_ids": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
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
		Description: "This is parity groups output",
		Elem: &schema.Resource{
			Schema: ParityGroupsInfoSchema,
		},
	},
}
