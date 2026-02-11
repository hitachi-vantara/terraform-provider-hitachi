package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ------------------- SnapshotGroup Info Schema -------------------
func snapshotGroupDetailSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"snapshot_group_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The name of the snapshot group.",
		},
		"snapshot_group_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique ID of the snapshot group.",
		},
		"snapshots": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Detailed list of snapshot pairs within this group.",
			Elem: &schema.Resource{
				Schema: snapshotDetailSchema(), // Reuses your existing snapshotDetailSchema
			},
		},
	}
}

// ------------------- Datasource Get SnapshotGroup Schema -------------------

func datasourceVspSnapshotGroupInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"snapshot_group_name": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Filter by a specific snapshot group name.",
			ValidateFunc: validation.StringLenBetween(1, 32),
		},
	}
}

func datasourceVspSnapshotGroupOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"snapshot_group": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of snapshot groups matching the criteria.",
			Elem: &schema.Resource{
				Schema: snapshotGroupDetailSchema(),
			},
		},
		"snapshot_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of snapshot pairs returned.",
		},
	}
}

func DatasourceVspSnapshotGroupSchema() map[string]*schema.Schema {
	s := datasourceVspSnapshotGroupInputSchema()
	for k, v := range datasourceVspSnapshotGroupOutputSchema() {
		s[k] = v
	}
	return s
}

// ------------------- Datasource Get SnapshotGroups Schema -------------------

func datasourceVspMultipleSnapshotGroupsInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"include_pairs": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether to include snapshot pair details within each group.",
		},
	}
}

func datasourceVspMultipleSnapshotGroupsOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"snapshot_groups": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of snapshot groups matching the criteria.",
			Elem: &schema.Resource{
				Schema: snapshotGroupDetailSchema(),
			},
		},
		"snapshotgroup_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of snapshot groups returned.",
		},
		"snapshot_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of snapshot pairs returned.",
		},
	}
}

func DatasourceVspMultipleSnapshotGroupsSchema() map[string]*schema.Schema {
	s := datasourceVspMultipleSnapshotGroupsInputSchema()
	for k, v := range datasourceVspMultipleSnapshotGroupsOutputSchema() {
		s[k] = v
	}
	return s
}

// ------------------- Resource Vsp Snapshot Input Schema -------------------

func resourceVspSnapshotGroupInputSchema() map[string]*schema.Schema {
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
			Description: "Action for the group. Options: read, split, resync, restore, clone, vclone, vrestore, update_retention.",
			ValidateFunc: validation.StringInSlice([]string{
				"", "read", "split", "resync", "restore", "clone", "vclone", "vrestore", "update_retention",
			}, true),
		},
		"snapshot_group_name": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "The name of the snapshot group to act upon.",
			ValidateFunc: validation.StringLenBetween(1, 32),
		},
		"auto_split": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Automatically split the pairs units of snapshot groups after resync or restore.",
		},
		"copy_speed": {
			Type:     schema.TypeString,
			Optional: true,
			// don't set a default here; let the API apply its own default
			Description:  "Copy speed at which the created pair is to be cloned: slower, medium, or faster.",
			ValidateFunc: validation.StringInSlice([]string{"slower", "medium", "faster"}, true),
		},
		"retention_period_hours": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The period for which snapshot data of the Thin Image Advanced pair is to be retained in hours (1-12288). Requires auto_split to be true.",
			ValidateFunc: validation.IntBetween(0, 12288),
		},
	}
}

func resourceVspSnapshotGroupOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"snapshot_group": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "The state of the snapshot group after the action.",
			Elem: &schema.Resource{
				Schema: snapshotGroupDetailSchema(),
			},
		},
		"vclones": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem:     &schema.Resource{Schema: vcloneDetailSchema()},
		},
	}
}

func ResourceVspSnapshotGroupSchema() map[string]*schema.Schema {
	schemaMap := resourceVspSnapshotGroupInputSchema()
	for k, v := range resourceVspSnapshotGroupOutputSchema() {
		schemaMap[k] = v
	}
	return schemaMap
}
