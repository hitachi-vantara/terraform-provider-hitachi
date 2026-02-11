package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ComputeNodesSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block system.",
	},
	// output
	"compute_node": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Compute node output",
		Elem: &schema.Resource{
			Schema: NodeInfoSchema,
		},
	},
}

var VolumeNodeInfoSchema = map[string]*schema.Schema{
	"id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ID of the compute node",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the compute node",
	},
	"os_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "OS type of the compute node",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity of the compute node",
	},
	"used_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity of the compute node",
	},
	"volume_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total volume count attached to the compute node",
	},
	"lun": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "LUN ID of the compute node",
	},
}

var NodeInfoSchema = map[string]*schema.Schema{
	"id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ID of the compute node",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the compute node",
	},
	"os_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Operating system type",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity",
	},
	"used_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity",
	},
	"volume_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Volume count",
	},
}
