package terraform

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)
var IsUUID = validation.StringMatch(uuidRegex, "must be a valid UUID")

var ResourceVssbConfigurationFileSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The address (IP or hostname) of the VOSB system's REST API.",
	},

	"download_existconfig_only": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If true, skips creation and only downloads the latest existing configuration file. Requires `download_path`. All other parameters are ignored.",
	},

	"create_only": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If true, creates a new configuration file but does not download it. Ignored if `download_existconfig_only` is true.",
	},

	"download_path": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Path to save the downloaded configuration file. Ignored if no download occurs. Can be a directory or a specific file path.",
	},

	"create_configuration_file_param": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Parameters for creating a configuration file (relevant for Google Cloud and Azure only). Ignored for AWS and Bare Metal.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{

				"expected_cloud_provider": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "baremetal",
					Description: `Specifies the expected cloud provider type. Valid values: "google", "azure", "aws", "baremetal".

	- Used to validate combinations of inputs based on the deployment environment.
	- If set to "google" or "azure", specific parameters may be required for certain operations.
	- If set to "aws" or "baremetal" (default), other cloud-specific inputs are ignored. These behave identically.
	- Note: The actual cloud provider is determined by the VOSB system at the "vosb_address" endpoint.
	If there's a mismatch, the request still proceeds and behaves according to the actual environment.`,
					ValidateFunc: validation.StringInSlice([]string{
						"google", "azure", "aws", "baremetal",
					}, false),
				},

				"export_file_type": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "Normal",
					Description: `Specifies the type of configuration file to generate. Default: 'Normal'.
	Valid values: Normal, AddStorageNodes, ReplaceStorageNode, AddDrives, ReplaceDrive.

	Determines which additional parameters are needed:
	- Normal: All other parameters are ignored.
	- AddStorageNodes: Requires 'machine_image_id' and 'address_setting'.
	- ReplaceStorageNode:
		- Google Cloud: Requires 'machine_image_id' and 'node_id'. Optionally 'recover_single_node'.
		- Azure: Requires 'machine_image_id'.
	- AddDrives: Requires 'number_of_drives'.
	- ReplaceDrive (Google Cloud only): Requires 'drive_id' or 'recover_single_drive'.

	Note:
	- Ignored in AWS and Bare Metal environments.
	- Used only in Google Cloud or Azure to control behavior.`,
					ValidateFunc: validation.StringInSlice([]string{
						"Normal",
						"AddStorageNodes",
						"ReplaceStorageNode",
						"AddDrives",
						"ReplaceDrive",
					}, false),
				},

				"machine_image_id": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: "ID or URI of the VM image used in AddStorageNodes or ReplaceStorageNode operations.",
				},

				"number_of_drives": &schema.Schema{
					Type:         schema.TypeInt,
					Optional:     true,
					Description:  "Number of drives to install per node in AddDrives. Must be between 6 and 24.",
					ValidateFunc: validation.IntBetween(6, 24),
				},

				"recover_single_drive": &schema.Schema{
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Whether to recover a removed drive during a ReplaceDrive operation.",
				},

				"drive_id": &schema.Schema{
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "UUID of the drive to replace. Must not be set if `recover_single_drive` is true.",
					ValidateFunc: IsUUID,
				},

				"recover_single_node": &schema.Schema{
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Whether to recover a storage node during a ReplaceStorageNode operation.",
				},

				"node_id": &schema.Schema{
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "UUID of the storage node to replace. Required for ReplaceStorageNode.",
					ValidateFunc: IsUUID,
				},

				"address_setting": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					MinItems:    1,
					MaxItems:    6,
					Description: "IP settings to be assigned to storage nodes being added. Mandatory if export_file_type is AddStorageNodes.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"index": &schema.Schema{
								Type:         schema.TypeInt,
								Required:     true,
								Description:  "The ID of the node to be added. Must be 1 to 6.",
								ValidateFunc: validation.IntBetween(1, 6),
							},
							"control_port_ipv4_address": &schema.Schema{
								Type:         schema.TypeString,
								Required:     true,
								Description:  "IPv4 address of the control port.",
								ValidateFunc: validation.IsIPv4Address,
							},
							"internode_port_ipv4_address": &schema.Schema{
								Type:         schema.TypeString,
								Required:     true,
								Description:  "IPv4 address of the internode port.",
								ValidateFunc: validation.IsIPv4Address,
							},
							"compute_port_ipv4_address": &schema.Schema{
								Type:         schema.TypeString,
								Required:     true,
								Description:  "IPv4 address of the compute port.",
								ValidateFunc: validation.IsIPv4Address,
							},
							"compute_port_ipv6_address": &schema.Schema{
								Type:         schema.TypeString,
								Optional:     true,
								Description:  "IPv6 address of the compute port (Azure only).",
								ValidateFunc: validation.IsIPv6Address,
							},
						},
					},
				},
			},
		},
	},

	// output
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The status of the configuration file operation.",
	},

	"output_file_path": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full path where the downloaded file is saved. Appends .tar.gz if `download_path` is a filename without extension.",
	},
}
