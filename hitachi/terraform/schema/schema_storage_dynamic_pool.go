package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var DynamicPoolInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of the storage system",
	},
	"pool_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Pool ID of the storage system",
	},
	"pool_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool status of the storage system",
	},
	"used_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity rate",
	},
	"used_physical_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used physical capacity rate",
	},
	"snapshot_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Snapshot count",
	},
	"pool_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool name",
	},
	"available_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available volume capacity",
	},
	"available_physical_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available physical volume capacity",
	},
	"total_pool_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total pool capacity",
	},
	"total_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total physical capacity",
	},
	"num_of_ldevs": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of ldevs",
	},
	"first_ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "First ldev ID",
	},
	"first_ldev_id_hex": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "First ldev ID in hexadecimal",
	},
	"warning_threshold": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Warning threshold",
	},
	"depletion_threshold": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Depletion threshold",
	},
	"virtual_volume_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Virtual volume capacity rate",
	},
	"is_mainframe": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Is mainframe pool",
	},
	"is_shrinking": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Is shrinking pool",
	},
	"located_volume_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of located volume count",
	},
	"total_located_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of located capacity",
	},
	"blocking_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Blocking mode of pool",
	},
	"total_reserved_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of reserved capacity",
	},
	"reserved_volume_count": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of reserved volume count",
	},
	"pool_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Pool type",
	},
	"duplication_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Duplication number",
	},
	"effective_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Effective capacity",
	},
	"data_reduction_accelerate_comp_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Data reduction accelerate comp capacity",
	},
	"data_reduction_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Data reduction capacity",
	},
	"data_reduction_before_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Data reduction before capacity",
	},
	"data_reduction_accelerate_comp_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Data reduction accelerate comp rate",
	},
	"compression_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Compression rate",
	},
	"duplication_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Duplication rate",
	},
	"data_reduction_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Data reduction rate",
	},
	"snapshot_used_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Snapshot used capacity",
	},
	"suspend_snapshot": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Checks if suspend snapshot",
	},
	"data_reduction_accelerate_comp_including_system_data": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Data reduction accelerate comp including system data",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_reduction_capacity_available": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Is reduction capacity available",
				},
				"reduction_capacity": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Reduction capacity",
				},
				"is_reduction_rate_available": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Is reduction rate available",
				},
				"reduction_rate": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Reduction rate",
				},
			},
		},
	},
	"data_reduction_including_system_data": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Data reduction including system data",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_reduction_capacity_available": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Is reduction capacity available",
				},
				"reduction_capacity": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Reduction capacity",
				},
				"is_reduction_rate_available": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Is reduction rate available",
				},
				"reduction_rate": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Reduction rate",
				},
			},
		},
	},
	"capacities_excluding_system_data": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Capacities excluding system data",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"used_virtual_volume_capacity": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Used virtual volume capacity",
				},
				"compressed_capacity": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Compressed capacity",
				},
				"deduped_capacity": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Deduped capacity",
				},
				"reclaimed_capacity": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Reclaimed capacity",
				},
				"system_data_capacity": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "System data capacity",
				},
				"pre_used_capacity": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Pre used capacity",
				},
				"pre_compressed_capacity": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Pre compressed capacity",
				},
				"pre_dedupred_capacity": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Pre dedupred capacity",
				},
			},
		},
	},
	"efficiency": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Efficiency information (returned when detailInfoType includes efficiency)",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_calculated": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Whether efficiency is calculated",
				},
				"total_ratio": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Total efficiency ratio",
				},
				"compression_ratio": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Compression ratio",
				},
				"snapshot_ratio": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Snapshot ratio",
				},
				"provisioning_rate": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Provisioning rate",
				},
				"calculation_start_time": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Calculation start time",
				},
				"calculation_end_time": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Calculation end time",
				},
				"dedupe_and_compression": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Deduplication and compression efficiency",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"total_ratio": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Total ratio for dedup and compression",
							},
							"compression_ratio": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Compression ratio for dedup and compression",
							},
							"dedupe_ratio": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Deduplication ratio",
							},
							"reclaim_ratio": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Reclaim ratio",
							},
						},
					},
				},
				"accelerated_compression": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Accelerated compression efficiency",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"total_ratio": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Total ratio for accelerated compression",
							},
							"compression_ratio": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Compression ratio for accelerated compression",
							},
							"reclaim_ratio": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Reclaim ratio for accelerated compression",
							},
						},
					},
				},
			},
		},
	},
	"formatted_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Formatted capacity (returned when detailInfoType includes formattedCapacity)",
	},
	"auto_add_pool_vol": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Auto add pool volume setting (returned when detailInfoType includes autoAddPoolVol)",
	},
}

var DataDynamicPoolSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"pool_id": &schema.Schema{
		Type:          schema.TypeInt,
		Optional:      true,
		Computed:      true,
		ConflictsWith: []string{"pool_name"},
		Description:   "Pool ID of the storage system. Either `pool_id` or `pool_name` must be specified.",
	},
	"pool_name": &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Computed:      true,
		ConflictsWith: []string{"pool_id"},
		Description:   "Pool name of the storage system. Either `pool_name` or `pool_id` must be specified.",
	},
	// output
	"dynamic_pools": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Information about the dynamic pool.",
		Elem: &schema.Resource{
			Schema: DynamicPoolInfoSchema,
		},
	},
}

var DataDynamicPoolsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"pool_type": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filter pools by pool type. Supported values: `DP`, `HTI`. When not specified, returns all pools.",
		ValidateFunc: validation.StringInSlice([]string{
			"DP",
			"HTI",
		}, false),
	},
	"include_detail_info": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Include detailed information for dynamic pools. When set to true, additional detailed fields will be populated.",
	},
	"include_cache_info": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Include cache information for dynamic pools. When set to true, cache-related fields will be populated.",
	},
	"is_mainframe": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Filter pools by mainframe type. When set to true, returns only mainframe pools. When set to false, returns only non-mainframe pools. When not specified, returns all pools.",
	},
	// output
	"dynamic_pools": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "List of all dynamic pools retrieved from the storage system.",
		Elem: &schema.Resource{
			Schema: DynamicPoolInfoSchema,
		},
	},
}
