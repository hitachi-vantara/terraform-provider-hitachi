package terraform

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var SupportedHostModesHostModeSchema = map[string]*schema.Schema{
	"host_mode_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Host mode number",
	},
	"host_mode_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Identification name of the host mode",
	},
	"host_mode_display": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Host mode value (value to be used to specify the host mode)",
	},
}

var SupportedHostModesHostModeOptionSchema = map[string]*schema.Schema{
	"host_mode_option_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Host mode option number",
	},
	"host_mode_option_description": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Description of the host mode option",
	},
	"scope": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Scope of the host mode option (for example, HostGroup)",
	},
	"required_host_modes": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of host mode IDs that are required when using this host mode option.",
		Elem:        &schema.Schema{Type: schema.TypeInt},
	},
}

var DataSupportedHostModesSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system",
	},
	"host_modes": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The following attributes related to the host mode are output.",
		Elem:        &schema.Resource{Schema: SupportedHostModesHostModeSchema},
	},
	"host_mode_options": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The following attributes related to the host mode option are output.",
		Elem:        &schema.Resource{Schema: SupportedHostModesHostModeOptionSchema},
	},
}
