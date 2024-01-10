package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var InfraStoragePoolSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
	},
	"resource_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource Group Id of the Pool",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Pool Id",
	},
	"ldev_ids": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of Ldev Ids",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the Pool",
	},
	"depletion_threshold_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Depletion Threshold Rate",
	},
	"dp_volumes": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "DP Volumes of the Pool",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"logical_unit_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Logical Unit Id",
				},
				"size": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Size",
				},
			},
		},
	},
	"free_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free Capacity",
	},
	"free_capacity_in_units": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Free Capacity With Units",
	},
	"replication_depletion_alert_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Replication Depletion Alert Rate",
	},
	"replication_usage_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Replication Usage Rate",
	},
	"resource_group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Resource Group Id of the Pool",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status",
	},
	"subscription_limit_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Subscription Limit Rate",
	},
	"subscription_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Subscription Rate",
	},
	"subscription_warning_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Subscription Warning Rate",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total Capacity",
	},
	"total_capacity_in_unit": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Total Capacity In Unit",
	},
	"type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port type",
	},
	"virtual_volume_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Virtual Volume Count",
	},
	"warning_threshold_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Warning Threshold Rate",
	},
	"deduplication_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Deduplication Enabled",
	},
}

var DataInfraStoragePoolsSchema = map[string]*schema.Schema{

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
	"pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Name of the storage pool",
	},

	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Id of the storage pool",
	},

	// output
	"storage_pools": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is pools output",
		Elem: &schema.Resource{
			Schema: InfraStoragePoolSchema,
		},
	},
}
