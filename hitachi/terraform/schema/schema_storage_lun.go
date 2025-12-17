package terraform

import (
	// "regexp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var LunInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of the storage system",
	},

	"ldev_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "LDEV ID",
	},
	// vnext
	// "ldev_hex": {
	// 	Type:        schema.TypeString,
	// 	Computed:    true,
	// 	Description: "LDEV ID in hexadecimal format",
	// },
	"virtual_ldev_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Virtual LDEV ID assigned by the array",
	},
	"clpr_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "CLPR (Cache Logical Partition) ID the volume belongs to",
	},
	"emulation_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Emulation type of the LDEV (e.g., OPEN-V)",
	},

	"byte_format_capacity": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Human-readable capacity string (e.g., 2G, 512M)",
	},
	"block_capacity": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Volume capacity expressed in blocks (1 block = 512 bytes)",
	},

	// --- PORTS ---
	"num_ports": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of internal ports associated with this volume",
	},
	"ports": {
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"port_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Port ID (e.g., CL1-A)",
				},
				"hostgroup_number": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Internal host group number on the port",
				},
				"hostgroup_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the host group",
				},
				"lun": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "LUN number assigned to the host group",
				},
			},
		},
	},

	// --- ATTRIBUTES ---
	"attributes": {
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "List of volume attributes and flags",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},

	"label": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User-defined volume label",
	},
	"status": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Current status of the volume (e.g., NORMAL)",
	},

	"mpblade_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "MP Blade ID currently serving the volume",
	},
	"ssid": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Subsystem ID",
	},
	"pool_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Pool ID that the volume belongs to",
	},

	"num_of_used_block": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of used blocks",
	},
	"is_full_allocation_enabled": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicates whether full allocation is enabled",
	},
	"resource_group_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Resource group ID",
	},

	// --- DATA REDUCTION ---
	"data_reduction_status": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Data reduction status (ENABLED / DISABLED / PROCESSING)",
	},
	"data_reduction_mode": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Data reduction mode (compression, compression_deduplication)",
	},
	"data_reduction_process_mode": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Process mode for capacity saving (inline or post_process)",
	},
	"data_reduction_progress_rate": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Progress rate of data reduction operations (%)",
	},

	// --- ALUA ---
	"is_alua_enabled": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicates whether ALUA is enabled",
	},

	// --- NAA ID ---
	"naa_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "NAA ID (global identifier for the volume)",
	},

	// --- COMPRESSION ACCELERATION ---
	"is_compression_acceleration_enabled": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicates whether compression accelerator is enabled",
	},
	"compression_acceleration_status": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status of the compression accelerator",
	},

	// --- RAID INFORMATION ---
	"raid_level": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "RAID level of the parity group",
	},
	"raid_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "RAID type (internal/external/virtual)",
	},
	"num_of_parity_groups": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of parity groups the volume is composed of",
	},
	"parity_group_ids": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of parity group IDs",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"drive_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Drive type backing the volume (e.g., SAS, SSD)",
	},
	"drive_byte_format_capacity": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Drive capacity in human-readable units",
	},
	"drive_block_capacity": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Drive capacity in blocks",
	},

	// --- COMPOSING / SNAPSHOT / EXTERNAL ---
	"composing_pool_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Composing pool ID for tiered storage or HDT",
	},
	"snapshot_pool_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Snapshot pool ID used for copy-on-write snapshots",
	},

	"external_vendor_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Vendor ID of the external volume",
	},
	"external_product_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Product ID of the external volume",
	},
	"external_volume_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Internal numeric ID of the external volume",
	},
	"external_volume_id_string": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "External volume ID in string format",
	},

	// --- EXTERNAL PORTS ---
	"num_of_external_ports": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of external ports",
	},
	"external_ports": {
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Description: "List of external ports connected to the volume",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"port_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "External port ID",
				},
				"hostgroup_number": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "External host group number",
				},
				"lun": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "External LUN number",
				},
				"wwn": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "World Wide Name of the port",
				},
			},
		},
	},

	// --- QUORUM ---
	"quorum_disk_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Quorum disk ID (for GAD / HA clusters)",
	},
	"quorum_storage_serial_number": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Serial number of the quorum storage system",
	},
	"quorum_storage_type_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Type ID of the quorum storage system",
	},

	// --- NVME ---
	"namespace_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "NVMe namespace ID",
	},
	"nvm_subsystem_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "NVMe subsystem ID",
	},

	// --- HDT / RELOCATION / TIERING ---
	"is_relocation_enabled": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicates whether HDT relocation is enabled",
	},
	"tier_level": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Current tier level of the volume (0–3)",
	},

	"used_capacity_per_tier_level1": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity on Tier 1",
	},
	"used_capacity_per_tier_level2": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity on Tier 2",
	},
	"used_capacity_per_tier_level3": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity on Tier 3",
	},

	"tier_level_for_new_page_allocation": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Tier level where new pages are allocated",
	},

	// --- OPERATION ---
	"operation_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Current volume operation (e.g., formatting, relocation)",
	},
	"preparing_operation_progress_rate": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Percentage progress of the current operation",
	},

	// --- CAPACITIES ---
	"total_capacity_in_mb": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total volume capacity (MB)",
	},
	"free_capacity_in_mb": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Free capacity (MB)",
	},
	"used_capacity_in_mb": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity (MB)",
	},
}

var DataLunSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"ldev_id": {
		Type:          schema.TypeInt,
		Optional:      true,
		Description:   "LDEV ID.",
		ValidateFunc: validation.IntBetween(0, 65535),
	},

	// vnext
	// "ldev_id": {
	// 	Type:          schema.TypeInt,
	// 	Optional:      true,
	// 	Description:   "LDEV ID. Only one of ldev_id or ldev_hex may be specified, not both.",
	// 	ConflictsWith: []string{"ldev_hex"},
	// 	ValidateFunc: validation.IntBetween(0, 65535),
	// },
	// "ldev_hex": {
	// 	Type:          schema.TypeString,
	// 	Optional:      true,
	// 	Description:   "LDEV ID in hexadecimal format. Only one of ldev_id or ldev_hex may be specified, not both.",
	// 	ConflictsWith: []string{"ldev_id"},
	// 	// Validation: hex string, 1–4 hex chars (0x prefix optional)
	// 	ValidateFunc: validation.StringMatch(
	// 		regexp.MustCompile(`^(0x)?[A-Fa-f0-9]{1,4}$`),
	// 		"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
	// 	),
	// },
	// output
	"volume": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Volume output",
		Elem: &schema.Resource{
			Schema: LunInfoSchema,
		},
	},
}

var DataLunsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"start_ldev_id": &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 65535),
		Description:  "Start LDEV ID.",
	},
	"end_ldev_id": &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 65535),
		Description:  "End LDEV ID",
	},

	// vnext
	// "start_ldev_id": &schema.Schema{
	// 	Type:         schema.TypeInt,
	// 	Optional:     true,
	// 	ValidateFunc: validation.IntBetween(0, 65535),
	// 	Description:  "Start LDEV ID. Only one of start_ldev_id or start_ldev_hex may be specified, not both",
	// 	ConflictsWith: []string{"start_ldev_hex"},
	// },
	// "start_ldev_hex": {
	// 	Type:          schema.TypeString,
	// 	Optional:      true,
	// 	Description:   "Start LDEV ID in hexadecimal format. Only one of start_ldev_id or start_ldev_hex may be specified, not both.",
	// 	ConflictsWith: []string{"start_ldev_id"},
	// 	// Validation: hex string, 1–4 hex chars (0x prefix optional)
	// 	ValidateFunc: validation.StringMatch(
	// 		regexp.MustCompile(`^(0x)?[A-Fa-f0-9]{1,4}$`),
	// 		"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
	// 	),
	// },
	// "end_ldev_id": &schema.Schema{
	// 	Type:         schema.TypeInt,
	// 	Optional:     true,
	// 	ValidateFunc: validation.IntBetween(0, 65535),
	// 	Description:  "End LDEV ID. Only one of end_ldev_id or end_ldev_hex may be specified, not both.",
	// 	ConflictsWith: []string{"end_ldev_hex"},
	// },
	// "end_ldev_hex": {
	// 	Type:          schema.TypeString,
	// 	Optional:      true,
	// 	Description:   "End LDEV ID in hexadecimal format. Only one of end_ldev_id or end_ldev_hex may be specified, not both.",
	// 	ConflictsWith: []string{"end_ldev_id"},
	// 	// Validation: hex string, 1–4 hex chars (0x prefix optional)
	// 	ValidateFunc: validation.StringMatch(
	// 		regexp.MustCompile(`^(0x)?[A-Fa-f0-9]{1,4}$`),
	// 		"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
	// 	),
	// },

	"undefined_ldev": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If set to true, returns the LUNs that are not allocated",
	},
	// output
	"volumes": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Volumes output",
		Elem: &schema.Resource{
			Schema: LunInfoSchema,
		},
	},
}

var ResourceLunSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},

	"ldev_id": {
		Type:          schema.TypeInt,
		Optional:      true,
		Description:   "LDEV ID.",
		ValidateFunc: validation.IntBetween(0, 65535),
	},

	// vnext
	// // --- LDEV ID / LDEV HEX (mutually exclusive) ---
	// "ldev_id": {
	// 	Type:          schema.TypeInt,
	// 	Optional:      true,
	// 	Description:   "LDEV ID. Only one of ldev_id or ldev_hex may be specified, not both.",
	// 	ConflictsWith: []string{"ldev_hex"},
	// 	ValidateFunc: validation.IntBetween(0, 65535),
	// },
	// "ldev_hex": {
	// 	Type:          schema.TypeString,
	// 	Optional:      true,
	// 	Description:   "LDEV ID in hexadecimal format. Only one of ldev_id or ldev_hex may be specified, not both.",
	// 	ConflictsWith: []string{"ldev_id"},
	// 	// Validation: hex string, 1–4 hex chars (0x prefix optional)
	// 	ValidateFunc: validation.StringMatch(
	// 		regexp.MustCompile(`^(0x)?[A-Fa-f0-9]{1,4}$`),
	// 		"must be a valid hexadecimal LDEV value between 0x0 and 0xFFFF",
	// 	),
	// },

	// --- POOL / PARITY GROUP SELECTION (exactly one) ---
	"pool_id": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      -1,
		Description:  "Pool ID. One of pool_id, pool_name, paritygroup_id, external_paritygroup_id must be set.",
		ExactlyOneOf: []string{"pool_id", "pool_name", "paritygroup_id", "external_paritygroup_id"},
	},
	"pool_name": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "Pool name. One of pool_id, pool_name, paritygroup_id, external_paritygroup_id must be set.",
		ExactlyOneOf: []string{"pool_id", "pool_name", "paritygroup_id", "external_paritygroup_id"},
	},
	"paritygroup_id": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "Parity group ID. One of pool_id, pool_name, paritygroup_id, external_paritygroup_id must be set.",
		ExactlyOneOf: []string{"pool_id", "pool_name", "paritygroup_id", "external_paritygroup_id"},
	},
	"external_paritygroup_id": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "External parity group ID. One of pool_id, pool_name, paritygroup_id, external_paritygroup_id must be set.",
		ExactlyOneOf: []string{"pool_id", "pool_name", "paritygroup_id", "external_paritygroup_id"},
	},

	// --- OTHER FIELDS ---
	"size_gb": {
		Type:         schema.TypeFloat,
		Required:     true,
		Description:  "Size of the volume (GB) (supports decimal values like 1.5). Max 262,144 GB (256 TiB).",
		ValidateFunc: validation.FloatBetween(1, 262144),
	},
	"name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Name of the volume",
	},
	"capacity_saving": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "disabled",
		Description: "Capacity saving mode: compression_deduplication, compression, or disabled.",
		ValidateFunc: validation.StringInSlice([]string{
			"compression_deduplication",
			"compression",
			"disabled",
		}, true),
	},
	"is_data_reduction_shared_volume_enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Create a data reduction shared volume (TI Advanced). Must specify pool_id or pool_name and capacity_saving != disabled if true. Optional on create; ignored on update.",
	},
	"is_compression_acceleration_enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Enable the compression accelerator. If omitted and capacity saving is enabled, accelerator will be auto-enabled when available.",
	},
	"is_alua_enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Enable ALUA. Optional on update; error on create.",
	},
	"data_reduction_process_mode": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Capacity-saving mode. Allowed values: inline, post_process. Can only be used when capacity saving is enabled. Optional on update; error on create.",
		ValidateFunc: validation.StringInSlice([]string{
			"inline",
			"post_process",
		}, true),
	},

	// OUTPUT
	"volume": {
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Volume output",
		Elem: &schema.Resource{
			Schema: LunInfoSchema,
		},
	},
}
