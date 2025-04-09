package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DataVssbPortSchema = map[string]*schema.Schema{
	"vosb_block_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of the REST API server on Virtual Storage Software block.",
	},
	"port": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Port Id of the port, or iqn of the port or port name of the port",
	},
	// output
	"vosb_port": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is port output",
		Elem: &schema.Resource{
			Schema: VssbPortInfoSchema,
		},
	},
}

var VssbPortInfoSchema = map[string]*schema.Schema{
	// "id": &schema.Schema{
	// 	Type:        schema.TypeString,
	// 	Computed:    true,
	// 	Description: "The ID of the compute port.",
	// },
	"protocol": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `x ∈ { "FC" , "iSCSI" }
		The protocol for connecting compute ports.`,
	},
	"type": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `x ∈ { "Target" , "Initiator" , "Universal" }
		The type of the compute port.`,
	},
	"nickname": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `The compute port nickname. Each compute port must have its own unique nickname.
		must match /^[a-zA-Z0-9!#\$%&'\+\-\.=@\^_\{\}~\(\)\[\]:]{1,32}$/`,
	},
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `The WWN of the allocation destination compute port of the target operation for FC connection, or the iSCSI name for iSCSI connections.
		The same name cannot be used for multiple compute ports.`,
	},
	"configured_port_speed": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `x ∈ { "16G" , "32G" , "Auto" }
		Link speed setting. The actual link speed is determined based on this setting.`,
	},
	"port_speed": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: ` x ∈ { "1G" , "2G" , "4G" , "8G" , "10G" , "16G" , "25G" , "32G" , "40G" , "Unknown" , "LinkDown", "DependsOnHypervisor" }
		Actual link speed (unit: bps).
		If configuredPortSpeed is Auto, a value is output as per the actual cable or switch specifications.
		Unknown: The status is unknown.
		LinkDown: Link down occurred.
		DependsOnHypervisor: Depends on the hypervisor setting.`,
	},
	"por_speed_duplex": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `x ∈ { "10Mbps Half" , "10Mbps Full" , "100Mbps Half" , "100Mbps Full" , "1Gbps Half" , "1Gbps Full" , "2.5Gbps Full" , "5Gbps Full" , "10Gbps Full" , "20Gbps Full" , "25Gbps Full" , "40Gbps Full" , "50Gbps Full" , "56Gbps Full" , "100Gbps Full" , "200Gbps Full" , "400Gbps Full" , "1G" , "8G" , "10G" , "16G" , "25G" , "32G" , "40G" , "Unknown" , "LinkDown" , "DependsOnHypervisor" }
		Actual link speed and duplex settings of the physical port used for communication. Only link speed is displayed for FC connection configuration.
		If configuredPortSpeed is Auto, a value is output as per the actual cable or switch specifications.
		Unknown: The status is unknown.
		LinkDown: Link down occurred.
		DependsOnHypervisor: Depends on the hypervisor setting.
		(Virtual machine) DependsOnHypervisor is always output for iSCSI connection configuration.`,
	},
	"protection_domain_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the protection domain to which the volume is belonging.",
	},
	"storage_node_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the storage node that has compute ports.",
	},
	"interface_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "An interface name, which is unique within a storage node that contains compute ports, control ports, and ports between storage nodes. Example: eth0, eth1",
	},
	"status_summary": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `x ∈ { "Normal", "Warning", "Error" }
		The summary of the compute port status.
			o Normal: No action by the user is required.
			o Warning: Although immediate action by the user is not required, some action may have to be taken.
			o Error: Immediate action by the user is required.`,
	},
	"status": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: `x ∈ { "Normal" , "Error" , "MaintenanceBlockage" }
		The status of the compute port.
			o Normal: Available.
			o Error: Unavailable.
			o MaintenanceBlockage: Unavailable (during maintenance blockade).`,
	},
	"fc_information": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Description: "fcTarget: object nullable",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"connection_type": {
					Computed: true,
					Type:     schema.TypeString,
					Description: ` x ∈ { "PointToPoint" }
					Network connection type.`,
				},
				"sfp_data_transfer_rate": {
					Computed: true,
					Type:     schema.TypeString,
					Description: ` x ∈ { "8G" , "16G" , "32G" , "Unknown" }
					SFP data transfer rate (unit: bps). "Unknown" is output if the SFP extension port cannot be recognized or nothing is connected to the SFP extension port.`,
				},
				"physical_wwn": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Physical WWN of the compute port if 'name' of the compute port is a logical WWN.",
				},
			},
		},
	},
	"iscsi_information": &schema.Schema{
		Computed:    true,
		Type:        schema.TypeList,
		Description: "iscsiUniversal: object nullable",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_mode": {
					Computed: true,
					Type:     schema.TypeString,
					Description: `x ∈ { "ipv4" , "ipv4v6" }
					Enables or disables IPv4/IPv6.
						o ipv4: Enables IPv4 only.
						o ipv4v6: Enables both IPv4 and IPv6.`,
				},
				"delayed_ack": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Whether TCP delayed ACKs are used. When 'true' is specified, TCP delayed ACKs are used.",
				},
				"mtu_size": {
					Computed: true,
					Type:     schema.TypeInt,
					Description: `{ x ∈ ℤ | 1500 ≤ x ≤ 9000 }
					The MTU size of Ethernet (unit: byte).`,
				},
				"mac_address": {
					Computed: true,
					Type:     schema.TypeString,
					Description: `(17 chars) , must match /^([a-f0-9]{2}:){5}[a-f0-9]{2}$/
					MAC address.`,
				},
				"is_isns_client_enabled": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "iSNS client function. Specifying 'true' enables the iSNS client function.",
				},
				"ipv4_information": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeList,
					Description: "ipv4InformationOfUniversal: object",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"address": {
								Computed: true,
								Type:     schema.TypeString,
								Description: `(7 to 15 chars)
								IP address (IPv4).`,
							},
							"subnet_mask": {
								Computed: true,
								Type:     schema.TypeString,
								Description: `(7 to 15 chars)
								Subnet mask (IPv4).`,
							},
							"default_gateway": {
								Computed: true,
								Type:     schema.TypeString,
								Description: `(up to 15 chars)
								The IP address of the default gateway (IPv4).`,
							},
						},
					},
				},
				"ipv6_information": &schema.Schema{
					Computed:    true,
					Type:        schema.TypeList,
					Description: "ipv6InformationOfUniversal: object",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"linklocal_address_mode": {
								Computed: true,
								Type:     schema.TypeString,
								Description: `x ∈ { "Auto" , "Manual" }
								Setting mode for link local addresses.`,
							},
							"linklocal_address": {
								Computed: true,
								Type:     schema.TypeString,
								Description: `(up to 39 chars)
								IPv6 link local address.
								An empty string "" is output if no address is set.`,
							},
							"global_address_mode": {
								Computed: true,
								Type:     schema.TypeString,
								Description: `x ∈ { "Auto" , "Manual" }
								Setting mode for IPv6 global addresses.`,
							},
							"global_address_1": {
								Computed: true,
								Type:     schema.TypeString,
								Description: `(up to 39 chars)
								IPv6 global address 1.
								An empty string "" is output if no address is set.`,
							},
							"subnet_prefix_length_1": {
								Computed: true,
								Type:     schema.TypeInt,
								Description: `{ x ∈ ℤ | 0 ≤ x ≤ 128 }
								Subnet prefix length of IPv6 global address 1.`,
							},
							"default_gateway": {
								Computed: true,
								Type:     schema.TypeString,
								Description: `(up to 39 chars)
								The IP address of the default gateway (IPv6).
								An empty string "" is output if no address is set.`,
							},
						},
					},
				},
				"isns_servers": &schema.Schema{
					Computed: true,
					Type:     schema.TypeList,
					Description: `iSNS server as the connection destination in the iSNS client function.
								ITEMS
								o isnsServerOfIscsiUniversal: object`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"index": {
								Computed: true,
								Type:     schema.TypeInt,
								Description: ` { x ∈ ℤ | 1 ≤ x ≤ 1 }
								The ID of the iSNS server.`,
							},
							"server_name": {
								Computed: true,
								Type:     schema.TypeString,
								Description: `(up to 45 chars) , must match /^$|^(([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^((([0-9a-f]{1,4}:){7}([0-9a-f]{1,4}))|(([0-9a-f]{1,4}:){6}:)|(([0-9a-f]{1,4}:){5}(((:[0-9a-f]{1,4}))|:))|(([0-9a-f]{1,4}:){4}(((:[0-9a-f]{1,4}){1,2})|:))|(([0-9a-f]{1,4}:){3}(((:[0-9a-f]{1,4}){1,3})|:))|(([0-9a-f]{1,4}:){2}(((:[0-9a-f]{1,4}){1,4})|:))|(([0-9a-f]{1,4}:){1}(((:[0-9a-f]{1,4}){1,5})|:))|(:(((:[0-9a-f]{1,4}){1,6})|:)))$/
								IP address (IPv4 or IPv6) setting of the iSNS server. An empty string "" is output if no address is set.`,
							},
							"port": {
								Computed: true,
								Type:     schema.TypeInt,
								Description: `{ x ∈ ℤ | 1 ≤ x ≤ 65536 }
								TCP port number of the iSNS server.`,
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
		Description: "Information about the authentication settings for the compute port for the target operation.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auth_mode": &schema.Schema{
					Computed: true,
					Type:     schema.TypeString,
					Description: ` x ∈ { "CHAP" , "CHAPComplyingWithInitiatorSetting" , "None" }
					Authentication scheme of the compute port.
					o CHAP: CHAP authentication.
					o CHAPComplyingWithInitiatorSetting: Complies with the setting of the compute node. If the setting is "CHAP", CHAP authentication is performed. If the setting is "None", no authentication is required.
					o None: No authentication is performed.`,
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
					Description: "The ID of the CHAP user.",
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
