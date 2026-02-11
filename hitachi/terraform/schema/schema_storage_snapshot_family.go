package terraform

import (
	"regexp"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ------------------- Datasource Virtual Clone Parent Volume Schema -------------------

func datasourceVspVirtualCloneParentVolumeInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(1),
		},
	}
}

func datasourceVspVirtualCloneParentVolumeOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"parent_volumes": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of LDEVs that are currently acting as parents for virtual clones (TIA).",
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"parent_volumes_hex": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of LDEVs in hex that are currently acting as parents for virtual clones (TIA).",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"parent_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of virtual clone parent volumes found.",
		},
	}
}

func DatasourceVspVirtualCloneParentVolumeSchema() map[string]*schema.Schema {
	s := datasourceVspVirtualCloneParentVolumeInputSchema()

	for k, v := range datasourceVspVirtualCloneParentVolumeOutputSchema() {
		s[k] = v
	}

	return s
}

// ------------------- Datasource Snapshot Family Schema -------------------

func datasourceVspSnapshotFamilyInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"ldev_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Description:  "The LDEV ID of the volume to query the snapshot family for. Only one of ldev_id or ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"ldev_id_hex"},
			ValidateFunc:  validation.IntBetween(0, 65535),
		},
		"ldev_id_hex": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:  "The LDEV ID in hex of the volume to query the snapshot family for. Only one of ldev_id or ldev_id_hex may be specified, not both.",
			ConflictsWith: []string{"ldev_id"},
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
				"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
			),
		},
	}
}

func datasourceVspSnapshotFamilyOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"family_members": snapshotFamilySchema(),
		"total_members": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The total number of volumes found in this snapshot family.",
		},
	}
}

func DatasourceVspSnapshotFamilySchema() map[string]*schema.Schema {
	s := datasourceVspSnapshotFamilyInputSchema()

	for k, v := range datasourceVspSnapshotFamilyOutputSchema() {
		s[k] = v
	}

	return s
}

func snapshotFamilySchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "List of volumes belonging to the snapshot family, including parents and virtual clones.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ldev_id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "LDEV number of the volume in the family.",
				},
				"ldev_id_hex": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "LDEV number of the volume in the family in hexadecimal format",
				},
				"snapshot_group_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the snapshot group.",
				},
				"primary_or_secondary": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Indicates if the volume is a P-VOL or S-VOL.",
				},
				"status": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Current status of the snapshot pair.",
				},
				"pvol_ldev_id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "LDEV number of the Primary Volume.",
				},
				"svol_ldev_id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "LDEV number of the Secondary Volume.",
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
				"mirror_unit_id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Mirror Unit (MU) number.",
				},
				"pool_id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Pool ID associated with the snapshot.",
				},
				"is_virtual_clone_volume": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "True if the volume is a virtual clone (V-Clone).",
				},
				"is_virtual_clone_parent_volume": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "True if the volume is a parent of a virtual clone.",
				},
				"split_time": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The time when the snapshot was split.",
				},
				"parent_ldev_id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "LDEV ID of the immediate parent in the cascade.",
				},
				"snapshot_group_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Internal ID of the snapshot group.",
				},
				"snapshot_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The unique identifier for the snapshot (usually Pvol,Mu).",
				},
			},
		},
	}
}
