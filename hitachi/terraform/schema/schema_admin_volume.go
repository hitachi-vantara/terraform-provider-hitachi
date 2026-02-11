package terraform

import (
	"fmt"
	"regexp"

	utils "terraform-provider-hitachi/hitachi/common/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// separate input and output schema for better readability, then combine them

// ------------------- Volumes Info Schema -------------------
func volumesInfoListSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"volumes_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of volumes.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"volume_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Volume ID.",
					},
					"volume_id_hex": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Volume ID in hexadecimal format.",
					},
					"nickname": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "(Nullable) Nickname.",
					},
					"pool_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Pool ID.",
					},
					"pool_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "(Nullable) Pool name. Not returned if the volume is in a transitional state.",
					},
					"total_capacity_mb": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Total capacity in MiB. Max 256TiB = 268,435,456MiB = 549,755,813,888 blocks.",
					},
					"used_capacity_mb": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Used capacity in MiB. Always 0 if volume is in a transitional state.",
					},
					"saving_setting": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Capacity reduction function setting.",
					},
					"is_data_reduction_share_enabled": {
						Type:        schema.TypeBool,
						Computed:    true,
						Description: "(Nullable) Shared data reduction setting. Only output if saving_setting != DISABLE.",
					},
					"compression_acceleration": {
						Type:        schema.TypeBool,
						Computed:    true,
						Description: "(Nullable) Compression accelerator enable/disable. Output if saving_setting != DISABLE.",
					},
					"capacity_saving_status": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Capacity reduction function status.",
					},
					"number_of_connecting_servers": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Number of connected servers.",
					},
					"number_of_snapshots": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Number of snapshots.",
					},
					"volume_types": {
						Type:        schema.TypeList,
						Computed:    true,
						Description: "Volume types.",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"volume_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of volumes returned.",
		},
		"total_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total number of volumes in storage.",
		},
	}
}

// ------------------- Volume Info Schema -------------------
func volumeInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"volume_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Volume ID.",
		},
		"volume_id_hex": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Volume ID in hexadecimal format.",
		},

		"nickname": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "(Nullable) Nickname.",
		},
		"pool_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Pool ID.",
		},
		"pool_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "(Nullable) Pool name. Not returned if the volume is in a transitional state.",
		},
		"total_capacity_mb": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total capacity in MiB.",
		},
		"used_capacity_mb": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Used capacity in MiB.",
		},
		"free_capacity_mb": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Free capacity in MiB.",
		},
		"reserved_capacity_mb": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Reserved capacity in MiB.",
		},
		"saving_setting": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Capacity reduction function setting.",
		},
		"is_data_reduction_share_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "(Nullable) Shared data reduction setting.",
		},
		"compression_acceleration": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "(Nullable) Compression accelerator enable/disable.",
		},
		"compression_acceleration_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "(Nullable) Actual compression state of the data stored in the volume.",
		},
		"capacity_saving_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Capacity reduction function status.",
		},
		"capacity_saving_progress": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "(Nullable) Capacity reduction migration progress (%)",
		},
		"number_of_connecting_servers": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of connected servers.",
		},
		"number_of_snapshots": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of snapshots.",
		},
		"luns": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "(Nullable) LUN list. Not displayed if volumeTypes is Namespace.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"lun": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "LU number",
					},
					"server_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Server ID",
					},
					"port_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Assigned port(s)",
					},
				},
			},
		},
		"volume_types": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Volume types.",
		},
	}
}

// ------------------- Datasource Get Volumes Schema -------------------
func datasourceAdminMultipleVolumesInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of the storage system",
		},
		"pool_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Filter volumes by Pool ID.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"pool_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Filter volumes by Pool Name (partial match).",
		},
		"server_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Filter volumes by Server ID.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"server_nickname": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Filter volumes by Server Nickname (partial match).",
		},
		"nickname": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Filter volumes by Volume Nickname (partial match).",
		},
		"min_total_capacity_mb": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Minimum volume size in MiB (47 to 268435456).",
			ValidateFunc: validation.IntBetween(47, 268435456),
		},
		"max_total_capacity_mb": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Maximum volume size in MiB (47 to 268435456).",
			ValidateFunc: validation.IntBetween(47, 268435456),
		},
		"min_used_capacity_mb": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Minimum used capacity in MiB (0 to 268435456).",
			ValidateFunc: validation.IntBetween(0, 268435456),
		},
		"max_used_capacity_mb": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Maximum used capacity in MiB (0 to 268435456).",
			ValidateFunc: validation.IntBetween(0, 268435456),
		},

		// --- Start Volume ID / Volume HEX (mutually exclusive) ---
		"start_volume_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Description:   "Starting Volume ID to display. Default is 0. Only one of start_volume_id or start_volume_id_hex may be specified, not both.",
			ConflictsWith: []string{"start_volume_id_hex"},
			ValidateFunc:  validation.IntBetween(0, 65535),
		},
		"start_volume_id_hex": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "Starting Volume ID (in hexadecimal format) to display. Default is 0. Only one of start_volume_id or start_volume_id_hex may be specified, not both.",
			ConflictsWith: []string{"start_volume_id"},
			// Validation: hex string, 1–4 hex chars (0x prefix optional)
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
				"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
			),
		},

		"requested_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      500,
			Description:  "Number of volumes to display. Default is 500.",
			ValidateFunc: validation.IntAtLeast(0),
		},
	}
}

func datasourceAdminMultipleVolumesOutputSchema() map[string]*schema.Schema {
	return volumesInfoListSchema()
}

func DatasourceAdminMultipleVolumesSchema() map[string]*schema.Schema {
	schema := datasourceAdminMultipleVolumesInputSchema()

	for k, v := range datasourceAdminMultipleVolumesOutputSchema() {
		schema[k] = v
	}

	return schema
}

