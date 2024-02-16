package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var InfraHostGroupInfoSchema = map[string]*schema.Schema{
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
	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of hostGroup",
	},
	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "ID of hostGroup",
	},
	"resource_group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "ID of Resource Group",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Assigned port of hostGroup",
	},
	"lun_paths": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "HostGroup lun paths with lun IDs and ldev IDs",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"lun_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Lun Path Lun ID",
				},
				"ldev_id": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Lun Path Ldev ID",
				},
			},
		},
	},
	"host_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Mode of hostgroup",
	},
	"wwns": &schema.Schema{
		Computed:    true,
		Description: "WWN list of hostGroup.",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "WWN ID",
				},
				"name": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "WWN Name",
				},
			},
		},
	},
	"host_mode_options": &schema.Schema{
		Computed:    true,
		Description: "The host mode options list",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host_mode_option": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Host mode option",
				},
				"host_mode_option_number": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Host mode option number",
				},
			},
		},
	},
}

var DataInfraHostGroupSchema = map[string]*schema.Schema{

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

	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port name",
	},

	"hostgroup_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Name of the hostGroup",
	},

	"hostgroup_number": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Unique ID of the hostGroup",
	},

	// output
	"hostgroup": &schema.Schema{
		Type:        schema.TypeList,
		Description: "This is hostgroups output",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: InfraHostGroupInfoSchema,
		},
	},
}

var DataInfraHostGroupsSchema = map[string]*schema.Schema{

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

	"port_ids": &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "A list of port IDs on storage",
	},

	"total_hostgroup_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of hostGroups in storage",
	},

	// output
	"hostgroups": &schema.Schema{
		Type:        schema.TypeList,
		Description: "This is hostgroups output",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: InfraHostGroupInfoSchema,
		},
	},
}

var ResourceInfraHostGroupSchema = map[string]*schema.Schema{

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
	"system": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The serial number of the preferred UCP system",
	},
	// output
	"hostgroup": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Response data of the created hostGroup",
		Elem: &schema.Resource{
			Schema: InfraHostGroupInfoSchema,
		},
	},
}
