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
		Description: "List of port attributes of the storage system",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"port_speed": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port speed of the storage system",
	},
	"loop_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Loop ID of the storage system",
	},
	"fabric_mode": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Fabric mode of the storage system",
	},
	"port_connection": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port connection of the storage system",
	},
	"lun_security_setting": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Lun security setting of the storage system",
	},
	"wwn": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "WWN of the port",
	},
	"port_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port mode of the storage system",
	},
}

var StoragePortsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Port ID of the storage system",
	},
	"include_detail_info": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Include detailed information (logins, portMode, loginHostNqn) for storage ports. When set to true, additional detailed fields will be populated.",
	},
	"port_type": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Port type filter. Allowed values: FIBRE, SCSI, ISCSI, NVME_TCP, ENAS, ESCON, FICON. If omitted, information about ports of all port types will be obtained. Cannot be used with port_id.",
	},
	"port_attributes": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Port attribute filter. Allowed values: TAR (Target port), MCU (Initiator port), RCU (RCU target port), ELUN (External port). If omitted, information about all port attributes will be obtained. Cannot be used with port_id.",
	},
	// output
	"total_port_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of ports on the storage system",
	},
	"ports": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Ports output",
		Elem: &schema.Resource{
			Schema: StoragePortSchema,
		},
	},
}
