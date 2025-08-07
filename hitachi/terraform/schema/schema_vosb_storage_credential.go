package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Schema for UserGroup
var VssbStorageUserGroupSchema = map[string]*schema.Schema{
	"user_group_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The ID of the user group.",
	},
	"user_group_object_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The object ID of the user group.",
	},
}

// Schema for Privilege
var VssbStorageUserPrivilegeSchema = map[string]*schema.Schema{
	"scope": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The scope of the privilege.",
	},
	"role_names": &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "List of role names associated with the privilege.",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
}

// Schema for User
var VssbStorageUserInfoSchema = map[string]*schema.Schema{
	"user_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The user ID.",
	},
	"user_object_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The object ID of the user.",
	},
	"password_expiration_time": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The password expiration time for the user.",
		// You could also use time.Time and handle it programmatically
	},
	"is_enabled": &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Whether the user is enabled.",
	},
	"user_groups": &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "List of user groups associated with the user.",
		Elem:        &schema.Resource{Schema: VssbStorageUserGroupSchema},
	},
	"is_built_in": &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Indicates if the user is a built-in user.",
	},
	"authentication": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The authentication method for the user.",
	},
	"role_names": &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "List of role names assigned to the user.",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"is_enabled_console_login": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Whether the user is enabled for console login. Can be null.",
	},
	"vps_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The virtual private server ID associated with the user.",
	},
	"privileges": &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "List of privileges assigned to the user.",
		Elem:        &schema.Resource{Schema: VssbStorageUserPrivilegeSchema},
	},
}

var ResourceVssbChangeUserPasswordSchema = map[string]*schema.Schema{
	// input
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the REST API server on Virtual Storage Software block.",
	},
	"user_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The user ID for the password change.",
	},
	"current_password": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The current password of the user.",
		Sensitive:   true,
	},
	"new_password": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The new password for the user.",
		Sensitive:   true,
	},
	// Output field for User
	"user_info": &schema.Schema{
		Type:        schema.TypeList, // Changed from TypeMap to TypeList
		Computed:    true,
		Description: "This is the user output information",
		Elem: &schema.Resource{
			Schema: VssbStorageUserInfoSchema,
		},
	}, "status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The status of the password change operation.",
	},
}
