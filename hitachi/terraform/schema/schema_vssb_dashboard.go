package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var VssbDashboardInfoSchema = map[string]*schema.Schema{

	"health_status": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Type:        schema.TypeList,
		Description: "Displays the health status",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Object type.",
				},
				"status": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Health status.",
				},
			},
		},
	},
	"fault_domain_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of fault domains.",
	},
	"volume_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of volumes.",
	},
	"compute_node_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of compute nodes",
	},
	"compute_port_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of compute ports",
	},
	"drive_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of drives",
	},
	"storage_node_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of storage nodes",
	},
	"storage_pool_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of storage nodes",
	},
	"total_capacity_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total pool capacity in MiB",
	},
	"used_capacity_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used pool capacity in MiB",
	},
	"free_capacity_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free pool capacity in MiB",
	},
	"total_capacity_gb": &schema.Schema{
		Type:        schema.TypeFloat,
		Computed:    true,
		Description: "Total pool capacity in MiB",
	},
	"used_capacity_gb": &schema.Schema{
		Type:        schema.TypeFloat,
		Computed:    true,
		Description: "Used pool capacity in MiB",
	},
	"free_capacity_gb": &schema.Schema{
		Type:        schema.TypeFloat,
		Computed:    true,
		Description: "Free pool capacity in MiB",
	},
	"total_efficiency": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Indicates the effect of volume creation and snapshot functions on capacity consumption.",
	},
	"data_reduction": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Ratio of the data capacity before and after data reduction (unit: %).",
	},
}
var DataVssbDashboardSchema = map[string]*schema.Schema{

	"vss_block_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the REST API server on Virtual Storage Software block.",
	},
	// output
	"dashboard_info": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: VssbDashboardInfoSchema,
		},
	},
}
