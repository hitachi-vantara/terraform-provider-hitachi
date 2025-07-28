package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ProviderSchema = map[string]*schema.Schema{

	"hitachi_vosb_provider": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Hitachi VSP One SDS Block (VOSB) is a storage software product that builds and sets up a virtual storage system from multiple general-purpose servers. The system offers a high-performance, high-capacity block storage service with high reliability.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"vosb_address": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Host name or the IP address (IPv4) of VSP One SDS Block.",
				},
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Username of the VSP One SDS Block",
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Password of the VSP One SDS Block",
				},
			},
		},
		DefaultFunc: schema.EnvDefaultFunc("HITACHI_VOSB_PROVIDER", nil),
	},
	"san_storage_system": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Hitachi VSP One series and Hitachi VSP 5000 series are enterprise storage solutions designed to provide reliable and scalable block storage for a variety of environments. Both systems focus on simplifying data storage management while ensuring high availability and data integrity.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"serial": &schema.Schema{
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Serial number storage",
				},
				"management_ip": &schema.Schema{ // svp_ip for VSP-5000 or controller0 for other models
					Type:        schema.TypeString,
					Required:    true,
					Description: "Management IP for VSP-5000 series",
				},
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Username for VSP server",
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Password for VSP server",
				},
			},
		},
		DefaultFunc: schema.EnvDefaultFunc("HITACHI_SAN_STORAGE_SYSTEMS", nil),
	},
}
