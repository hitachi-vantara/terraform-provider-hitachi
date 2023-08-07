package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var IscsiTargetInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a storage serial number",
	},
	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a iscsi target number of the resource",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a port id of the resource",
	},
	"iscsi_target_alias": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a iscsi target alias of the resource",
	},
	"iscsi_target_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a iscsi target name of the resource",
	},
	"iscsi_target_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a iscsi target id of the resource",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's host mode of the resource",
	},
	"host_mode_options": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "It's a list of host mode options of the resource",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"ldevs": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Description: "It's a list of ldev ids of the resource",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"luns": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Description: "It's a list of luns of the resource",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"lun_paths": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "It's a list of lun_paths of the resource",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ldev_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "It's a ldev id of the resource",
				},
				"lun_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "It's a lun id of the resource",
				},
			},
		},
	},
	"initiators": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "It's a list of initiators of the resource",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"initiator_name": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "It's a initiator name of the resource",
				},
				"initiator_nickname": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "It's a initiator nickname of the resource",
				},
			},
		},
	},
}

var IscsiTargetsInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Storage serial number of the ISCSI target",
	},
	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "ISCSI target number of the ISCSI target",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port ID of the ISCSI target",
	},
	"iscsi_target_alias": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ISCSI target alias of the ISCSI target",
	},
	"iscsi_target_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ISCSI target ID of the ISCSI target",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Host mode of the ISCSI target",
	},
}

var DataIscsiTargetSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage device is required",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "It's a port id",
	},
	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "It's a iscsi target number",
	},
	// output
	"iscsitarget": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: IscsiTargetInfoSchema,
		},
	},
}

var DataIscsiTargetsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage server",
	},
	"port_ids": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "List of Port IDs which need to be fetched from the storage server",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"total_iscsi_target_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of ISCSI Target IDs which need to be fetched from the storage server",
	},
	// output
	"iscsitargets": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "More information about the ISCSI Target",
		Elem: &schema.Resource{
			Schema: IscsiTargetsInfoSchema,
		},
	},
}

var ResourceIscsiTargetSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage device is required",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port id in which the resource to be created",
	},
	"iscsi_target_number": &schema.Schema{ // similar to hostgroup number
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Resource will be created based on iscsi target number",
	},
	"iscsi_target_alias": &schema.Schema{ // similar to hostgroup name
		Type:        schema.TypeString,
		Required:    true,
		Description: "It's a iscsi target alias",
	},
	"iscsi_target_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "It's a iscsi target name",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Host mode value to be given to create the resource",
	},
	"host_mode_options": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Host mode options can be passed to create the resource",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"lun": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Lun input for the resource",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ldev_id": {
					Optional:    true,
					Type:        schema.TypeInt,
					Description: "Ldev id to create the resource",
				},
				"lun_id": {
					Optional:    true,
					Type:        schema.TypeInt,
					Description: "Lun id to create the resource",
				},
			},
		},
	},
	"initiator": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Initiator input for the resource",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"initiator_name": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "Initiator name to create the resource",
				},
				"initiator_nickname": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "Initiator nickname to create the resource",
				},
			},
		},
	},
	// output
	"iscsitarget": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: IscsiTargetInfoSchema,
		},
	},
}
