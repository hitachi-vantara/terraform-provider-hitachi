package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var StoragePortSchema = map[string]*schema.Schema{
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port ID",
	},
	"port_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port type",
	},
	"port_attributes": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of port attributes of the storage device",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"port_speed": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port speed of the storage device",
	},
	"loop_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Loop ID of the storage device",
	},
	"fabric_mode": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Fabric mode of the storage device",
	},
	"port_connection": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port connection of the storage device",
	},
	"lun_security_setting": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Lun security setting of the storage device",
	},
	"wwn": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "WWN of the port",
	},
}

var StoragePortsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Port ID of the storage device",
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
			Schema: StoragePortSchema,
		},
	},
}
