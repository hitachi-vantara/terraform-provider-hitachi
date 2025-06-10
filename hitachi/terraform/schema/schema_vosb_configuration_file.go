package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ResourceVssbConfigurationFileSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The VOSB address.",
	},
	"download_existconfig_only": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If set to true, the provider skips the creation step and only downloads the latest existing configuration file. When true, all other parameters are ignored.",
	},
	"create_only": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If set to true, the provider creates a new configuration file but does not download it afterward. This flag is ignored if download_existconfig_only is true.",
	},
	"download_path": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the local file path to save the configuration file. Ignored if no download occurs. It could be a directory or file path",
	}, 
	// output
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The status of the storage node operation.",
	},
	"output_file_path": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full path to where the downloaded file will end up. Appends .tar.gz if download_path is a filename without extension",
	},
}
