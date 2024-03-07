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

var PartnerInfraStoragePortInfoSchema = map[string]*schema.Schema{
	"resource_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource Id",
	},
	"type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Types of the resource",
	},
	"storage_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage Id",
	},
	"entitlement_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Entitlement status of the storage port",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Port Id",
	},
	"port_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Type of the port",
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
	/*
		"partner_id": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "partner Id id  of the storage device",
		},
		"subscriber_id": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "subscriber Id id  of the storage device",
		},
	*/
}

var DataInfraStoragePortsSchema = map[string]*schema.Schema{

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

	//output
	"partner_ports": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is partners volumes output",
		Elem: &schema.Resource{
			Schema: PartnerInfraStoragePortInfoSchema,
		},
	},
}
