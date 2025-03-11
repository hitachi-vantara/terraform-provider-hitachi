package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ProviderSchema = map[string]*schema.Schema{

	"hitachi_vss_block_provider": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Hitachi Virtual Storage Software Block (VSS block) is a storage software product that builds and sets up a virtual storage system from multiple general-purpose servers. The system offers a high-performance, high-capacity block storage service with high reliability.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"vss_block_address": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Host name or the IP address (IPv4) of Virtual Storage Software block.",
				},
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Username of the Virtual Storage Software block",
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Password of the Virtual Storage Software block",
				},
			},
		},
		DefaultFunc: schema.EnvDefaultFunc("HITACHI_VSS_BLOCK_PROVIDER", nil),
	},
	"san_storage_system": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Hitachi VSP 5000 series reliably delivers more data faster than ever for open systems and mainframe applications. VSP 5000 series provides response times as low as 39 microseconds and can be configured with up to 69 PB of raw capacity, with scalability to handle up to 33 million IOPS. All VSP 5000 models are backed by the industry’s most comprehensive 100% data availability guarantee to ensure that your operations are always up and running.",
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
