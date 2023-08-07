package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var LunInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a storage serial number",
	},
	"ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a ldev id of volume",
	},
	"clpr_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a clpr id of volume",
	},
	"emulation_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a emulation type of volume",
	},
	"num_ports": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of ports available on volume",
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
					Description: "It's a port id",
				},
				"hostgroup_id": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "It's a hostgroup id",
				},
				"hostgroup_name": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "It's a hostgroup name",
				},
				"lun_id": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "It's a lun id",
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
		Description: "List of attributes of volume",
	},
	"paritygroup_id": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "It's a parity group id of volume",
	},
	"label": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a label of volume",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's status of volume",
	},
	"mpblade_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a mpblade id of volume",
	},
	"ss_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a ss id of volume",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a pool id of volume",
	},
	"is_full_allocation_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "It checks whether full allocation is enabled on volume",
	},
	"resourcegroup_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a resource group id of volume",
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
		Description: "It checks whether alua is enabled on volume",
	},
	"naa_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a naa id of volume",
	},
	"free_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It shows free capacity of volume in MB",
	},
	"used_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It shows used capacity of volume in MB",
	},
	"total_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It shows total capacity of volume in MB",
	},
}

var DataLunSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage device is required",
	},
	"ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "It's a ldev id of lun",
	},
	// output
	"volume": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: LunInfoSchema,
		},
	},
}

var DataLunsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage device is required",
	},
	"start_ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "It's a start ldev id of lun",
	},
	"end_ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "It's a end ldev id of lun",
	},
	"undefined_ldev": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If set to true, returns not allocated luns else otherwise",
	},
	// output
	"volumes": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: LunInfoSchema,
		},
	},
}

var ResourceLunSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage device is required",
	},
	"ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Ldev Id to be created",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Pool id in which volume is to be created",
	},
	"paritygroup_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Parity group id in which volume is to be created",
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
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: LunInfoSchema,
		},
	},
}
