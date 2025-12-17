package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var StoragePoolSchema = map[string]*schema.Schema{
	"pool_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "ID of the pool",
	},
	"pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the pool",
	},
	"protection_domain_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Protection domain ID of the pool",
	},
	"status_summary": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status summary of the pool",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status of the pool",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity of the pool",
	},
	"total_raw_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total raw capacity of the pool",
	},
	"used_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity of the pool",
	},
	"free_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free capacity of the pool",
	},
	"total_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total physical capacity of the pool",
	},
	"meta_data_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Meta data physical capacity of the pool",
	},
	"reserved_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Reserved physical capacity of the pool",
	},
	"usable_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Usable physical capacity of the pool",
	},
	"blocked_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Blocked physical capacity of the pool",
	},
	"capacity_manage": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Capacity manage information",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"used_capacity_rate": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Used capacity rate of the pool",
				},
				"maximum_reserve_rate": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Maximum reserve rate of the pool",
				},
				"threshold_warning": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Threshold warning of the pool",
				},
				"threshold_depletion": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Threshold depletion of the pool",
				},
				"threshold_storage_controller_depletion": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Threshold storage controller depletion of the pool",
				},
			},
		},
	},
	"saving_effects": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Saving effects information",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"efficiency_data_reduction": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Efficiency data reduction of saving effects",
				},
				"pre_capacity_data_reduction": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Pre capacity data reduction of saving effects",
				},
				"post_capacity_data_reduction": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Post capacity data reduction of saving effects",
				},
				"total_efficiency_status": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Total efficiency status of saving effects",
				},
				"data_reduction_without_system_data_status": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Data reduction without system data status of saving effects",
				},
				"total_efficiency": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Total efficiency of saving effects",
				},
				"data_reduction_without_system_data": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Data reduction without system data of saving effects",
				},
				"pre_capacity_data_reduction_without_system_data": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Pre capacity data reduction without system data of saving effects",
				},
				"post_capacity_data_reduction_without_system_data": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Post capacity data reduction without system data of saving effects",
				},
				"calculation_start_time": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Calculation start time of saving effects",
				},
				"calculation_end_time": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Calculation end time of saving effects",
				},
			},
		},
	},
	"number_of_volumes": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of volumes on the pool",
	},
	"redundant_policy": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Number of volumes on the pool",
	},
	"redundant_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Redundant type of the pool",
	},
	"data_redundancy": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Data redundancy of the pool",
	},
	"storage_controller_capacities_general_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage controller capacities general status of the pool",
	},
	"total_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total volume capacity of the pool",
	},
	"provisioned_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Provisioned volume capacity of the pool",
	},
	"other_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Other volume capacity of the pool",
	},
	"temporary_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Temporary volume capacity of the pool",
	},
	"rebuild_capacity_policy": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Rebuild capacity policy of the pool",
	},
	"rebuild_capacity_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Rebuild capacity status of the pool",
	},
	"rebuild_capacity_resource_setting": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Rebuild capacity resource setting information",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"number_of_tolerable_drive_failures": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Number of tolerable drive failures of the pool",
				},
			},
		},
	},
	"rebuildable_resources": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Rebuildable resources information",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"number_of_drives": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "Number of drives of the pool",
				},
			},
		},
	},
}

var DatasourceVssbStoragePoolsSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block.",
	},
	"storage_pool_names": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "List of pool names to be retrieved from the storage system",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// output
	"storage_pools": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Storage pools output",
		Elem: &schema.Resource{
			Schema: StoragePoolSchema,
		},
	},
}

var ResourceVssbStoragePoolSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of VSP One SDS Block.",
	},
	"storage_pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Storage pool name",
	},
	"add_all_offline_drives": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false, // Defaults to false if not provided
		Description: "Flag to indicate if all offline drives should be added for expansion",
	},
	"drive_ids": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		// Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of specific offline drive IDs for expansion of the storage pool",
	},
	// Output:
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The status of the storage pool operation.",
	},
}