// ------------------- Datasource Get One Volume Schema -------------------
func datasourceAdminOneVolumeInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of the storage system",
		},
		// --- Volume ID / Volume HEX (mutually exclusive) ---
		"volume_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Description:   "ID of the volume to retrieve. Only one of volume_id or volume_id_hex may be specified, not both.",
			ConflictsWith: []string{"volume_id_hex"},
			ValidateFunc:  validation.IntBetween(0, 65535),
		},
		"volume_id_hex": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "ID of the volume to retrieve in hexadecimal format. Only one of volume_id or volume_id_hex may be specified, not both.",
			ConflictsWith: []string{"volume_id"},
			// Validation: hex string, 1–4 hex chars (0x prefix optional)
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
				"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
			),
		},
	}
}

func datasourceAdminOneVolumeOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"volume_info": { // returns only one volume inside a list
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Volume Info returned from API",
			Elem: &schema.Resource{
				Schema: volumeInfoSchema(),
			},
		},
	}
}

func DatasourceAdminOneVolumeSchema() map[string]*schema.Schema {
	schema := datasourceAdminOneVolumeInputSchema()

	for k, v := range datasourceAdminOneVolumeOutputSchema() {
		schema[k] = v
	}

	return schema
}

// ------------------- Resource Volume Schema -------------------
func resourceAdminVolumeInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of the storage system",
		},
		"capacity": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Volume capacity with unit suffix (M, G, or T). Example: 500M, 100G, or 1T. Maximum 256T (268,435,456MiB). Minimum 47M.",
			ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
				s := v.(string)
				miB, err := utils.ParseCapacityToMiB(s)
				if err != nil {
					es = append(es, fmt.Errorf("%q must be a valid capacity with unit suffix (M, G, or T): %v", k, err))
					return
				}
				if miB < 47 || miB > 268435456 {
					es = append(es, fmt.Errorf("%q must be between 47MiB and 256TiB (268,435,456MiB)", k))
				}
				return
			},
		},

		"number_of_volumes": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Number of volumes to create. Default is 1.",
			ValidateFunc: validation.IntAtLeast(1),
		},

		"nickname_param": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "Nickname configuration. The combined base name and numeric suffix must not exceed 32 characters.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"base_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Base nickname (e.g., 'AAA').",
						ValidateFunc: validateBaseName,
					},
					"start_number": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1, // sentinel for "not set"
						Description:  "Starting number suffix. Only set if you want a numeric suffix.",
						ValidateFunc: validation.IntAtLeast(-1),
					},
					"number_of_digits": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      1,
						Description:  "Number of digits in the numeric suffix (applies only if start_number is specified).",
						ValidateFunc: validation.IntBetween(1, 32),
					},
				},
			},
		},

		"capacity_saving": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "DISABLE",
			Description: "Capacity reduction setting. Valid values: DEDUPLICATION_AND_COMPRESSION, COMPRESSION, DISABLE.",
			ValidateFunc: validation.StringInSlice([]string{
				"DEDUPLICATION_AND_COMPRESSION",
				"COMPRESSION",
				"DISABLE",
			}, false),
		},

		"is_data_reduction_share_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable shared data reduction. Only allowed if capacity_saving != DISABLE.",
		},

		"pool_id": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Pool ID in which the volume is created.",
			ValidateFunc: validation.IntAtLeast(0),
		},

		// --- Volume ID / Volume HEX (mutually exclusive) ---
		"volume_id": {
			Type:     schema.TypeInt,
			Optional: true,
			Description: `Specifies the ID of the volume to modify. Required for updates. Not allowed for create.

	- Only one of volume_id or volume_id_hex may be specified, not both.",
    - "number_of_volumes" must not be set; otherwise, an error will occur.
    - During update, only one volume is modified but all volume information is returned.
    - Note: Terraform destroy uses the state file, not the .tf file, so it deletes all tracked volumes.`,

			ConflictsWith: []string{"volume_id_hex"},
			ValidateFunc:  validation.IntBetween(0, 65535),
		},
		"volume_id_hex": {
			Type:     schema.TypeString,
			Optional: true,
			Description: `Specifies the ID of the volume to modify in hexadecimal format. Required for updates. Not allowed for create.

	- Only one of volume_id or volume_id_hex may be specified, not both.",
    - "number_of_volumes" must not be set; otherwise, an error will occur.
    - During update, only one volume is modified but all volume information is returned.
    - Note: Terraform destroy uses the state file, not the .tf file, so it deletes all tracked volumes.`,

			ConflictsWith: []string{"volume_id"},
			// Validation: hex string, 1–4 hex chars (0x prefix optional)
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
				"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
			),
		},

		"compression_acceleration": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable or disable compression acceleration. Optional on update; ignored on create.",
		},
	}
}

func resourceAdminVolumeOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"volume_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "Number of volumes created.",
		},
		"volumes_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Volume Info returned from API",
			Elem: &schema.Resource{
				Schema: volumeInfoSchema(),
			},
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Comma-separated list of created Volume IDs.",
		},
	}
}

func ResourceAdminVolumeSchema() map[string]*schema.Schema {
	schema := resourceAdminVolumeInputSchema()

	for k, v := range resourceAdminVolumeOutputSchema() {
		schema[k] = v
	}

	return schema
}

// ------------------- Helpers -------------------
func validateBaseName(v interface{}, k string) (ws []string, es []error) {
	baseName := v.(string)
	length := len(baseName)

	if length < 1 || length > 32 {
		es = append(es, fmt.Errorf(
			"%q must be between 1 and 32 characters long, got %q (length %d)",
			"nickname_param.0.base_name", baseName, length,
		))
	}

	return
}
