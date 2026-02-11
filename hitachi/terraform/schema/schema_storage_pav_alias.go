package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var PavAliasSchema = map[string]*schema.Schema{
	"cu_number": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "CU number",
	},
	"ldev_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "LDEV number",
	},
	"pav_attribute": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "PAV attribute of LDEV (BASE or ALIAS)",
	},
	"c_base_volume_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "LDEV number of the base device assigned within the storage system (only for alias devices)",
	},
	"s_base_volume_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "LDEV number of the base device set by the user (only for alias devices)",
	},
}

var DataPavAliasSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"cu_number": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "CU number (0-255). If omitted, returns all PAV aliases.",
		ConflictsWith: []string{
			"base_ldev_id",
			"alias_ldev_ids",
		},
		ValidateFunc: validation.IntBetween(0, 255),
	},
	"base_ldev_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Base LDEV ID. When set, the provider returns alias entries whose c_base_volume_id (or s_base_volume_id) matches this base LDEV.",
		ConflictsWith: []string{
			"cu_number",
			"alias_ldev_ids",
		},
		ValidateFunc: validation.IntBetween(0, 65279),
	},
	"alias_ldev_ids": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Alias LDEV IDs. When set, the provider returns entries whose ldev_id matches any of the provided IDs.",
		ConflictsWith: []string{
			"cu_number",
			"base_ldev_id",
		},
		Elem: &schema.Schema{
			Type:         schema.TypeInt,
			ValidateFunc: validation.IntBetween(0, 65279),
		},
	},
	"pav_aliases": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "PAV alias information.",
		Elem:        &schema.Resource{Schema: PavAliasSchema},
	},
}

// ResourceVspPavLdevSchema defines schema for hitachi_vsp_pav_ldev resource.
var ResourceVspPavLdevSchema = map[string]*schema.Schema{
	"serial": {
		Type:     schema.TypeInt,
		Required: true,

		Description: "Serial number of the storage system",
	},
	"base_ldev_id": {
		Type:         schema.TypeInt,
		Required:     true,
		Description:  "LDEV number of the base device.",
		ValidateFunc: validation.IntBetween(0, 65279),
	},
	"alias_ldev_ids": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		MaxItems:    255,
		Description: "LDEV numbers to be assigned/unassigned as alias devices. Specify LDEVs that have the same CU number as the base device.",
		Elem: &schema.Schema{
			Type:         schema.TypeInt,
			ValidateFunc: validation.IntBetween(0, 65279),
		},
	},
	"pav_aliases": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "PAV alias information for the configured base/alias LDEVs.",
		Elem:        &schema.Resource{Schema: PavAliasSchema},
	},
}
