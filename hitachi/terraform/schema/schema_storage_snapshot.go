package terraform

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ------------------- Snapshots Info Schema (SnapshotAll) -------------------
func snapshotsRangeListSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"snapshots": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of all Thin Image pairs in the storage system.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"snapshot_replication_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The unique identifier of the snapshot replication (Format: pvolLdevId,muNumber).",
					},
					"pvol_ldev_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "P-VOL LDEV ID.",
					},
					"pvol_ldev_id_hex": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "P-VOL LDEV ID in hexadecimal format",
					},
					"mirror_unit_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "MU number of the P-VOL.",
					},
					"snapshot_group_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of the snapshot group.",
					},
					"snapshot_pool_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "ID of the pool where snapshot data was created.",
					},
					"svol_ldev_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "S-VOL LDEV ID.",
					},
					"svol_ldev_id_hex": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "S-VOL LDEV ID in hexadecimal format",
					},
					"consistency_group_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Consistency group ID.",
					},
					"status": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Status of the pair (e.g., PAIR, PSUS, COPY).",
					},
					"concordance_rate": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Concordance rate for pairs.",
					},
					"is_redirect_on_write": {
						Type:        schema.TypeBool,
						Computed:    true,
						Description: "True if the pair is Thin Image Advanced.",
					},
					"is_clone": {
						Type:        schema.TypeBool,
						Computed:    true,
						Description: "True if the pair has the clone attribute.",
					},
					"can_cascade": {
						Type:        schema.TypeBool,
						Computed:    true,
						Description: "True if the pair can be a cascaded pair.",
					},
					"split_time": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Time when snapshot data was created.",
					},
				},
			},
		},
		"snapshot_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of snapshots returned.",
		},
	}
}

// ------------------- Snapshot Detail Schema (Snapshot) -------------------
func snapshotDetailSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pvol_ldev_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "LDEV number of the Primary Volume (P-VOL).",
		},
		"pvol_ldev_id_hex": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "P-VOL LDEV ID in hexadecimal format",
		},
		"mirror_unit_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "MU number of the P-VOL.",
		},
		"svol_ldev_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "LDEV number of the Secondary Volume (S-VOL).",
		},
		"svol_ldev_id_hex": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "S-VOL LDEV ID in hexadecimal format",
		},
		"snapshot_pool_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the pool where snapshot data is created.",
		},
		"snapshot_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Identifier in the format 'pvolLdevId,muNumber'.",
		},
		"snapshot_group_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the snapshot group.",
		},
		"primary_or_secondary": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Attribute of the LDEV (P-VOL or S-VOL).",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Current pair status.",
		},
		"is_redirect_on_write": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if the pair is Thin Image Advanced.",
		},
		"is_consistency_group": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if pair is created in CTG mode.",
		},
		"is_written_in_svol": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if data was written to the S-VOL when status was PSUS/PFUS.",
		},
		"is_clone": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if the pair has the clone attribute.",
		},
		"can_cascade": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if the pair can be a cascaded pair.",
		},
		"snapshot_data_read_only": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if the snapshot data has the read-only attribute.",
		},
		"concordance_rate": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Concordance rate between P-VOL and S-VOL.",
		},
		"progress_rate": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Progress of processing (%).",
		},
		"split_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Timestamp when snapshot data was created.",
		},
		"pvol_processing_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Processing status of the P-VOL pair (E: In Progress, N: Not In Progress).",
		},
		"svol_processing_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Processing status of the S-VOL pair (E: In Progress, N: Not In Progress).",
		},
		"is_virtual_clone_parent_volume": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if the volume has the vClone Parent attribute. For VSP One B20/B85 only.",
		},
		"is_virtual_clone_volume": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if the volume has the vClone attribute or vClone Parent attribute. For VSP One B20/B85 only",
		},
		"retention_period_hours": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The remaining period for which snapshot data of the Thin Image Advanced pair is retained (in hours) if set.",
		},
	}
}

// ------------------- Datasource Get Snapshot Schema -------------------

func datasourceVspSnapshotInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"snapshot_group_name": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Name of the snapshot group (1-32 characters).",
			ValidateFunc: validation.StringLenBetween(1, 32),
		},
		"pvol_ldev_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Description:   "P-VOL LDEV ID. Only one of pvol_ldev_id or pvol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"pvol_ldev_id_hex"},
			ValidateFunc:  validation.IntBetween(0, 65535),
		},
		"pvol_ldev_id_hex": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "P-VOL LDEV ID in hexadecimal format. Only one of pvol_ldev_id or pvol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"pvol_ldev_id"},
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
				"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
			),
		},
		"svol_ldev_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Description:   "S-VOL LDEV ID. Only one of svol_ldev_id or svol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"svol_ldev_id_hex"},
			ValidateFunc:  validation.IntBetween(0, 65535),
		},
		"svol_ldev_id_hex": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "S-VOL LDEV ID in hexadecimal format. Only one of svol_ldev_id or svol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"svol_ldev_id"},
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
				"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
			),
		},
		"mirror_unit_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "MU number of the P-VOL.",
			ValidateFunc: validation.IntAtLeast(0),
		},
	}
}

func datasourceVspSnapshotOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"snapshots": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of snapshots matching the criteria.",
			Elem: &schema.Resource{
				Schema: snapshotDetailSchema(),
			},
		},
		"snapshot_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of snapshots returned.",
		},
	}
}

func DatasourceVspSnapshotSchema() map[string]*schema.Schema {
	s := datasourceVspSnapshotInputSchema()

	for k, v := range datasourceVspSnapshotOutputSchema() {
		s[k] = v
	}

	return s
}

// ------------------- Datasource Get Snapshots Range Schema -------------------

func datasourceVspSnapshotRangeInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system. Only VSP 5000 series is supported.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"start_pvol_ldev_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Default:       0,
			Description:   "The starting P-VOL LDEV number of the range. Defaults to 0. Only one of start_pvol_ldev_id or start_pvol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"start_pvol_ldev_id_hex"},
			ValidateFunc:  validation.IntBetween(0, 65535),
		},
		"start_pvol_ldev_id_hex": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "The starting P-VOL LDEV number of the range in hexadecimal format. Defaults to 0. Only one of start_pvol_ldev_id or start_pvol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"start_pvol_ldev_id"},
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
				"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
			),
		},
		"end_pvol_ldev_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Description:   "The ending P-VOL LDEV number of the range. Defaults to the maximum LDEV number. Only one of end_pvol_ldev_id or end_pvol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"end_pvol_ldev_id_hex"},
			ValidateFunc:  validation.IntBetween(0, 65535),
		},

		"end_pvol_ldev_id_hex": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "The ending P-VOL LDEV number of the range in hexadecimal format. Defaults to the maximum LDEV number. Only one of end_pvol_ldev_id or end_pvol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"end_pvol_ldev_id"},
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
				"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
			),
		},
	}
}

func datasourceVspSnapshotRangeOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"snapshots": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of snapshots matching the criteria.",
			Elem: &schema.Resource{
				Schema: snapshotDetailSchema(),
			},
		},
		"snapshot_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of snapshots returned.",
		},
	}
}

func DatasourceVspSnapshotRangeSchema() map[string]*schema.Schema {
	s := datasourceVspSnapshotRangeInputSchema()

	for k, v := range datasourceVspSnapshotRangeOutputSchema() {
		s[k] = v
	}

	return s
}

// ------------------- Resource Vsp Snapshot Input Schema -------------------

func resourceVspSnapshotInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Storage system serial number.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"state": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "read",
			Description: "read, create, split, resync, restore, clone (TI Std), vclone (TIA), vrestore (TIA), defrag (TI Std), deletetree (TI Std)",
			ValidateFunc: validation.StringInSlice([]string{
				"", "read", "create", "split", "resync", "restore", "clone", "vclone", "vrestore", "defrag", "deletetree",
			}, true),
		},
		"snapshot_group_name": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Snapshot group name (1-32 characters), case sensitive. Required for create. If you specify a new group name, a snapshot group is also created at the same time.",
			ValidateFunc: validation.StringLenBetween(1, 32),
		},
		"snapshot_pool_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Pool ID for snapshot data to be created. Use HTI for TI Standard. TIA requires the same HDP pool as the P-VOL. Required for create.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"pvol_ldev_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Description:   "P-VOL LDEV ID. For TIA, must be a DRS volume and reside in the snapshot_pool_id. Only one of pvol_ldev_id or pvol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"pvol_ldev_id_hex"},
			ValidateFunc:  validation.IntBetween(0, 65535),
		},
		"pvol_ldev_id_hex": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "P-VOL LDEV ID in hexadecimal format. For TIA, must be a DRS volume and reside in the snapshot_pool_id. Only one of pvol_ldev_id or pvol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"pvol_ldev_id"},
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
				"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
			),
		},
		"svol_ldev_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Description:   "S-VOL LDEV ID. For TIA, must be a DRS volume and reside in the snapshot_pool_id. Required if is_clone or can_cascade is true. Pair actions (split, resync, etc.) require an S-VOL. Only one of svol_ldev_id or svol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"svol_ldev_id_hex"},
			ValidateFunc:  validation.IntBetween(0, 65535),
		},
		"svol_ldev_id_hex": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "S-VOL LDEV ID. For TIA, must be a DRS volume and reside in the snapshot_pool_id. Required if is_clone or can_cascade is true. Pair actions (split, resync, etc.) require an S-VOL. Only one of svol_ldev_id or svol_ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"svol_ldev_id"},
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
				"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
			),
		},
		"mirror_unit_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Mirror Unit (MU) number. Assigned automatically during create if omitted.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"auto_split": {
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			Description:   "Automatically split the pair during create, resync, or restore. Incompatible with CTG or clones. Cannot set with is_consistency_group, is_clone, or auto_clone.",
			ConflictsWith: []string{"is_consistency_group", "is_clone", "auto_clone"},
		},
		"is_consistency_group": {
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			Description:   "Enable snapshot group in Consistency Group (CTG) mode. Restricts individual split actions.",
			ConflictsWith: []string{"auto_split"},
		},
		"is_clone": {
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			Description:   "Enable the clone attribute. Must be true to use the 'clone' action. Only for TI Std. Must not specify auto_split or is_consistency_group. Requires can_cascade to be true",
			ConflictsWith: []string{"auto_split"},
		},
		"auto_clone": {
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			Description:   "Automatically clone the pair after creation. Requires is_clone to be true.",
			ConflictsWith: []string{"auto_split"},
		},
		"copy_speed": {
			Type:     schema.TypeString,
			Optional: true,
			// don't set a default here; let the API apply its own default
			Description:  "Copy speed at which the created pair is to be cloned: slower, medium, or faster. Requires is_clone and auto_clone to be true.",
			ValidateFunc: validation.StringInSlice([]string{"slower", "medium", "faster"}, true),
		},
		"can_cascade": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "True if the pair can be a cascaded pair. Default is true.",
		},
		"is_data_reduction_force_copy": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Must be true if P-VOL has capacity saving enabled.",
		},
		"retention_period_hours": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The period for which snapshot data of the Thin Image Advanced pair is to be retained in hours (1-12288). Requires auto_split to be true.",
			ValidateFunc: validation.IntBetween(0, 12288),
		},
		"defrag_operation": {
			Type:     schema.TypeString,
			Optional: true,
			Description: "Needed for state defrag. Starts or stops the deletion of garbage data (defragmentation) for the snapshot tree. " +
				"Valid only for TI Standard on VSP 5000 series; not supported for TIA. " +
				"Options: 'start' to begin defragmentation, 'stop' to cancel a running operation.",
			ValidateFunc: validation.StringInSlice([]string{"start", "stop"}, true),
		},
	}
}

func resourceVspSnapshotOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"snapshot": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem:     &schema.Resource{Schema: snapshotDetailSchema()},
		},
		"vclone": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem:     &schema.Resource{Schema: vcloneDetailSchema()},
		},
		"additional_info": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem:     &schema.Resource{Schema: snapshotAdditionalInfoSchema()},
		},
	}
}

func ResourceVspSnapshotSchema() map[string]*schema.Schema {
	schemaMap := resourceVspSnapshotInputSchema()

	for k, v := range resourceVspSnapshotOutputSchema() {
		schemaMap[k] = v
	}

	return schemaMap
}

func snapshotAdditionalInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Storage system serial number.",
		},
		"is_thin_image_advanced": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Indicates if the pair uses Thin Image Advanced (Redirect-on-Write) architecture.",
		},
		"snapshot_pool_type": { // Fixed typo: snaphot -> snapshot
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The storage pool type used for snapshot data (e.g., HDP, HTI).",
		},
		"pvol_attributes": {
			Type:        schema.TypeList, // Changed to array of strings
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of attributes associated with the Primary Volume (P-VOL).",
		},
		"svol_attributes": {
			Type:        schema.TypeList, // Changed to array of strings
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of attributes associated with the Secondary Volume (S-VOL).",
		},
	}
}

func vcloneDetailSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// --- Identity (Hardware Confirmation) ---
		"pvol_ldev_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The LDEV number of the Parent volume (VCP).",
		},
		"svol_ldev_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The LDEV number of the Clone volume (VC).",
		},
		"pvol_ldev_id_hex": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "P-VOL LDEV ID in hexadecimal format",
		},
		"svol_ldev_id_hex": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "S-VOL LDEV ID in hexadecimal format",
		},

		// --- Attribute Flags ---
		"is_virtual_clone_volume": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if this LDEV is a promoted Secondary volume (VC).",
		},
		"is_virtual_clone_parent_volume": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if the Primary LDEV is acting as a vClone Parent (VCP).",
		},

		// --- Metadata ---
		"split_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The timestamp when the vClone relationship was established or promoted.",
		},
		"pool_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The Pool ID where the vClone data blocks reside.",
		},
	}
}
