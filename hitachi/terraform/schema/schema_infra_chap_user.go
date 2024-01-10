package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var InfraIscsiChapUserInfoSchema = map[string]*schema.Schema{
	"storage_serial_number": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Serial number of storage",
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
		Description: `Type of CHAP user name
		o target : CHAP user name of the iSCSI target side
		o initiator : CHAP user name of the host bus adapter (iSCSI initiator) side`,
	},
	"chap_user_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "CHAP user name.",
	},
	"chap_user_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Object ID for the CHAP user",
	},
}

var DataInfraIscsiChapUserSchema = map[string]*schema.Schema{

	"storage_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Unique ID of the storage device",
	},

	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Serial number of storage",
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
		Description: `Type of CHAP user name
		o target : CHAP user name of the iSCSI target side
		o initiator : CHAP user name of the host bus adapter (iSCSI initiator) side`,
	},
	"chap_user_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "CHAP user name.",
	},

	// output
	"chap_user": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "This is chap user output",
		Elem: &schema.Resource{
			Schema: InfraIscsiChapUserInfoSchema,
		},
	},
}

var DataInfraIscsiChapUsersSchema = map[string]*schema.Schema{

	"storage_id": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Unique ID of the storage device",
	},

	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Serial number of storage",
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
		Optional:    true,
		Description: "This is chap users output",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		/*
			Elem: &schema.Resource{
				Schema: IscsiChapUserInfoSchema,
			},
		*/
	},
}

var ResourceInfraIscsiChapUserSchema = map[string]*schema.Schema{
	"serial": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of storage",
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
		Description: `Type of CHAP user name
			o target : CHAP user name of the iSCSI target side
			o initiator : CHAP user name of the host bus adapter (iSCSI initiator) side`,
	},
	"chap_user_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "CHAP user name.",
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
		Optional:    true,
		Description: "This is chap user output",
		Elem: &schema.Resource{
			Schema: InfraIscsiChapUserInfoSchema,
		},
	},
}
