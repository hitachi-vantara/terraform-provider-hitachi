package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DataVssbStoragePortSchema = map[string]*schema.Schema{
	"vss_block_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the REST API server on Virtual Storage Software block.",
	},
	"port_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Port name of the storage device",
		Default:     "",
	},
	// output
	"ports": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is ports output",
		Elem: &schema.Resource{
			Schema: VssbStoragePortInfoSchema,
		},
	},
}

var VssbStoragePortInfoSchema = map[string]*schema.Schema{
	// "id": &schema.Schema{
	// 	Type:        schema.TypeString,
	// 	Computed:    true,
	// 	Description: "Id of port",
	// },
	"protocol": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Protocol of port",
	},
	"type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Type of port",
	},
	"nickname": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Nickname of port",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of port",
	},
	"configured_port_speed": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Configured port speed of port",
	},
	"port_speed": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port speed of port",
	},
	"por_speed_duplex": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Por speed duplex of port",
	},
	"protection_domain_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Protection domain ID of port",
	},
	"storage_node_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage node ID of port",
	},
	"interface_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Interface name of port",
	},
	"status_summary": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "status summary of port",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status of port",
	},
	"fc_information": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "FC information of port",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"connection_type": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Connection type of FC port",
				},
				"sfp_data_transfer_rate": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Data transfer rate of FC port",
				},
				"physical_wwn": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "WWN of port",
				},
			},
		},
	},
	"iscsi_information": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Iscsi information of port",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_mode": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "IP mode of iSCSI port",
				},
				"delayed_ack": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Delayed ack of iSCSI port",
				},
				"mtu_size": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "MTU size of iSCSI port",
				},
				"mac_address": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "MAC address of iSCSI port",
				},
				"is_isns_client_enabled": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Checks if is isns client enabled of iSCSI port",
				},
				"ipv4_information": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Ipv4 information",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"address": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Address of ipv4 information",
							},
							"subnet_mask": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Subnet mask of ipv4 information",
							},
							"default_gateway": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Default gateway of ipv4 information",
							},
						},
					},
				},
				"ipv6_information": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Ipv6 information",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"linklocal_address_mode": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Linklocal address mode of ipv6 information",
							},
							"linklocal_address": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Linklocal address of ipv6 information",
							},
							"global_address_mode": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Global address mode of ipv6 information",
							},
							"global_address_1": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Global address of ipv6 information",
							},
							"subnet_prefix_length_1": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Subnet prefix length of ipv6 information",
							},
							"default_gateway": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Default gateway of ipv6 information",
							},
						},
					},
				},
				"isns_servers": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Isns servers information",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"index": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Index of isns server",
							},
							"server_name": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Server name of isns server",
							},
							"port": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Port of isns server",
							},
						},
					},
				},
			},
		},
	},
	"port_auth_settings": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Port auth settings information",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auth_mode": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Auth mode of Port",
				},
				"is_discovery_chap_auth": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Is discovery chap auth of Port",
				},
				"is_mutual_chap_auth": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Is mutual chap of Port",
				},
			},
		},
	},
}
