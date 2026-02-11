package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var SpaceSchema = map[string]*schema.Schema{
	"partition_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Partition number within the parity group",
	},
	"ldev_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "LDEV ID associated with this partition",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status of the partition",
	},
	"lba_location": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "LBA (Logical Block Address) location",
	},
	"lba_size": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "LBA (Logical Block Address) size",
	},
}

var ParityGroupsInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of the storage system",
	},
	"parity_group_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Parity group ID",
	},
	"group_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Group type of the parity group",
	},
	"num_of_ldevs": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total number of LDEVs in parity group",
	},
	"used_capacity_rate": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Used capacity rate of the parity group",
	},
	"available_volume_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available volume capacity of the parity group",
	},
	"raid_level": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "RAID level of the parity group",
	},
	"raid_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "RAID type of the parity group",
	},
	"clpr_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "CLPR ID of the parity group",
	},
	"drive_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Drive type of the parity group",
	},
	"drive_type_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Drive type name of the parity group",
	},
	"is_copy_back_mode_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicates whether copy back mode is enabled for the parity group",
	},
	"is_encryption_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicates whether encryption is enabled for the parity group",
	},
	"total_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Total capacity of the parity group",
	},
	"physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Physical capacity of the parity group",
	},
	"available_physical_capacity": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available physical capacity of the parity group",
	},
	"is_accelerated_compression_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicates whether accelerated compression is enabled for the parity group",
	},
	"emulation_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Emulation type of the parity group",
	},
	"spaces": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Space/partition details within the parity group",
		Elem: &schema.Resource{
			Schema: SpaceSchema,
		},
	},
	"available_volume_capacity_in_kb": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Available volume capacity of the parity group in kb",
	},
}

var DataParityGroupSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"parity_group_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Parity group ID to retrieve",
	},
	// output
	"parity_group": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Parity group output",
		Elem: &schema.Resource{
			Schema: ParityGroupsInfoSchema,
		},
	},
}

var DataParityGroupsSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"drive_type_name": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Description: "Drive type.\n\n" +
			"For VSP One B20:\n" +
			"- SSD(QLC)\n" +
			"- SSD\n\n" +
			"For VSP 5000 series:\n" +
			"- SAS\n" +
			"- SSD(MLC)\n" +
			"- SSD(FMC)\n" +
			"- SSD\n" +
			"- SCM\n\n" +
			"For VSP E series:\n" +
			"- SAS\n" +
			"- SSD(MLC)\n" +
			"- SSD\n\n" +
			"For VSP G350, G370, G700, G900, VSP F350, F370, F700, F900:\n" +
			"- SAS\n" +
			"- SSD(MLC)\n" +
			"- SSD(FMC)\n" +
			"- SSD(RI)",
		ValidateFunc: validation.StringInSlice([]string{"SAS", "SCM", "SSD", "SSD(FMC)", "SSD(MLC)", "SSD(QLC)", "SSD(RI)"}, true),
	},
	"clpr_id": &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Description:  "CLPR number",
		ValidateFunc: validation.IntAtLeast(0),
	},
	"parity_group_ids": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "List of parity group IDs to retrieve",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"include_detail_info": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Include detailed information (FMC) for parity groups. When set to true, additional detailed fields will be populated.",
	},
	"include_cache_info": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Include cache information (class) for parity groups. When set to true, class-related fields will be populated.",
	},
	"parity_groups": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Parity groups output",
		Elem: &schema.Resource{
			Schema: ParityGroupsInfoSchema,
		},
	},
}
