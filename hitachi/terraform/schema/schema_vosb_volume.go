package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DataVolumeSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block and Cloud system.",
	},
	"compute_node_name": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	// output
	"volumes": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: VolumeInfoSchema,
		},
	},
}

var ResourceVolumeSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block and Cloud system.",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the volume server",
	},
	"storage_pool": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Storage pool name of the storage system. Required for the create operation.",
	},
	"capacity_gb": &schema.Schema{
		Type:        schema.TypeFloat,
		Optional:    true,
		Description: "Capacity of the volume to be created in Gigabytes. Required for the create operation.",
	},
	"compute_nodes": &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true, // needed to get rid of '/* of string */' in the output
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "List of compute nodes to be attached to the volume. To remove all the nodes from the volume declare compute_nodes = []",
	},
	"nick_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Nickname of the volume",
	},
	// output
	"volume": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Output information about the volume",
		Elem: &schema.Resource{
			Schema: VolumeInfoSchema,
		},
	},
}

var VolumeNodeSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block and Cloud system.",
	},
	"volume_name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	// output
	"volume": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: VolumeInfoSchema,
		},
	},
}

var VolumeInfoSchema = map[string]*schema.Schema{
	"saving_effects": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Type:        schema.TypeList,
		Description: "Volumes saving effects information",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"system_data_capacity": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "System data capacity",
				},
				"pre_capacity_data_reduction_without_system_data": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Pre-capacity data reduction without system data",
				},
				"post_capacity_data_reduction": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Post-capacity data reduction without system data",
				},
			},
		},
	},

	"id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ID of the volume resource",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the created volume resource",
	},
	"nick_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Nickname of the created volume resource",
	},
	"volume_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Volume number of the created volume resource",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool ID of the created volume resource",
	},
	"pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool Name of the created volume resource",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity of the created volume resource",
	},
	"used_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity of the created volume resource",
	},

	"number_of_connecting_servers": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of connected servers to create the volume resource",
	},
	"number_of_snapshots": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of snapshots to create the volume resource",
	},
	"protection_domain_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Protection domain ID of the volume resource",
	},
	"full_allocated": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Allocation status of the volume resource",
	},
	"volume_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Volume type of the volume resource",
	},
	"status_summary": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status summary of the volume resource",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status of the volume resource",
	},
	"storage_controller_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage controller ID of the volume",
	},
	"snapshot_attribute": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Snapshot attribute of the volume",
	},

	"snapshot_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Snapshot status of the volume resource",
	},
	"saving_setting": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Setting of the volume resource",
	},
	"saving_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Saving mode of the volume resource",
	},
	"data_reduction_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Data reduction status of the volume resource",
	},
	"data_reduction_progress_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Data reduction progress rate of the volume",
	},
	"compute_nodes": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "List of the compute nodes attached to the volume",
		Elem: &schema.Resource{
			Schema: VolumeNodeInfoSchema,
		},
	},
}
