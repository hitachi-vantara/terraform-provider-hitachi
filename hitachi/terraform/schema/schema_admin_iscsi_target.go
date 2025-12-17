package terraform

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ------------------- iSCSI Targets Info Schema -------------------
func iscsiTargetsInfoListSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"iscsi_targets_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of iSCSI targets.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"port_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Assigned port ID.",
					},
					"target_iscsi_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "(nullable) Target port iSCSI name (for iSCSI).",
					},
				},
			},
		},
		"iscsi_targets_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of iSCSI target ports returned.",
		},
	}
}

// ------------------- iSCSI Target Info Schema -------------------
func iscsiTargetInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Assigned port ID.",
		},
		"target_iscsi_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "(nullable) Target port iSCSI name (for iSCSI).",
		},
	}
}

// ------------------- Datasource: Multiple iSCSI Targets -------------------
func datasourceAdminMultipleIscsiTargetsInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"server_id": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "The ID of the server whose iSCSI target ports are to be retrieved.",
			ValidateFunc: validation.IntAtLeast(0),
		},
	}
}

func datasourceAdminMultipleIscsiTargetsOutputSchema() map[string]*schema.Schema {
	return iscsiTargetsInfoListSchema()
}

func DatasourceAdminMultipleIscsiTargetsSchema() map[string]*schema.Schema {
	s := datasourceAdminMultipleIscsiTargetsInputSchema()
	for k, v := range datasourceAdminMultipleIscsiTargetsOutputSchema() {
		s[k] = v
	}
	return s
}

// ------------------- Datasource: One iSCSI Target -------------------
func datasourceAdminOneIscsiTargetInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"server_id": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "The ID of the server.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"port_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The port ID of the iSCSI target (e.g., CL1-A, CL2-B).",
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Za-z0-9-]+$`),
				"port_id must contain only letters, numbers, and hyphens (e.g., CL1-A)",
			),
		},
	}
}

func datasourceAdminOneIscsiTargetOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"iscsi_target_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "iSCSI Target info returned from API.",
			Elem: &schema.Resource{
				Schema: iscsiTargetInfoSchema(),
			},
		},
	}
}

func DatasourceAdminOneIscsiTargetSchema() map[string]*schema.Schema {
	s := datasourceAdminOneIscsiTargetInputSchema()
	for k, v := range datasourceAdminOneIscsiTargetOutputSchema() {
		s[k] = v
	}
	return s
}

// ------------------- Resource: iSCSI Target -------------------
func resourceAdminIscsiTargetInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"server_id": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "The ID of the server.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"port_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The port ID of the iSCSI target (e.g., CL1-A, CL2-B).",
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Za-z0-9-]+$`),
				"port_id must contain only letters, numbers, and hyphens (e.g., CL1-A)",
			),
		},
		"target_iscsi_name": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "For changing target port iSCSI name.",
			ValidateFunc: validation.StringLenBetween(1, 255),
		},
	}
}

func resourceAdminIscsiTargetOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"iscsi_target_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "iSCSI Target info returned from API.",
			Elem: &schema.Resource{
				Schema: iscsiTargetInfoSchema(),
			},
		},
	}
}

func ResourceAdminIscsiTargetSchema() map[string]*schema.Schema {
	s := resourceAdminIscsiTargetInputSchema()
	for k, v := range resourceAdminIscsiTargetOutputSchema() {
		s[k] = v
	}
	return s
}
