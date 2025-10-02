package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var VssbComputePortInfoSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the port",
	},
	"authentication_settings": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: `Authentication scheme of the compute port.

	- CHAP: CHAP authentication.
	- CHAPComplyingWithInitiatorSetting: Complies with the setting of the compute node. If the setting is "CHAP", CHAP authentication is performed. If the setting is "None", no authentication is required.
	- None: No authentication is performed.`,
	},
	"target_chap_users": &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,

		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "List of chap users to be attached to the compute port.",
	},
}

var ResourceComputePortSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block.",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the port",
	},
	"authentication_settings": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		Description: `Authentication scheme of the compute port.

	- CHAP: CHAP authentication.
	- CHAPComplyingWithInitiatorSetting: Complies with the setting of the compute node. If the setting is "CHAP", CHAP authentication is performed. If the setting is "None", no authentication is required.
	- None: No authentication is performed.`,
	},
	"target_chap_users": &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "List of chap users to be attached to the compute port.",
	},
	// output
	"compute_port": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Outputs information about the volume",
		Elem: &schema.Resource{
			Schema: VssbIscsiPortAuthInfoSchema,
		},
	},
}

var DataSourceVssbComputePortSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block.",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the port",
	},

	// output
	"compute_port": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional:    true,
		Description: "Outputs information about the compute port",
		Elem: &schema.Resource{
			Schema: VssbIscsiPortAuthInfoSchema,
		},
	},
}

var VssbIscsiPortAuthInfoSchema = map[string]*schema.Schema{
	"port_auth_settings": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Information about the authentication settings for the compute port for the target operation.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auth_mode": &schema.Schema{
					Computed: true,
					Type:     schema.TypeString,
					Description: `Authentication scheme of the compute port.

	- CHAP: CHAP authentication.
	- CHAPComplyingWithInitiatorSetting: Complies with the setting of the compute node. If the setting is "CHAP", CHAP authentication is performed. If the setting is "None", no authentication is required.
	- None: No authentication is performed.`,
				},
				"is_discovery_chap_auth": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Enables or disables CHAP authentication at the time of discovery in iSCSI connection. Enables CHAP authentication at the time of discovery when true is specified.",
				},
				"is_mutual_chap_auth": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Enables or disables mutual CHAP authentication. Enables mutual CHAP authentication when true is specified.",
				},
			},
		},
	},
	"chap_users": &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"chap_user_id": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeString,
					Description: "The unique ID of the CHAP user associated with CHAP authentication settings.",
				},
				"target_chap_user_name": &schema.Schema{
					Computed: true,
					Type:     schema.TypeString,
					Description: `CHAP user name used for CHAP authentication on the compute port (i.e., target side).
					(1 to 223 chars) , must match /^[a-zA-Z0-9\.:@_\-\+=\[\]~ ]{1,223}$/`,
				},
				"initiator_chap_user_name": &schema.Schema{
					Computed: true,
					Type:     schema.TypeString,
					Description: `CHAP user name used for CHAP authentication on the initiator port of the compute node in mutual CHAP authentication.
					(1 to 223 chars) , must match /^[a-zA-Z0-9\.:@_\-\+=\[\]~ ]{1,223}$/`,
				},
			},
		},
	},
}
