package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DataComputeNodeSchema = map[string]*schema.Schema{
	"vss_block_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The  VSS block address of the storage server",
	},
	"compute_node_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the compute node to be fetched",
	},
	// output
	"compute_nodes": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "All the information if the selected compute node",
		Elem: &schema.Resource{
			Schema: ComputeNodeInfoSchema,
		},
	},
}

var ComputeNodeInfoSchema = map[string]*schema.Schema{

	"id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The id of the compute node",
	},
	"nickname": &schema.Schema{ //this is the actual nickname
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The name of the compute node",
	},
	"os_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The OS type of the compute node",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The total capacity of the compute node in megabytes",
	},
	"used_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The used capacity of the compute node in megabytes",
	},
	"number_of_volumes": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The number of volumes attached in the compute node",
	},
	"number_of_paths": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of ISCSI connection initiated to the compute node",
	},

	"paths": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Path/ISCSI connections details for the compute node",
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
					Description: "Port IDs of the ISCSI connection",
				},
			},
		},
	},
	"port_details": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Type:        schema.TypeList,
		Description: "Port Details of the ISCSI connection",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"port_id": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Port ID of the connection",
				},
				"iscsi_initiator": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Iscsi initiator name.",
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
	"vss_block_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the REST API server on Virtual Storage Software block.",
	},
	"compute_node_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the compute node to be created",
	},
	"os_type": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Type of the OS to be selected while creating the compute node like Linux, Windows,VmWare, etc.",
	},
	"iscsi_connection": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Details of ICSI connection to the compute node",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"iscsi_initiator": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "Name of the ICSI initiator",
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
