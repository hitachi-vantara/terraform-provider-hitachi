package terraform

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func nvmSubsystemDetailSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"nvm_subsystem_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "NVM subsystem ID.",
		},
		"virtual_nvm_subsystem_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Virtual NVM subsystem ID. Output for VSP 5000 series storage systems.",
		},
		"nvm_subsystem_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "NVM subsystem name.",
		},
		"resource_group_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Resource group ID of the resource group to which the NVM subsystem belongs.",
		},
		"namespace_security_setting": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Namespace security settings (Enable, Disable, or Unknown).",
		},
		"t10_pi_mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Status of the T10 PI mode of the port (Enable, Disable, or Unknown).",
		},
		"host_mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Host mode of the NVM subsystem.",
		},
		"host_mode_options": {
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Description: "List of host mode option numbers set for the NVM subsystem.",
		},
		"nvm_subsystem_nqn": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Subsystem NQN. If virtualized, outputs the virtualized NQN.",
		},
		"namespaces": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of namespaces created in the NVM subsystem.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"namespace_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The Namespace ID.",
					},
					"ldev_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "LDEV number associated with the namespace.",
					},
					"ldev_id_hex": {
						Computed:    true,
						Type:        schema.TypeString,
						Description: "LDEV hexadecimal number associated with the namespace.",
					},
					"namespace_nickname": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Nickname for the namespace.",
					},
					"paths": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of Host NQNs that have access to this namespace.",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"port_ids": {
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of registered port IDs (e.g., CL1-A).",
		},
		"host_nqns": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of Host NQNs registered to this subsystem.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host_nqn_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The unique identifier for the Host NQN.",
					},
					"host_nqn": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The NVMe Qualified Name (NQN) of the host.",
					},
					"host_nqn_nickname": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The nickname assigned to the Host NQN.",
					},
				},
			},
		},
	}
}

func DatasourceVspNvmSubsystemSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"nvm_subsystem_id": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "NVM subsystem ID.",
			ValidateFunc: validation.IntBetween(0, 65535),
		},
		// --- Output ---
		"nvm_subsystems": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "NVM subsystem info.",
			Elem: &schema.Resource{
				Schema: nvmSubsystemDetailSchema(),
			},
		},
		"nvm_subsystem_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total number of NVM subsystems returned.",
		},
	}
}

func DatasourceVspNvmSubsystemsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"nvm_subsystem_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Filter by NVM subsystem ID.",
			ValidateFunc: validation.IntBetween(0, 65535),
		},
		"nvm_subsystem_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Filter by NVM subsystem name (case-insensitive).",
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 32),
				validation.StringDoesNotMatch(regexp.MustCompile(`^-|^\s|\s$`), "name cannot start with a hyphen/space or end with a space"),
			),
		},
		"include_ports": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, include ports information in the output.",
		},
		"include_host_nqns": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, include host NQNs and host NQN paths information in the output.",
		},
		"include_namespaces": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, include namespaces information in the output.",
		},
		// --- Output ---
		"nvm_subsystems": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of NVMe subsystems matching the filter criteria.",
			Elem: &schema.Resource{
				Schema: nvmSubsystemDetailSchema(),
			},
		},
		"nvm_subsystem_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of NVM subsystems returned.",
		},
	}
}

func ResourceVspNvmSubsystemSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		// -------------------------------------------------
		// Identity Fields
		// -------------------------------------------------

		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true, // identity
			Description: "Storage system serial number.",
		},

		"nvm_subsystem_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true, // allows auto-assignment
			ForceNew:     true, // changing it recreates resource
			Description:  "NVMe ID (0-2047). If omitted, the lowest available ID is used.",
			ValidateFunc: validation.IntBetween(0, 2047),
		},

		// -------------------------------------------------
		// Configurable Fields
		// -------------------------------------------------

		"virtual_nvm_subsystem_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true, // keeps import safe
			Description:  "Virtual NVM subsystem ID (0-2047). Applicable for VSP 5000 series.",
			ValidateFunc: validation.IntBetween(0, 2047),
		},

		"nvm_subsystem_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true, // important for safe import
			Description: "NVM subsystem name (1-32 characters). Cannot start/end with space or hyphen.",
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 32),
				validation.StringDoesNotMatch(regexp.MustCompile(`^-|^\s|\s$`),
					"name cannot start with a hyphen/space or end with a space"),
				validation.StringDoesNotMatch(
					regexp.MustCompile(`(?i)^nvmss_id_.*\(default_name\)$`),
					"reserved default name pattern not allowed"),
			),
		},

		"host_mode": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true, // prevents import diff
			Description: "Host mode. Options: LINUX/IRIX, VMWARE, VMWARE_EX, AIX.",
			ValidateFunc: validation.StringInSlice(
				[]string{"LINUX/IRIX", "VMWARE", "VMWARE_EX", "AIX"},
				false,
			),
		},

		"host_mode_options": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true, // critical for import safety
			Description: "List of host mode option numbers.",
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},

		"enable_namespace_security": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true, // avoids default drift
			Description: "Enable or disable namespace security.",
		},

		"ports": {
			Type:        schema.TypeSet, // Set avoids ordering diffs
			Optional:    true,
			Computed:    true,
			Description: "List of port IDs (FC-NVMe or NVMe over TCP).",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"host_nqns": {
			Type:        schema.TypeSet, // Use Set instead of List
			Optional:    true,
			Computed:    true,
			Description: "List of Host NQNs to register.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nqn": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The NVMe Qualified Name (NQN) of the host.",
					},
					"nickname": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "An optional user-defined alias for the Host NQN to simplify identification in the storage array.",
					},
				},
			},
		},

		"namespaces": {
			Type:        schema.TypeSet, // Set avoids order diff
			Optional:    true,
			Computed:    true,
			Description: "List of namespaces within the subsystem.",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ldev_id": {
						Type:         schema.TypeInt,
						Optional:     true,
						Default:      -1,
						Description:  "LDEV ID to be presented as an NVMe namespace. Only one of ldev_id or ldev_id_hex may be specified, not both.",
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"ldev_id_hex": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "LDEV ID in hexadecimal to be presented as an NVMe namespace. Only one of ldev_id or ldev_id_hex may be specified, not both.",
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^(0[xX])?[A-Fa-f0-9]{1,4}$`),
							"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
						),
					},
					"nickname": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "A descriptive name for the namespace, often used to identify the volume's purpose.",
					},
					"path": {
						Type:        schema.TypeSet, // avoid ordering diff
						Optional:    true,
						Computed:    true,
						Description: "List of Host NQNs to map to this namespace.",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},

		// -------------------------------------------------
		// Output-Only
		// -------------------------------------------------

		"nvm_subsystems": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "NVM subsystem info.",
			Elem: &schema.Resource{
				Schema: nvmSubsystemDetailSchema(),
			},
		},
	}
}
