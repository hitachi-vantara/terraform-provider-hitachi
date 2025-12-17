package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DataVssbStoragePortSchema = map[string]*schema.Schema{
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the VSP One SDS Block system.",
	},
	"port_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Port name of the storage system",
		Default:     "",
	},
	// output
	"ports": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "Ports output",
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
		Description: "Protocol of the port",
	},
	"type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Type of the port",
	},
	"nickname": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Nickname of the port",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the port",
	},
	"configured_port_speed": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Configured port speed of the port",
	},
	"port_speed": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port speed of the port",
	},
	"por_speed_duplex": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Port speed duplex of the port",
	},
	"protection_domain_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Protection domain ID of the port",
	},
	"storage_node_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Storage node ID of the port",
	},
	"interface_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Interface name of the port",
	},
	"status_summary": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "status summary of the port",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Status of the port",
	},
	"fc_information": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Fibre Channel information of the port",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"connection_type": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Connection type of the Fibre Channel port",
				},
				"sfp_data_transfer_rate": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Data transfer rate of the Fibre Channel port",
				},
				"physical_wwn": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "WWN of the port",
				},
			},
		},
	},
	"iscsi_information": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Optional:    true,
		Description: "iSCSI information of the port",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_mode": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "IP mode of the iSCSI port",
				},
				"delayed_ack": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Delayed ACK of the iSCSI port",
				},
				"mtu_size": {
					Computed:    true,
					Type:        schema.TypeInt,
					Description: "MTU size of the iSCSI port",
				},
				"mac_address": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "MAC address of the iSCSI port",
				},
				"is_isns_client_enabled": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Checks if is iSNS client enabled of the iSCSI port",
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
					Description: "IPv6 information",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"linklocal_address_mode": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Linklocal address mode of IPv6 information",
							},
							"linklocal_address": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Linklocal address of IPv6 information",
							},
							"global_address_mode": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Global address mode of IPv6 information",
							},
							"global_address_1": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Global address of IPv6 information",
							},
							"subnet_prefix_length_1": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Subnet prefix length of IPv6 information",
							},
							"default_gateway": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Default gateway of IPv6 information",
							},
						},
					},
				},
				"isns_servers": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeList,
					Optional:    true,
					Description: "iSNS server information",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"index": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Index of iSNS server",
							},
							"server_name": {
								Computed:    true,
								Type:        schema.TypeString,
								Description: "Server name of iSNSs server",
							},
							"port": {
								Computed:    true,
								Type:        schema.TypeInt,
								Description: "Port of iSNS server",
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
					Description: "Auth mode of the port",
				},
				"is_discovery_chap_auth": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Is discovery CHAP auth of the port",
				},
				"is_mutual_chap_auth": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Is mutual CHAP of the port",
				},
			},
		},
	},
}
