package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var IscsiChapUserInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The serial number of the storage",
	},
	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Target ID of the iSCSI target.",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port number",
	},
	"chap_user_type": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `Type of the CHAP user name
		o target : The CHAP user name of the iSCSI target side
		o initiator : The CHAP user name of the host bus adapter (iSCSI initiator) side`,
	},
	"chap_user_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The CHAP user name.",
	},
	"chap_user_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Object ID for the CHAP user",
	},
}

var DataIscsiChapUserSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The serial number of the storage",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port number",
	},
	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Target ID of the iSCSI target.",
	},
	"chap_user_type": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: `Type of the CHAP user name
		o target : The CHAP user name of the iSCSI target side
		o initiator : The CHAP user name of the host bus adapter (iSCSI initiator) side`,
	},
	"chap_user_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The CHAP user name.",
	},

	// output
	"chap_user": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: IscsiChapUserInfoSchema,
		},
	},
}

var DataIscsiChapUsersSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The serial number of the storage",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port number",
	},
	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Target ID of the iSCSI target.",
	},
	// output
	"chap_users": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: IscsiChapUserInfoSchema,
		},
	},
}

var ResourceIscsiChapUserSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The serial number of the storage",
	},
	"port_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Port number",
	},
	"iscsi_target_number": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Target ID of the iSCSI target.",
	},
	"chap_user_type": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: `Type of the CHAP user name
			o target : The CHAP user name of the iSCSI target side
			o initiator : The CHAP user name of the host bus adapter (iSCSI initiator) side`,
	},
	"chap_user_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The CHAP user name.",
	},
	"chap_user_password": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Description: `Specify a secret consisting of 12 to 32 characters for the specified CHAP user.
			If you specify a null character, the password is reset.`,
	},

	// output
	"chap_user": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "This is output schema",
		Elem: &schema.Resource{
			Schema: IscsiChapUserInfoSchema,
		},
	},
}
