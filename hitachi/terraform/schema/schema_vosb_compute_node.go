package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DataComputeNodeSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block system.",
	},
	"compute_node_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The name to be retrieved",
	},
	// output
	"compute_nodes": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "All the information of the selected compute node",
		Elem: &schema.Resource{
			Schema: ComputeNodeInfoSchema,
		},
	},
}

var ComputeNodeInfoSchema = map[string]*schema.Schema{

	"id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ID",
	},
	"nickname": &schema.Schema{ //this is the actual nickname
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Nickname",
	},
	"os_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "OS type",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity in MB",
	},
	"used_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity in MB",
	},
	"number_of_volumes": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The number of volumes attached to the compute node",
	},
	"number_of_paths": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of iSCSI connection initiated to the compute node",
	},

	"paths": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Path/iSCSI connections details",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"protocol": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Type of Protocol of the attached connection",
				},
				"hba_name": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Name of the HBA connection",
				},
				"port_ids": &schema.Schema{
					Computed: true,
					Type:     schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Port IDs of the iSCSI connection",
				},
			},
		},
	},
	"port_details": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Type:        schema.TypeList,
		Description: "Port Details of the iSCSI connection",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"port_id": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Port ID of the connection",
				},
				"target_port_identifier": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "WWN (for FC) or iSCSI name (for iSCSI) of the allocation destination compute port of the target operation",
				},
				"port_name": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Port name of the connection",
				},
			},
		},
	},
}

var ResourceVssbStorageComputeNodeSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of VSP One SDS Block.",
	},
	"compute_node_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name to be created",
	},
	"os_type": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Type of the OS to be selected like Linux, Windows, VmWare, etc., it is required while creating compute node.",
	},
	"iscsi_connection": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Details of iSCSI connection to the compute node",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"iscsi_initiator": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "Name of the iSCSI initiator",
				},
				"port_names": &schema.Schema{
					Optional:    true,
					Type:        schema.TypeList,
					Description: "List of port names to connect to the compute node",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	},
	"fc_connection": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Details of Fiber channel connections",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host_wwn": {
					Optional:    true,
					Description: "Host WWN Names which need to be attached to the compute node",
					Type:        schema.TypeString,
				},
			},
		},
	},
	// output
	"compute_nodes": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "Additional information about the compute node",
		Elem: &schema.Resource{
			Schema: ComputeNodeInfoSchema,
		},
	},
}
