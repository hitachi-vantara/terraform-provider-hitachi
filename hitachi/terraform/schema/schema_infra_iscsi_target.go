package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var InfraIscsiTargetInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
	},
	"resource_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource ID",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port ID",
	},
	"host_mode": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "Host Mode",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host_common_settings": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "host_common_settings",
				},
				"host_middle_ware": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "host_middle_ware",
				},
				"host_mode": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "host_mode",
				},
				"host_mode_options": &schema.Schema{
					Computed:    true,
					Optional:    true,
					Description: "host_mode_options",
					Type:        schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"df_option": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Syslog Server Address",
							},
							"is_ams_legal": {
								Computed:    true,
								Type:        schema.TypeBool,
								Description: "is_ams_legal",
							},
							"is_df": {
								Computed:    true,
								Type:        schema.TypeBool,
								Description: "is_df",
							},
							"is_hus_legal": {
								Computed:    true,
								Type:        schema.TypeBool,
								Description: "is_hus_legal",
							},
							"is_raid": {
								Computed:    true,
								Type:        schema.TypeBool,
								Description: "is_raid",
							},
							"raid_option": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "raid_option",
							},
							"raid_option_number": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "raid_option_number",
							},
						},
					},
				},
				"host_platform_option": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "host_platform_option",
				},
				"is_df": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "is_df",
				},
				"is_raid": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "is_raid",
				},
				"raid_host_mode_char": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "raid_host_mode_char",
				},
			},
		},
	},
	"resource_group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "ID of Resource Group",
	},
	"target_user": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Target User",
	},
	"iqn": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "IQN",
	},
	"iqn_initiators": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of iqn_initiators",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"logical_units": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "logical_units",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host_lun_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "host_lun_id",
				},
				"logical_unit_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "logical_unit_id",
				},
				"logical_unit_id_hex_format": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "logical_unit_id_hex_format",
				},
			},
		},
	},
	"auth_params": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "Storage Efficiency Statistics",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_chap_enabled": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "is_chap_enabled",
				},
				"is_chap_required": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "is_chap_required",
				},
				"is_mutual_auth": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "is_mutual_auth",
				},
				"authentication_mode": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "authentication_mode",
				},
			},
		},
	},
	"chap_users": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of chap_users",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"iscsi_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "iscsi_name",
	},
	"iscsi_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "ID of iscsi",
	},
}

var DataInfraIscsiTargetSchema = map[string]*schema.Schema{
	"storage_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Unique ID of the storage device",
	},

	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Serial Number of the storage device",
	},

	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port name",
	},

	"iscsi_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "ISCSI name",
	},

	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "ISCSI Id",
	},

	// output
	"iscsi_target": &schema.Schema{
		Type:        schema.TypeList,
		Description: "This is iscsi_target output",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: InfraIscsiTargetInfoSchema,
		},
	},
}

var DataInfraIscsiTargetsSchema = map[string]*schema.Schema{

	"storage_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Unique ID of the storage device",
	},

	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Serial Number of the storage device",
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
		Description: "Total number of iSCSI Target IDs which need to be fetched from the storage server",
	},

	// output
	"iscsi_targets": &schema.Schema{
		Type:        schema.TypeList,
		Description: "This is iscsi_targets output",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: InfraIscsiTargetInfoSchema,
		},
	},
}
