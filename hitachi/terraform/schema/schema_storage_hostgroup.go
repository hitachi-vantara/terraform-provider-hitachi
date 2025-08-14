package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var HostGroupInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
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
		Description: "Serial number of storage",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "HostGroup number",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port ID on storage",
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
		Description: "Serial number of storage",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port ID on storage",
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
		Description: "This is hostGroup output",
		Elem: &schema.Resource{
			Schema: HostGroupInfoSchema,
		},
	},
}

var DataHostGroupsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage",
	},
	"port_ids": &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "A list of port IDs on storage",
	},
	// output
	"hostgroups": &schema.Schema{
		Type:        schema.TypeList,
		Description: "This is hostgroups output",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: HostGroupsInfoSchema,
		},
	},
	"total_hostgroup_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of hostGroups in storage",
	},
}

var ResourceHostGroupSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Port ID to be specified",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
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
		Description: "The number of host mode options to be given to create the group",
	},
	"lun": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Properties of Luns to create the group",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ldev_id": {
					Optional:    true,
					Type:        schema.TypeInt,
					Default:     -1,
					Description: "Ldev ID of Lun",
				},
				"lun": {
					Optional:    true,
					Type:        schema.TypeInt,
					Description: "Lun ID of Lun",
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
