package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var InfraStorageDeviceSchema = map[string]*schema.Schema{
	"resource_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource  ID",
	},
	"serial_number": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Serial Number",
	},
	"management_address": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Management Address",
	},

	"controller_address": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Controller Address",
	},

	"username": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Username",
	},

	"systems": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of UCP Systems",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},

	"device_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Device Type",
	},
	"model": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Model",
	},
	"microcode_version": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Microcode Version",
	},
	"total_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total Capacity",
	},
	"free_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free Capacity",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Total Capacity",
	},
	"free_capacity": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Free Capacity",
	},

	"resource_state": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource State",
	},

	"tags": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of tags",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},

	"health_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Health Status",
	},

	"operational_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Operational Status",
	},

	"free_gad_consistency_group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free GAD Consistency Group Id",
	},
	"free_local_clone_consistency_group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free Local Clone Consistency Group Id",
	},
	"free_remote_clone_consistency_group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free Remote Clone Consistency Group Id",
	},

	"storage_efficiency_stat": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "Storage Efficiency Statistics",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"compression_ratio": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Compression Ratio",
				},
				"start_time": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Start Time",
				},
				"end_time": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "End Time",
				},
				"provisioning_rate": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Provisioning Rate",
				},
				"snapshot_ratio": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Snapshot Ratio",
				},
				"total_ratio": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Snapshot Ratio",
				},
				"accel_comp_efficiency_stat": &schema.Schema{
					Computed:    true,
					Optional:    true,
					Description: "accel_comp_efficiency_stat",
					Type:        schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"compression_ratio": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Compression Ratio",
							},
							"dedupe_ratio": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Dedupe Ratio",
							},
							"reclaim_ratio": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Reclaim Ratio",
							},
							"total_ratio": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Total Ratio",
							},
						},
					},
				},
				"dedupe_comp_efficiency_stat": &schema.Schema{
					Computed:    true,
					Optional:    true,
					Description: "dedupe_comp_efficiency_stat",
					Type:        schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"compression_ratio": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Compression Ratio",
							},
							"dedupe_ratio": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Dedupe Ratio",
							},
							"reclaim_ratio": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Reclaim Ratio",
							},
							"total_ratio": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Total Ratio",
							},
						},
					},
				},
			},
		},
	},
	"syslog_config": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "Storage Efficiency Statistics",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"detailed": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Compression Ratio",
				},
				"syslog_servers": &schema.Schema{
					Computed:    true,
					Optional:    true,
					Description: "Syslog Servers",
					Type:        schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Id",
							},
							"syslog_server_address": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Syslog Server Address",
							},
							"syslog_server_port": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Syslog Server Port",
							},
						},
					},
				},
			},
		},
	},
	"storage_device_licenses": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "storage_device_licenses",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_enabled": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "is_enabled",
				},
				"is_installed": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "is_installed",
				},
				"type": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "License Type",
				},
				"name": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "License Name",
				},
			},
		},
	},
	"device_limits": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "device_limits",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"health_description": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Health Description",
				},
				"is_unified": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "is_unified",
				},
				"external_group_number_range": &schema.Schema{
					Computed:    true,
					Optional:    true,
					Description: "external_group_number_range",
					Type:        schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"is_valid": {
								Computed:    true,
								Type:        schema.TypeBool,
								Description: "is_valid",
							},
							"max_value": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "max_value",
							},
							"min_value": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Reclaim Ratio",
							},
						},
					},
				},
				"external_group_sub_number_range": &schema.Schema{
					Computed:    true,
					Optional:    true,
					Description: "external_group_sub_number_range",
					Type:        schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"is_valid": {
								Computed:    true,
								Type:        schema.TypeBool,
								Description: "is_valid",
							},
							"max_value": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "max_value",
							},
							"min_value": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Reclaim Ratio",
							},
						},
					},
				},
				"parity_group_number_range": &schema.Schema{
					Computed:    true,
					Optional:    true,
					Description: "parity_group_number_range",
					Type:        schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"is_valid": {
								Computed:    true,
								Type:        schema.TypeBool,
								Description: "is_valid",
							},
							"max_value": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "max_value",
							},
							"min_value": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Reclaim Ratio",
							},
						},
					},
				},
				"parity_group_sub_number_range": &schema.Schema{
					Computed:    true,
					Optional:    true,
					Description: "parity_group_sub_number_range",
					Type:        schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"is_valid": {
								Computed:    true,
								Type:        schema.TypeBool,
								Description: "is_valid",
							},
							"max_value": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "max_value",
							},
							"min_value": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Reclaim Ratio",
							},
						},
					},
				},
			},
		},
	},
}

var DataInfraStorageDevicesSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Serial Number of the storage device",
	},
	// output
	"storage_devices": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is storage devices output",
		Elem: &schema.Resource{
			Schema: InfraStorageDeviceSchema,
		},
	},
}

var ResourceInfraStorageDeviceSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial Number of the storage device",
	},

	"management_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Management Address",
	},

	"username": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Management Address",
	},

	"password": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Management Address",
	},

	"gateway_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Gateway Address",
	},

	"out_of_band": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Out of band",
	},

	"system": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "ID of the UCP System",
	},

	// output
	"storage_devices": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is storage devices output",
		Elem: &schema.Resource{
			Schema: InfraStorageDeviceSchema,
		},
	},
}
