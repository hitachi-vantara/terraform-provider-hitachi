package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var HostGroupInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The serial number of the storage",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Created HostGroup number",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Assigned port ID of the HostGroup",
	},
	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the HostGroup",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Mode of the hostgroup",
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
		Description: "The number of ldevs Ids in the host group",
	},
	"hg_luns": &schema.Schema{
		Computed: true,
		Type:     schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
		Description: "Host group lun IDs",
	},
	"lun_paths": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "HostGroup lun paths with lun ids and ldev ids",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hg_lun_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Lun Paths Host group ID",
				},
				"ldev_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Lun Paths Ldev ID",
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
		Description: "WWNs list for the created host group.",
	},
	"wwns_detail": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Details of wwns for the created host group including id and name",
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
		Description: "It's a storage serial number",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "It's a hostgroup number",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a port id",
	},
	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's a name of hostgroup",
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "It's host mode",
	},
}

var DataHostGroupSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the HostGroup",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "PortId of the HostGroup",
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
		Description: "Additional information of HostGroup",
		Elem: &schema.Resource{
			Schema: HostGroupInfoSchema,
		},
	},
}

var DataHostGroupsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage device is required",
	},
	"port_ids": &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "A list of port ids of hostgroup",
	},
	"total_hostgroup_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of HostGroups in storage",
	},
	// output
	"hostgroups": &schema.Schema{
		Type:        schema.TypeList,
		Description: "This is output schema",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: HostGroupsInfoSchema,
		},
	},
}

var ResourceHostGroupSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The host group serial number need to be specified",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Port ID to be specified",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The Host Group number to be specified",
	},
	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Host Group name to be specified to create the group",
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
		Description: "The number of host mode options to be given to create the group.",
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
					Description: "Ldev ID of the Lun",
				},
				"lun": {
					Optional:    true,
					Type:        schema.TypeInt,
					Description: "Lun ID of the Lun",
				},
			},
		},
	},
	"wwn": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "World wide name of the Host group",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host_wwn": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "Name of the wwn resource",
				},
				"wwn_nickname": {
					Optional:    true,
					Type:        schema.TypeString,
					Description: "Nickname of the wwn resource",
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
