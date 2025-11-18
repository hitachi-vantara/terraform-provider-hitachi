package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var LunInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of the storage system",
	},
	"ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "LDEV ID",
	},
	"clpr_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "CLPR ID",
	},
	"emulation_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Emulation type",
	},
	"num_ports": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of ports available on the volume",
	},
	"ports": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"port_id": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Port ID",
				},
				"hostgroup_id": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Host group ID",
				},
				"hostgroup_name": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Host group name",
				},
				"lun_id": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "LUN ID",
				},
			},
		},
	},
	"attributes": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "List of attributes of the volume",
	},
	"paritygroup_id": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "Parity group ID",
	},
	"label": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Label",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status",
	},
	"mpblade_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "MP blade ID",
	},
	"ss_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "SS ID",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Pool ID",
	},
	"is_full_allocation_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks whether full allocation is enabled on the volume",
	},
	"resourcegroup_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Resource group ID of the volume",
	},
	/*
		"data_reduction_mode": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "It's data reduction mode of volume",
		},
	*/
	"is_alua_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks whether ALUA is enabled on the volume",
	},
	"naa_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "NAA ID",
	},
	"free_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free capacity in MB",
	},
	"used_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity in MB",
	},
	"total_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity in MB",
	},
}

var DataLunSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ValidateFunc: validation.IntBetween(0, 65535),
		Description: "LDEV ID",
	},
	// output
	"volume": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Volume output",
		Elem: &schema.Resource{
			Schema: LunInfoSchema,
		},
	},
}

var DataLunsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"start_ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ValidateFunc: validation.IntBetween(0, 65535),
		Description: "Starting LDEV ID",
	},
	"end_ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ValidateFunc: validation.IntBetween(0, 65535),
		Description: "Ending LDEV ID",
	},
	"undefined_ldev": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If set to true, returns not allocated LUNs",
	},
	// output
	"volumes": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Volumes output",
		Elem: &schema.Resource{
			Schema: LunInfoSchema,
		},
	},
}

var ResourceLunSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "LDEV ID",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Default: -1 ,
		Description: "Pool ID in which volume is to be created. **One of `paritygroup_id`, `pool_id`, or `pool_name` must be specified.**",
	},
	"pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Pool Name in which volume is to be created. **One of `paritygroup_id`, `pool_id`, or `pool_name` must be specified.**",
	},
	"paritygroup_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Parity group ID in which volume is to be created. **One of `paritygroup_id`, `pool_id`, or `pool_name` must be specified.**",
	},
	"size_gb": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Size of volume to be created in GB",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Name of volume to be created",
	},
	//Remove dedup from this version
	/*
		"dedup_mode": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Dedup mode of volume to be created",
		},
	*/

	// output
	"volume": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "Volume output",
		Elem: &schema.Resource{
			Schema: LunInfoSchema,
		},
	},
}
