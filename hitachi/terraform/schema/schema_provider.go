package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ProviderSchema = map[string]*schema.Schema{

	"hitachi_vosb_provider": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "VSP One SDS Block and Cloud combines VSP One SDS Block, which creates virtual storage systems from general-purpose servers, with VSP One SDS Cloud, which enables deployment on AWS, Google Cloud Platform (GCP), and Microsoft Azure.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"vosb_address": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Host name or the IP address (IPv4) of the VSP One SDS Block and Cloud system.",
				},
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Username of the VSP One SDS Block and Cloud system",
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Password of the VSP One SDS Block and Cloud system",
				},
			},
		},
		DefaultFunc: schema.EnvDefaultFunc("HITACHI_VOSB_PROVIDER", nil),
	},
	"san_storage_system": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "VSP One Block 20 series, VSP One Block 80 series, and VSP 5000 series are enterprise storage solutions designed to provide reliable and scalable block storage for a variety of environments. Both systems focus on simplifying data storage management while ensuring high availability and data integrity.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"serial": &schema.Schema{
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Serial number of the storage system",
				},
				"management_ip": &schema.Schema{ // svp_ip for VSP-5000 or controller0 for other models
					Type:        schema.TypeString,
					Required:    true,
					Description: "Management IP for the VSP 5000 series",
				},
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Username for the storage system",
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Password for the storage system",
				},
			},
		},
		DefaultFunc: schema.EnvDefaultFunc("HITACHI_SAN_STORAGE_SYSTEMS", nil),
	},
	"hitachi_vsp_one_provider": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "VSP One Block Administrator is a configuration management tool designed for VSP One Block storage systems, simplifying and streamlining storage management.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"serial": &schema.Schema{
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Serial number of the storage system",
				},
				"management_ip": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Management IP address for the VSP One Block Administrator",
				},
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "User name for the VSP One Block Administrator",
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Password for the VSP One Block Administrator",
				},
			},
		},
		DefaultFunc: schema.EnvDefaultFunc("HITACHI_VOSB_PROVIDER", nil),
	},
}
