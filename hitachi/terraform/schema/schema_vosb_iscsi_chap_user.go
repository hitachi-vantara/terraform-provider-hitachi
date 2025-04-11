package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var VssbIscsiChapUserInfoSchema = map[string]*schema.Schema{
	"chap_user_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the CHAP user.",
	},
	"target_chap_user_name": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `CHAP user name used for CHAP authentication on the compute port (i.e., target side).
		(1 to 223 chars) , must match /^[a-zA-Z0-9\.:@_\-\+=\[\]~ ]{1,223}$/`,
	},
	"initiator_chap_user_name": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `CHAP user name used for CHAP authentication on the initiator port of the compute node in mutual CHAP authentication.
		(1 to 223 chars) , must match /^[a-zA-Z0-9\.:@_\-\+=\[\]~ ]{1,223}$/`,
	},
}

var DataVssbIscsiChapUsersSchema = map[string]*schema.Schema{

	"vosb_block_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the REST API server on Virtual Storage Software block.",
	},
	"target_chap_user": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "CHAP user name or CHAP user ID used for CHAP authentication on the compute port (i.e., target side).",
	},
	// output
	"chap_users": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: VssbIscsiChapUserInfoSchema,
		},
	},
}

var DataVssbIscsiChapUserSchema = map[string]*schema.Schema{

	"vosb_block_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the REST API server on Virtual Storage Software block.",
	},
	"target_chap_user_name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: `CHAP user name used for CHAP authentication on the compute port (i.e., target side).
		(1 to 223 chars) , must match /^[a-zA-Z0-9\.:@_\-\+=\[\]~ ]{1,223}$/`,
	},
	// output
	"chap_user": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: VssbIscsiChapUserInfoSchema,
		},
	},
}

var ResourceVssbChapUserSchema = map[string]*schema.Schema{
	"vosb_block_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the REST API server on Virtual Storage Software block.",
	},
	"chap_user_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The ID of the CHAP user.",
	},
	"target_chap_user_name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: `CHAP user name used for CHAP authentication on the compute port (i.e., target side).
		(1 to 223 chars) , must match /^[a-zA-Z0-9\.:@_\-\+=\[\]~ ]{1,223}$/`,
	},
	"target_chap_user_secret": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: `CHAP secret used for CHAP authentication on the compute port (i.e., target side).
		(12 to 32 chars) , must match /^[a-zA-Z0-9\.:@_\-\+=\/\[\]~ ]{12,32}$/`,
	},
	"initiator_chap_user_name": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Description: `CHAP user name used for CHAP authentication on the initiator port of the compute node in mutual CHAP authentication.
		(1 to 223 chars) , must match /^[a-zA-Z0-9\.:@_\-\+=\[\]~ ]{1,223}$/`,
	},
	"initiator_chap_user_secret": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Description: `CHAP secret used for CHAP authentication on the initiator port of the compute node in mutual CHAP authentication.
		(12 to 32 chars) , must match /^[a-zA-Z0-9\.:@_\-\+=\/\[\]~ ]{12,32}$/`,
	},

	// output
	"chap_users": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is chap users output",
		Elem: &schema.Resource{
			Schema: VssbIscsiChapUserInfoSchema,
		},
	},
}
