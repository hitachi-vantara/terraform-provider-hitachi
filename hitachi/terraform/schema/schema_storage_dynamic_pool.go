package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DynamicPoolInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of the storage device",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Poop id of the storage device",
	},
	"pool_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool status of the storage device",
	},
	"used_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity rate of pool",
	},
	"used_physical_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used physical capacity rate of pool",
	},
	"snapshot_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Snapshot count of pool",
	},
	"pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "pool name",
	},
	"available_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available volume capacity on pool",
	},
	"available_physical_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available physical volume capacity on pool",
	},
	"total_pool_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total pool capacity on pool",
	},
	"total_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total physical capacity on pool",
	},
	"num_of_ldevs": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of ldevs on pool",
	},
	"first_ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "First ldev id on pool",
	},
	"warning_threshold": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Warning threshold on pool",
	},
	"depletion_threshold": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Depletion threshold on pool",
	},
	"virtual_volume_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Virtual volume capacity rate on pool",
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
		Description: "Total number of located volume count on pool",
	},
	"total_located_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of located capacity on pool",
	},
	"blocking_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Blocking mode of pool",
	},
	"total_reserved_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of reserved capacity on pool",
	},
	"reserved_volume_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of reserved volume count on pool",
	},
	"pool_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool type",
	},
	"duplication_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Duplication number of pool",
	},
	"duplication_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Duplication rate of pool",
	},
	"data_reduction_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Data reduction rate of pool",
	},
	"snapshot_used_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Snapshot used capacity of pool",
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
		Description: "Serial number of storage device is required",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Pool id of the storage device",
	},
	// output
	"dynamic_pools": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: DynamicPoolInfoSchema,
		},
	},
}
