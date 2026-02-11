package terraform

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var HostGroupInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of the storage system",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Created hostGroup number",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Assigned port ID of hostGroup",
	},
	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of hostGroup",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Mode of hostgroup",
	},
	"host_mode_options": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
		Description: "The number of items in the host mode options list",
	},
	"ldevs": &schema.Schema{
		Computed: true,
		Type:     schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
		Description: "The number of ldev IDs in hostGroup",
	},
	"hg_luns": &schema.Schema{
		Computed: true,
		Type:     schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
		Description: "HostGroup lun IDs",
	},
	"lun_paths": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "HostGroup lun paths with lun IDs and ldev IDs",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hg_lun_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Lun Path hostGroup ID",
				},
				"ldev_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Lun Path Ldev ID",
				},
				"ldev_id_hex": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Lun Path Ldev ID in hexadecimal",
				},
			},
		},
	},
	"wwns": &schema.Schema{
		Computed: true,
		Type:     schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "WWN list of hostGroup.",
	},
	"wwns_detail": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Details of wwns for the created hostGroup including ID and name",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"wwn": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "ID of the WWN resource",
				},
				"name": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Name of the WWN resource",
				},
			},
		},
	},
}

var HostGroupsInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of the storage system",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "HostGroup number",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port ID on the storage system",
	},
	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "HostGroup name",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Type of host mode",
	},
}

var DataHostGroupSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port ID on the storage system",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "HostGroup number",
	},
	// output
	"hostgroup": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "hostGroup output",
		Elem: &schema.Resource{
			Schema: HostGroupInfoSchema,
		},
	},
}

var DataHostGroupsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"port_ids": &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "A list of port IDs on the storage system",
	},
	// output
	"hostgroups": &schema.Schema{
		Type:        schema.TypeList,
		Description: "hostgroups output",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: HostGroupsInfoSchema,
		},
	},
	"total_hostgroup_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of hostGroups in the storage system",
	},
}

var ResourceHostGroupSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port ID to be specified",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "HostGroup number to be specified",
	},
	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "HostGroup name to be specified to create the group",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The Host Mode to create the group",
	},
	"host_mode_options": &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
		Description: "The number of host mode options to create the group",
	},
	"lun": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Properties of LUNs to create the group",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"lun": {
					Optional:    true,
					Type:        schema.TypeInt,
					Description: "LUN ID",
				},
				"ldev_id": {
					Optional:     true,
					Type:         schema.TypeInt,
					Default:      -1,
					Description:  "LDEV ID. Only one of ldev_id or ldev_id_hex may be specified.",
					ValidateFunc: validation.IntBetween(0, 65535),
				},
				"ldev_id_hex": {
					Optional: true,
					Type:     schema.TypeString,
					// Default:     "-0x01",
					Description: "LDEV ID in hexadecimal format. Only one of ldev_id or ldev_id_hex may be specified.",
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
						"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
					),
				},
			},
		},
	},
	"wwn": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "World wide name of hostGroup",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host_wwn": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "Name of wwn resource",
				},
				"wwn_nickname": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "Nickname of wwn resource",
				},
			},
		},
	},
	// output
	"hostgroup": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Response data of the created hostGroup",
		Elem: &schema.Resource{
			Schema: HostGroupInfoSchema,
		},
	},
}
