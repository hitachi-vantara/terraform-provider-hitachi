package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var IscsiTargetInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
	},
	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "iSCSI target number",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port ID on storage",
	},
	"iscsi_target_alias": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "iSCSI target alias",
	},
	"iscsi_target_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "iSCSI target name",
	},
	"iscsi_target_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "iSCSI target ID",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Type of host mode",
	},
	"host_mode_options": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of host mode options",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"ldevs": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Description: "List of ldev IDs",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"luns": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Description: "List of luns",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"lun_paths": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "List of lun_paths of the resource",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ldev_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Ldev ID of lun",
				},
				"lun_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Lun ID of lun",
				},
			},
		},
	},
	"initiators": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "List of initiators",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"initiator_name": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Initiator name",
				},
				"initiator_nickname": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Initiator nickname",
				},
			},
		},
	},
}

var IscsiTargetsInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
	},
	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "iSCSI target number",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port ID of the iSCSI target",
	},
	"iscsi_target_alias": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "iSCSI target alias of the iSCSI target",
	},
	"iscsi_target_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "iSCSI target name",
	},
	"iscsi_target_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "iSCSI target ID of the iSCSI target",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Type of host mode",
	},
	"host_mode_options": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of host mode options",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
}

var DataIscsiTargetSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port ID on storage",
	},
	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "iSCSI target number",
	},
	// output
	"iscsitarget": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is iSCSI target output",
		Elem: &schema.Resource{
			Schema: IscsiTargetInfoSchema,
		},
	},
}

var DataIscsiTargetsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage",
	},
	"port_ids": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "List of Port IDs which need to be fetched from the storage server",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"total_iscsi_target_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of iSCSI Target IDs which need to be fetched from the storage server",
	},
	// output
	"iscsitargets": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is iSCSI Target output",
		Elem: &schema.Resource{
			Schema: IscsiTargetsInfoSchema,
		},
	},
}

var ResourceIscsiTargetSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage is required",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port ID in which the resource to be created",
	},
	"iscsi_target_number": &schema.Schema{ // similar to hostgroup number
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Resource will be created based on iSCSI target number",
	},
	"iscsi_target_alias": &schema.Schema{ // similar to hostgroup name
		Type:        schema.TypeString,
		Required:    true,
		Description: "iSCSI target alias",
	},
	"iscsi_target_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "iSCSI target name",
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
					Description: "Ldev ID of lun",
				},
				"lun_id": {
					Optional:    true,
					Type:        schema.TypeInt,
					Description: "Lun ID of lun",
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
		Description: "This is iSCSI target output",
		Elem: &schema.Resource{
			Schema: IscsiTargetInfoSchema,
		},
	},
}
