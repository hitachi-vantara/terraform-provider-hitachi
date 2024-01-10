package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var InfraStoragePortSchema = map[string]*schema.Schema{
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port ID",
	},
	"type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port type",
	},
	"speed": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port speed of the storage device",
	},
	"resource_group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Resource Group Id of the Port",
	},
	"wwn": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "WWN of the port",
	},
	"resource_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource Group Id of the Port",
	},
	"attribute": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Attribute of the Port",
	},
	"connection_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Connection type of the Port",
	},
	"fabric_on": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether Fabric mode is on",
	},
	"mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port Mode",
	},
	"is_security_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether Security is enabled on the port",
	},
}

var InfraStoragePortsSchema = map[string]*schema.Schema{

	"storage_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Unique ID of the storage device",
	},

	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Serial number of storage",
	},

	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Port Id",
	},

	"total_port_count": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Total number of ports on the storage device",
	},

	// output
	"ports": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is ports output",
		Elem: &schema.Resource{
			Schema: InfraStoragePortSchema,
		},
	},
}
