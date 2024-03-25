package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DynamicPoolInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Pool ID of storage",
	},
	"pool_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool status of storage",
	},
	"used_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity rate",
	},
	"used_physical_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used physical capacity rate",
	},
	"snapshot_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Snapshot count",
	},
	"pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool name",
	},
	"available_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available volume capacity",
	},
	"available_physical_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available physical volume capacity",
	},
	"total_pool_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total pool capacity",
	},
	"total_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total physical capacity",
	},
	"num_of_ldevs": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of ldevs",
	},
	"first_ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "First ldev ID",
	},
	"warning_threshold": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Warning threshold",
	},
	"depletion_threshold": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Depletion threshold",
	},
	"virtual_volume_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Virtual volume capacity rate",
	},
	"is_mainframe": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Is mainframe pool",
	},
	"is_shrinking": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Is shrinking pool",
	},
	"located_volume_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of located volume count",
	},
	"total_located_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of located capacity",
	},
	"blocking_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Blocking mode of pool",
	},
	"total_reserved_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of reserved capacity",
	},
	"reserved_volume_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of reserved volume count",
	},
	"pool_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool type",
	},
	"duplication_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Duplication number",
	},
	"duplication_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Duplication rate",
	},
	"data_reduction_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Data reduction rate",
	},
	"snapshot_used_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Snapshot used capacity",
	},
	"suspend_snapshot": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks if suspend snapshot",
	},
}

var DataDynamicPoolsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Pool ID of the storage",
	},
	"pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Pool Name of the storage",
	},
	// output
	"dynamic_pools": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is dynamic pools output",
		Elem: &schema.Resource{
			Schema: DynamicPoolInfoSchema,
		},
	},
}
