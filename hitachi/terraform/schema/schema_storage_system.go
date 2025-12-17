package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var SSInfoSchema = map[string]*schema.Schema{
	"storage_device_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage system ID",
	},
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of the storage system",
	},
	"storage_device_model": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage system model",
	},
	"dkc_micro_code_version": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "DKC micro code version of the storage system",
	},
	"management_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage management IP address",
	},
	"svp_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage SVP IP address",
	},
	"controller1_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage controller1 IP address",
	},
	"controller2_ip": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage controller2 IP address",
	},
	"free_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free capacity in MB",
	},
	"used_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity in MB",
	},
	"total_capacity_in_mb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity in MB",
	},
}

// first 4 are inputs
var StorageSystemSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	// output
	"storage_system": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Storage system output",
		Elem: &schema.Resource{
			Schema: SSInfoSchema,
		},
	},
}

// Storage System VSP One Schema
var SSInfoAdminSchema = map[string]*schema.Schema{
	"model_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Model name",
	},
	"serial": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Serial number",
	},
	"nickname": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Nickname",
	},
	"number_of_total_volumes": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of created volumes",
	},
	"number_of_free_drives": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of available drives",
	},
	"number_of_total_servers": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of registered servers",
	},
	"total_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity (sum) of parity groups (MiB)",
	},
	"total_pool_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "POOL effective capacity (MiB)",
	},
	"total_pool_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "POOL total capacity (MiB)",
	},
	"used_pool_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "POOL total used capacity (MiB)",
	},
	"free_pool_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "POOL total free capacity (POOL effective capacity - POOL total used capacity) (MiB)",
	},
	"total_pool_capacity_with_ti_pool": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "POOL effective capacity (MiB) with TI pool. If unsupported, returns -1.",
	},
	"total_pool_physical_capacity_with_ti_pool": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "POOL total capacity (MiB) with TI pool. If unsupported, returns -1.",
	},
	"used_pool_capacity_with_ti_pool": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "POOL total used capacity (MiB) with TI pool. If unsupported, returns -1.",
	},
	"free_pool_capacity_with_ti_pool": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "POOL effective capacity - POOL total used capacity (MiB) with TI pool. If unsupported, returns -1.",
	},
	"estimated_configurable_pool_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Estimated configurable pool capacity (MiB) for Dynamic Provisioning pools. Only if withEstimatedConfigurableCapacities is true. If not valid, returns -1.",
	},
	"estimated_configurable_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Estimated configurable volume capacity (MiB) for volumes from Dynamic Provisioning pools. Only if withEstimatedConfigurableCapacities is true. If not valid, returns -1.",
	},
	"saving_effects": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Saving effects and efficiency data (see nested fields)",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"efficiency_data_reduction": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "System-wide controller-based capacity reduction effect ratio. Returns (n x 100) for n:1. If no pool exists or no user data is in the pool, returns -1.",
				},
				"pre_capacity_data_reduction": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "System-wide controller-based capacity before reduction (MiB)",
				},
				"post_capacity_data_reduction": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "System-wide controller-based capacity after reduction (MiB)",
				},
				"efficiency_fmd_saving": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Always returns -1.",
				},
				"pre_capacity_fmd_saving": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Always returns 0.",
				},
				"post_capacity_fmd_saving": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Always returns 0.",
				},
				"is_total_efficiency_support": &schema.Schema{
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "TotalEfficiency support flag",
				},
				"total_efficiency_status": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Status: NotSupported, Valid, CalculationInProgress, NoTargetData, Unknown",
				},
				"data_reduction_without_system_data_status": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Status: NotSupported, Valid, CalculationInProgress, NoTargetData, Unknown",
				},
				"software_saving_without_system_data_status": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Status: NotSupported, Valid, CalculationInProgress, NoTargetData, Unknown",
				},
				"total_efficiency": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Sum ratio of data reduction effect (excluding system data), snapshot effect, and provisioning effect (value multiplied by 100). If unsupported or status not valid, returns -1.",
				},
				"data_reduction_without_system_data": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Total reduction effect ratio by capacity reduction function (value multiplied by 100). If unsupported or status not valid, returns -1.",
				},
				"pre_capacity_data_reduction_without_system_data": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "System-wide controller-based capacity before reduction (MiB). If unsupported, returns -1.",
				},
				"post_capacity_data_reduction_without_system_data": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "System-wide controller-based capacity after reduction (MiB). If unsupported, returns -1.",
				},
				"software_saving_without_system_data": &schema.Schema{
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Data reduction ratio by capacity reduction function (value multiplied by 100). If unsupported or status not valid, returns -1.",
				},
				"calculation_start_time": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Calculation start date and time of. If unsupported or status not valid, returns null.",
				},
				"calculation_end_time": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Calculation end date and time. If unsupported or status not valid, returns null.",
				},
			},
		},
	},
	"gum_version": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "GUM Version. Example: 'A3-04-01/00'",
	},
	"esm_os_version": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ESM OS Version. Example: 'A3-04-01/00'",
	},
	"dkc_micro_version": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "DKC Main Version. Example: 'A3-04-01-40/00'",
	},
	"warning_led_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "LED Status. OFF: normal, ON: blocked/faulty, BLINK: unreferenced SIM, Unknown: unknown",
	},
	"esm_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ESM status",
	},
	"ip_address_ipv4_service": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "(nullable) IPv4 service address",
	},
	"ip_address_ipv4_ctl1": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "(nullable) IPv4 address (CTL01)",
	},
	"ip_address_ipv4_ctl2": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "(nullable) IPv4 address (CTL02)",
	},
	"ip_address_ipv6_service": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "(nullable) IPv6 service address",
	},
	"ip_address_ipv6_ctl1": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "(nullable) IPv6 address (CTL01)",
	},
	"ip_address_ipv6_ctl2": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "(nullable) IPv6 address (CTL02)",
	},
}

// first 4 are inputs
var StorageSystemAdminSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"with_estimated_configurable_capacities": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Flag to select whether to display configurable estimation information.",
		Default:     false,
	},
	// output
	"storage_system_admin": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "VSP One Storage system output",
		Elem: &schema.Resource{
			Schema: SSInfoAdminSchema,
		},
	},
}
