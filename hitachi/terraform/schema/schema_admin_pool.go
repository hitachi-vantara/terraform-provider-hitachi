package terraform

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// VSP One pool schemas follow ports pattern: split input/output and compose.

func adminPoolInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pool_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The unique identifier of the storage pool.",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the storage pool.",
		},
		"status": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The current status of the storage pool.",
		},
		"encryption_status": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The encryption status of the storage pool.",
		},
		"total_capacity": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The total capacity of the storage pool in MiB.",
		},
		"effective_capacity": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The effective capacity of the storage pool in MiB.",
		},
		"used_capacity": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The used capacity of the storage pool in MiB.",
		},
		"free_capacity": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The free capacity of the storage pool in MiB.",
		},
		"capacity_manage": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Capacity management information for the storage pool.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"used_capacity_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "The used capacity rate as a percentage.",
					},
					"threshold_warning": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "The warning threshold for capacity usage as a percentage.",
					},
					"threshold_depletion": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "The depletion threshold for capacity usage as a percentage.",
					},
				},
			},
		},
		"saving_effects": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Saving effects information for the storage pool.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"efficiency_data_reduction": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Data reduction efficiency.",
					},
					"efficiency_fmd_saving": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "FMD saving efficiency.",
					},
					"pre_capacity_fmd_saving": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Pre-FMD saving capacity.",
					},
					"post_capacity_fmd_saving": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Post-FMD saving capacity.",
					},
					"is_total_efficiency_support": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Whether total efficiency is supported.",
					},
					"total_efficiency_status": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Total efficiency status.",
					},
					"data_reduction_without_system_data_status": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Data reduction status without system data.",
					},
					"software_saving_without_system_data_status": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Software saving status without system data.",
					},
					"total_efficiency": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Total efficiency value.",
					},
					"data_reduction_without_system_data": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Data reduction value without system data.",
					},
					"software_saving_without_system_data": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Software saving value without system data.",
					},
					"calculation_start_time": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Calculation start time.",
					},
					"calculation_end_time": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Calculation end time.",
					},
				},
			},
		},
		"config_status": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Configuration status list.",
		},
		"number_of_volumes": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of volumes in the pool.",
		},
		"number_of_tiers": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of tiers in the pool.",
		},
		"number_of_drive_types": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of drive types in the pool.",
		},
		"drives": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Drive information for the storage pool.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"drive_type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The type of the drive.",
					},
					"drive_interface": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The interface type of the drive.",
					},
					"drive_rpm": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The RPM of the drive.",
					},
					"drive_capacity": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "The capacity of individual drives in GB.",
					},
					"display_drive_capacity": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Capacity of the drive and unit of measurement (GB or TB)",
					},
					"total_capacity": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Total capacity for this drive type in MiB.",
					},
					"number_of_drives": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Number of drives of this type.",
					},
					"locations": {
						Type:        schema.TypeList,
						Optional:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Description: "Drive location identifiers.",
					},
					"raid_level": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "RAID level for the drives.",
					},
					"parity_group_type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Parity group type for the drives.",
					},
				},
			},
		},
		"contains_capacity_saving_volume": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the pool contains capacity saving volumes.",
		},
	}
}

// -------- Resource VSP One Pool --------
func resourceAdminPoolInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of the storage system.",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "The name of the storage pool.",
			ValidateFunc: validatePoolName,
		},
		"encryption": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether encryption is enabled for the storage pool.",
		},
		"drive_configuration": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "The list of drives for the storage pool. Adding new drive blocks will trigger pool expansion.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"drive_type_code": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The type code of the drive (12 characters, e.g., SNB5B-R1R9NC).",
					},
					"data_drive_count": {
						Type:         schema.TypeInt,
						Required:     true,
						Description:  "The number of data drives (1-1440). For RAID5: minimum 5, for RAID6: minimum 9.",
						ValidateFunc: validation.IntBetween(1, 1440),
					},
					"raid_level": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The RAID level for the drives. Supported values: RAID5, RAID6.",
						ValidateFunc: validation.StringInSlice([]string{"RAID5", "RAID6"}, false),
					},
					"parity_group_type": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "DDP",
						Description:  "The parity group type for the drives. Must be DDP.",
						ValidateFunc: validation.StringInSlice([]string{"DDP"}, false),
					},
				},
			},
		},
		"threshold_warning": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Warning threshold for capacity usage as a percentage.",
			ValidateFunc: validation.IntBetween(0, 100),
		},
		"threshold_depletion": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Depletion threshold for capacity usage as a percentage.",
			ValidateFunc: validation.IntBetween(0, 100),
		},
	}
}

func resourceAdminPoolOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Data attribute containing pool information
		"data": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Pool data information.",
			Elem: &schema.Resource{
				Schema: adminPoolInfoSchema(),
			},
		},
	}
}

func ResourceAdminPoolSchema() map[string]*schema.Schema {
	schema := resourceAdminPoolInputSchema()

	for k, v := range resourceAdminPoolOutputSchema() {
		schema[k] = v
	}

	return schema
}

// -------- Datasource: VSP One Pool (single) --------
func datasourceAdminPoolInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of the storage system.",
		},
		"pool_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The unique identifier of the storage pool.",
		},
	}
}

func datasourceAdminPoolOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Data attribute containing pool information
		"data": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Pool data information.",
			Elem: &schema.Resource{
				Schema: adminPoolInfoSchema(),
			},
		},
	}
}

func DataSourceAdminPoolSchema() map[string]*schema.Schema {
	schema := datasourceAdminPoolInputSchema()

	for k, v := range datasourceAdminPoolOutputSchema() {
		schema[k] = v
	}

	return schema
}

// -------- Datasource: VSP One Pools (list) --------
func datasourceAdminPoolsInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of the storage system.",
		},
		"name_filter": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Filter pools by name.",
		},
		"status_filter": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Filter pools by status.",
		},
		"config_status_filter": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Filter pools by configuration status.",
		},
	}
}

func datasourceAdminPoolsOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pool_counts": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Number of storage pools returned.",
		},
		"data": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "List of storage pools.",
			Elem: &schema.Resource{
				Schema: adminPoolInfoSchema(),
			},
		},
	}
}

func DataSourceAdminPoolsSchema() map[string]*schema.Schema {
	schema := datasourceAdminPoolsInputSchema()

	for k, v := range datasourceAdminPoolsOutputSchema() {
		schema[k] = v
	}

	return schema
}

// validatePoolName validates pool name according to API requirements
func validatePoolName(val interface{}, key string) (warns []string, errs []error) {
	name := val.(string)

	if len(name) < 1 || len(name) > 32 {
		errs = append(errs, fmt.Errorf("%q must be between 1 and 32 characters, got %d", key, len(name)))
		return
	}

	if strings.HasPrefix(name, " ") || strings.HasSuffix(name, " ") {
		errs = append(errs, fmt.Errorf("%q cannot start or end with a space character", key))
		return
	}

	if strings.HasPrefix(name, "-") {
		errs = append(errs, fmt.Errorf("%q cannot start with a hyphen (-)", key))
		return
	}

	validPattern := regexp.MustCompile(`^[0-9A-Za-z ,./:@\\_-]+$`)
	if !validPattern.MatchString(name) {
		errs = append(errs, fmt.Errorf("%q can only contain alphanumeric characters, spaces, and the following symbols: , - . / : @ \\ _", key))
		return
	}
	return
}
