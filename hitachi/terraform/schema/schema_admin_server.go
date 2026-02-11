package terraform

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// AdminServerInfoSchema defines the shared schema for VSP One server information used by both resources and data sources
var AdminServerInfoSchema = map[string]*schema.Schema{
	"server_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Server ID assigned by the storage system.",
	},
	"nickname": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Server nickname.",
	},
	"protocol": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Protocol. FC or iSCSI. If is_reserved is true, Undefined.",
	},
	"os_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "OS type. If is_reserved is true, Undefined. Example values: Linux, VMware, HP-UX, OpenVMS, Tru64, Solaris, NetWare, Windows, AIX, Undefined, Unknown.",
	},
	"os_type_options": {
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeInt},
		Description: "OS type option list.",
	},
	"total_capacity": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Total assigned volume capacity (MiB).",
	},
	"used_capacity": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Volume usage (assigned capacity) (MiB).",
	},
	"number_of_volumes": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Number of assigned volumes.",
	},
	"number_of_paths": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Number of registered HBAs. Maximum 32*256.",
	},
	"paths": {
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "List of PATH information associated with HBA WWN.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hba_wwn": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Server HBA WWN (nullable).",
				},
				"iscsi_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Server iSCSI name (nullable).",
				},
				"port_ids": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "List of assigned ports.",
				},
			},
		},
	},
	"is_inconsistent": {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "Configuration information inconsistency flag.",
	},
	"is_reserved": {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "Indicates whether the server is for host group addition. If true, indicates a server reserved only with id and server nickname for adding host groups.",
	},
	"has_non_fullmesh_lu_paths": {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "True if any of attached volumes is not fully associated with ports which connect to the server.",
	},
	"has_unaligned_os_types": {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "OSType mismatch flag.",
	},
	"has_unaligned_os_type_options": {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "OSTypeOption mismatch flag.",
	},
}

// ResourceAdminServerSchema defines the schema for the VSP One server resource
var ResourceAdminServerSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage system.",
	},
	"server_nickname": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Server nickname.",
		ValidateFunc: validation.StringLenBetween(1, 32),
	},
	"protocol": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Protocol. FC or iSCSI. Required if is_reserved is false or omitted.",
		ValidateFunc: validation.StringInSlice([]string{
			"FC",
			"iSCSI",
		}, false),
	},
	"os_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "OS type. Required if is_reserved is false or omitted.",
		ValidateFunc: validation.StringInSlice([]string{
			"Linux",
			"HP-UX",
			"OpenVMS",
			"Solaris",
			"AIX",
			"VMware",
			"Windows",
		}, false),
	},
	"os_type_options": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "OS type option list. If osType is specified and this attribute is omitted, a value corresponding to the specified osType is automatically set.",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"is_reserved": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Indicates whether the server is for host group addition. If true, a server with only id and server nickname reserved for host group addition will be created.",
	},
	"keep_lun_config": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specify whether to delete the server information while retaining the volume assignment information during delete operations.",
	},

	// for add and sync hg
	"host_groups": {
		Type:     schema.TypeList,
		Optional: true,
		Description: `Specifies the list of existing host groups to associate with the server.
  Note:
  - There is no API to remove a host group from a server.
  - Any host group listed here will be added and remain associated even if removed from the list later.
  - It is not required to keep a host group in the list after adding it, but if it remains, it will be skipped on subsequent applies.
  - After host groups are added, their names are automatically synchronized with the server nickname.
  - All host groups must already exist.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"port_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Port ID associated with the host group.",
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^[A-Za-z0-9-]+$`),
						"port_id must contain only letters, numbers, and hyphens (e.g., CL1-A)",
					),
				},
				"host_group_id": {
					Type:         schema.TypeInt,
					Optional:     true,
					Description:  "Host Group ID. Either host_group_id or host_group_name must be specified, but not both.",
					ValidateFunc: validation.IntAtLeast(0),
				},

				"host_group_name": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Host Group Name. Either host_group_name or host_group_id must be specified, but not both.",
					ValidateFunc: validation.StringLenBetween(1, 64),
				},
			},
		},
	},
	// Output block
	"data": {
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Server information after creation/update.",
		Elem: &schema.Resource{
			Schema: AdminServerInfoSchema,
		},
	},
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Terraform resource ID in format 'serial-serverid'.",
	},
}

