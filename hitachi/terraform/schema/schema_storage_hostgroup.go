package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Description: "Created host group number",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Assigned port ID of the host group",
	},
	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the host group",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Mode of the host group",
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
		Description: "The number of LDEV IDs in the host group",
	},
	"hg_luns": &schema.Schema{
		Computed: true,
		Type:     schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
		Description: "Host group LUN IDs",
	},
	"lun_paths": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "Host group LUN paths with LUN IDs and LDEV IDs",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hg_lun_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "LUN path of the host group ID",
				},
				"ldev_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "LUN path of the LDEV ID",
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
		Description: "WWN list of the host group.",
	},
	"wwns_detail": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Details of WWNs for the created host group including ID and name",
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
		Description: "Host group number",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port ID on the storage system",
	},
	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Host group name",
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
		Description: "Host group number",
	},
	// output
	"hostgroup": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Host group output",
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
		Description: "Host groups output",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: HostGroupsInfoSchema,
		},
	},
	"total_hostgroup_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of host groups in the storage system",
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
		Description: "The Port ID to be specified",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Host group number to be specified",
	},
	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Host group name to be specified to create the group",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The host mode to create the group",
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
				"ldev_id": {
					Optional:    true,
					Type:        schema.TypeInt,
					Default:     -1,
					Description: "LDEV ID",
				},
				"lun": {
					Optional:    true,
					Type:        schema.TypeInt,
					Description: "LUN ID",
				},
			},
		},
	},
	"wwn": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "World wide name of the host group",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host_wwn": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "Name of the WWN resource",
				},
				"wwn_nickname": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "Nickname of the WWN resource",
				},
			},
		},
	},
	// output
	"hostgroup": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Response data of the created host group",
		Elem: &schema.Resource{
			Schema: HostGroupInfoSchema,
		},
	},
}
