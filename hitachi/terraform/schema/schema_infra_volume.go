package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var InfraVolumeInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
	},
	"resource_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource Id",
	},
	"deduplication_compression_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Deduplication Compression Mode",
	},
	"emulation_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Emulation type",
	},
	"format_or_shred_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Format Or Shred Rate",
	},
	"ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Ldev ID of lun",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name",
	},
	"parity_group_id": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "Parity group ID",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Pool ID",
	},
	"resource_group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Resource group ID of volume",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity",
	},
	"used_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity",
	},
	"virtual_storage_device_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Virtual Storage Device Id",
	},
	"stripe_size": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Stripe Size",
	},
	"type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Type",
	},
	"path_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Path Count",
	},
	"provision_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Provision Type",
	},
	"is_command_device": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks whether it is Command Device",
	},
	"logical_unit_id_hex_format": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Logical Unit Id in Hex Format",
	},
	"virtual_logical_unit_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Virtual Logical Unit Id",
	},
	"naa_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "NAA ID",
	},
	"dedup_compression_progress": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Dedup Compression Progress",
	},
	"dedup_compression_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Dedup Compression Status",
	},
	"is_alua_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks whether alua is enabled on volume",
	},
	"is_dynamic_pool_volume": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks whether it is a dynamic pool volume",
	},
	"is_journal_pool_volume": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks whether it is a journal pool volume",
	},
	"is_pool_volume": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks whether it is a pool volume",
	},
	"pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool Name",
	},
	"quorum_disk_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Quorum Disk Id",
	},
	"is_in_gad_pair": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks whether it is in a Gad Pair",
	},
	"is_vvol": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks whether it is a VVol",
	},

	/*
	 * The following properties came from 2.0
	 * Need to figure out how to find these properties in 2.5
	 */
	"clpr_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "clpr ID",
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
					Description: "Port ID",
				},
				"hostgroup_id": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "HostGroup ID",
				},
				"hostgroup_name": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "HostGroup name",
				},
				"lun_id": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Lun ID",
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
	"label": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Label",
	},

	"mpblade_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Mpblade ID",
	},
	"ss_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "SS ID",
	},
	"is_full_allocation_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "It checks whether full allocation is enabled on volume",
	},

	/*
		"data_reduction_mode": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "It's data reduction mode of volume",
		},
	*/

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

var PartnerInfraVolumeInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
	},
	"resource_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource Id",
	},
	"storage_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage Id",
	},
	"type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Type of resource",
	},
	"entitlement_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Entitlement Status of the volume",
	},
	"partner_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "partner Id id  of the volume",
	},
	"subscriber_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "subscriber Id id  of the volume",
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

	"ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "ldev id of the volume",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Pool Id of the volume",
	},
	"pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool name of the volume",
	},
}

var DataInfraVolumeSchema = map[string]*schema.Schema{
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
	"ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Ldev ID of lun",
	},
	"subscriber_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Id of the subscriber which is attached to the volume",
	},
	// output
	"volume": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is volume output",
		Elem: &schema.Resource{
			Schema: InfraVolumeInfoSchema,
		},
	},
	"subscriber_volume": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is subscriber's volume output",
		Elem: &schema.Resource{
			Schema: PartnerInfraVolumeInfoSchema,
		},
	},
}

var DataInfraVolumesSchema = map[string]*schema.Schema{
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
	"start_ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Start ldev ID of lun",
	},
	"end_ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "End ldev ID of lun",
	},
	"undefined_ldev": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If set to true, returns not allocated luns else otherwise",
	},
	"subscriber_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Id of the subscriber which is attached to the volume",
	},
	// output
	"volumes": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is volumes output",
		Elem: &schema.Resource{
			Schema: InfraVolumeInfoSchema,
		},
	},

	//output
	"partner_volumes": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is partners volumes output",
		Elem: &schema.Resource{
			Schema: PartnerInfraVolumeInfoSchema,
		},
	},
}

var ResourceInfraVolumeSchema = map[string]*schema.Schema{

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
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "name number of volume to be created",
	},

	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Pool ID in which volume is to be created",
	},
	"lun_id": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Lun ID in which volume is to be created",
	},
	"resource_group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Resource group id which volume is to be created",
	},
	"parity_group_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Parity group ID in which volume is to be created",
	},

	"subscriber_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Id of the subscriber which is attached to the volume",
	},

	"capacity": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "capacity(Size) of volume to be created",
	},
	//Remove dedup from this version

	"system": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "System name under volume to be created",
	},
	"deduplication_compression_mode": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "deduplicationCompressionMode of the volume to be created",
	},

	// output
	"volume": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "This is volume output",
		Elem: &schema.Resource{
			Schema: InfraVolumeInfoSchema,
		},
	},
}
