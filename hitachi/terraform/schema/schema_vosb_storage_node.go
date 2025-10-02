package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)


var DataVssbStorageNodeSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block.",
	},
	"node_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Name of the storage node",
		Default:     "",
	},
	// output
	"nodes": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Storage nodes output",
		Elem: &schema.Resource{
			Schema: VssbStorageNodeInfoSchema,
		},
	},
}

var VssbStorageNodeInfoSchema = map[string]*schema.Schema{
	"id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ID of node",
	},
	"bios_uuid": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Bios Uuid of node",
	},
	"fault_domain_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Fault Domain Id of node",
	},
	"fault_domain_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Fault Doma Name of node",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of node",
	},
	"cluster_role": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "cluster Role of node",
	},
	"drive_data_relocation_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Drive Data Relocation Status of node",
	},
	"control_port_ipv4_address": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Control Port Ipv4 Address of node",
	},
	"internode_port_ipv4_address": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Internode Port Ipv4 Address of node",
	},
	"software_version": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Software Version of node",
	},
	"serial_number": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Serial Number of node",
	},
	"memory": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Memory size of node",
	},
	"availability_zone_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Availability Zone Id of node",
	},
	"model_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "model name of node",
	},
	"protection_domain_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Protection Domain ID of node",
	},
	"status_summary": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "status summary of node",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status of node",
	},
	"insufficient_resources_for_rebuild_capacity": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		// MaxItems:    1,
		Description: "Insufficient resources for rebuild capacity of node",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"capacity_of_drive": {
					Optional:    true,
					Type:        schema.TypeInt,
					Default:     -1,
					Description: "Capacity Of Drive",
				},
				"number_of_drives": {
					Optional:    true,
					Type:        schema.TypeInt,
					Description: "Number Of Drives",
				},
			},
		},
	},
	"rebuildable_resources": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		// MaxItems:    1,
		Description: "Rebuildable resources information",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"number_of_drives": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Number of drives of node",
				},
			},
		},
	},
}

// ResourceVssbStorageNodeSchema defines the schema for the VSSB storage node resource.
var ResourceVssbStorageNodeSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block.",
	},
	"configuration_file": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Configuration File",
	},
	"exported_configuration_file": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Configuration file exported to add storagenodes",
	},
	"vm_configuration_file_s3_uri": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "S3 URI for VM configuration file (required for AWS deployments)",
	},
	"setup_user_password": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Setup User Password",
	},
	"expected_cloud_provider": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "baremetal",
		Description: `Specifies the expected cloud provider type. Valid values: "google", "azure", "baremetal".

	- Used to validate combinations of inputs based on the deployment environment.
	- If set to "google" or "azure", specific parameters may be required for certain operations.
	- If set to "baremetal" (default), other cloud-specific inputs are ignored.
	- Note: The actual cloud provider is determined by the VSP One SDS Block system at the "vosb_address" endpoint.
	If there's a mismatch, the request still proceeds and behaves according to the actual environment.`,
		ValidateFunc: validation.StringInSlice([]string{
			"google", "azure", "aws", "baremetal",
		}, false),
	},
	// "node_name": &schema.Schema{
	// 	Type:        schema.TypeString,
	// 	Required:    false,
	// 	Description: "Storage node name to be added",
	// },		
	// output
	"storage_nodes": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "Additional information about the compute node",
		Elem: &schema.Resource{
			Schema: VssbStorageNodeInfoSchema,
		},
	},
}