// --- Data source schema -----------------

// Schema for listing servers
var DataSourceAdminServerListSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The serial number of the storage.",
	},
	"nickname": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Server nickname (exact match). If specified, the list will be filtered by the specified conditions.",
	},
	"hba_wwn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "HBA WWN (exact match). If specified, the list will be filtered by the specified conditions.",
	},
	"iscsi_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "iSCSI name (exact match). If specified, the list will be filtered by the specified conditions.",
	},
	// Outputs
	"data": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of server information.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server_id": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "Server ID.",
				},
				"nickname": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Server nickname.",
				},
				"protocol": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Protocol. FC or iSCSI. If is_reserved is true, Undefined.",
				},
				"os_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "OS type. If is_reserved is true, Undefined. Example values: Linux, VMware, HP-UX, OpenVMS, Tru64, Solaris, NetWare, Windows, AIX, Undefined, Unknown.",
				},
				"total_capacity": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "Total assigned volume capacity (MiB).",
				},
				"used_capacity": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "Volume usage (assigned capacity) (MiB).",
				},
				"number_of_paths": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "Number of registered HBAs. Maximum 32*256.",
				},
				"is_inconsistent": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "Configuration information inconsistency flag.",
				},
				"modification_in_progress": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "Currently unused attribute.",
				},
				"is_reserved": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "Indicates whether the server is for host group addition. If true, indicates a server reserved only with id and server nickname for adding host groups.",
				},
				"has_unaligned_os_types": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "OSType mismatch flag.",
				},
			},
		},
	},
	"server_count": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Number of items stored in data.",
	},
}

// Schema for getting server info by ID
var DataSourceAdminServerInfoSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The serial number of the storage.",
	},
	"server_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Server ID.",
	},
	// Outputs
	"data": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Server information.",
		Elem: &schema.Resource{
			Schema: AdminServerInfoSchema,
		},
	},
}

// DataSourceAdminServerPathSchema defines the schema for getting server path info
var DataSourceAdminServerPathSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The serial number of the storage.",
	},
	"server_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Server ID.",
	},
	"hba_wwn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Server HBA WWN. Either hba_wwn or iscsi_name must be specified.",
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			return
		},
	},
	"iscsi_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Server iSCSI name. Either hba_wwn or iscsi_name must be specified.",
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			return
		},
	},
	"port_id": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port ID.",
	},
	// Outputs
	"data": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Server path information.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Path ID in format '<hbaWwn|iscsiName>,<portId>'.",
				},
				"server_id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Server ID.",
				},
				"hba_wwn": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Server HBA WWN (nullable).",
				},
				"iscsi_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Server iSCSI name (nullable).",
				},
				"port_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Assigned port ID.",
				},
			},
		},
	},
}

// ResourceAdminServerPathSchema defines the schema for the VSP One server path resource
var ResourceAdminServerPathSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The serial number of the storage.",
	},
	"server_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Server ID.",
	},
	"hba_wwn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Server HBA WWN. Either hba_wwn or iscsi_name must be specified, but not both.",
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			return
		},
	},
	"iscsi_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Server iSCSI name. Either hba_wwn or iscsi_name must be specified, but not both.",
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			return
		},
	},
	"port_ids": {
		Type:        schema.TypeList,
		Required:    true,
		Description: "List of port IDs to associate with the server path.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// Output attributes
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Path ID in format 'serial-server_id-hbaWwn|iscsiName'.",
	},
}
