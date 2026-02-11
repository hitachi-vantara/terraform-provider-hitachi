package terraform

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

// separate input and output schema for better readability, then combine them

// ------------------- IPv4 Information Schema -------------------
func ipv4InformationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv4 address.",
		},
		"subnet_mask": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv4 subnet mask.",
		},
		"default_gateway": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv4 default gateway.",
		},
	}
}

// ------------------- IPv6 Information Schema -------------------
func ipv6InformationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"linklocal": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv6 link-local setting.",
		},
		"linklocal_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv6 link-local address.",
		},
		"linklocal_address_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv6 link-local address status.",
		},
		"global": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv6 global setting.",
		},
		"global_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv6 global address.",
		},
		"global_address_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv6 global address status.",
		},
		"default_gateway": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv6 default gateway.",
		},
	}
}

// ------------------- FC Information Schema -------------------
func fcInformationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"al_pa": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Arbitrated Loop Physical Address.",
		},
		"fabric_switch_setting": {
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Description: "Fabric switch setting.",
		},
		"connection_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Connection type.",
		},
		"sfp_data_transfer_rate": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "SFP data transfer rate.",
		},
		"port_mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Port mode.",
		},
	}
}

// ------------------- iSCSI Information Schema -------------------
func iscsiInformationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vlan_use": {
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Description: "VLAN usage setting.",
		},
		"vlan_id": {
			Type:     schema.TypeInt,
			Computed: true,

			Description: "VLAN ID.",
		},
		"ip_mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP mode setting.",
		},
		"ipv4_information": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "IPv4 information.",
			Elem: &schema.Resource{
				Schema: ipv4InformationSchema(),
			},
		},
		"ipv6_information": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "IPv6 information.",
			Elem: &schema.Resource{
				Schema: ipv6InformationSchema(),
			},
		},
		"is_ipv6_updating": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "IPv6 updating status.",
		},
		"selective_ack": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Selective ACK setting.",
		},
		"delayed_ack": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Delayed ACK setting.",
		},
		"mtu_size": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "MTU size.",
		},
		"link_mtu_size": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Link MTU size.",
		},
		"virtual_port_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Virtual port enabled setting.",
		},
		"tcp_port": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "TCP port number.",
		},
		"window_size": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "TCP window size.",
		},
		"keep_alive_timer": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Keep alive timer setting.",
		},
		"isns_server_mode": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "iSNS server mode setting.",
		},
		"isns_server_ip_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "iSNS server IP address.",
		},
		"isns_server_port": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "iSNS server port number.",
		},
	}
}

// ------------------- NVMe-TCP Information Schema -------------------
func nvmeTcpInformationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vlan_use": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "VLAN usage setting.",
		},
		"vlan_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "VLAN ID.",
		},
		"ip_mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP mode setting.",
		},
		"ipv4_information": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "IPv4 information.",
			Elem: &schema.Resource{
				Schema: ipv4InformationSchema(),
			},
		},
		"ipv6_information": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "IPv6 information.",
			Elem: &schema.Resource{
				Schema: ipv6InformationSchema(),
			},
		},
		"is_ipv6_updating": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "IPv6 updating status.",
		},
		"selective_ack": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Selective ACK setting.",
		},
		"delayed_ack": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Delayed ACK setting.",
		},
		"mtu_size": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "MTU size.",
		},
		"link_mtu_size": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Link MTU size.",
		},
		"virtual_port_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Virtual port enabled setting.",
		},
		"tcp_port": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "TCP port number.",
		},
		"discovery_tcp_port": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Discovery TCP port number.",
		},
		"window_size": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "TCP window size.",
		},
	}
}

// ------------------- Ports Info Schema -------------------
func portsInfoListSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ports_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of ports.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Port ID.",
					},
					"protocol": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Port protocol.",
					},
					"port_wwn": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Port World Wide Name.",
					},
					"port_iscsi_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Port iSCSI name.",
					},
					"port_speed": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Port speed.",
					},
					"actual_port_speed": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Actual port speed.",
					},
					"port_security": {
						Type:        schema.TypeBool,
						Computed:    true,
						Optional:    true,
						Description: "Port security setting.",
					},
					"fc_information": {
						Type:        schema.TypeList,
						Computed:    true,
						Optional:    true,
						Description: "Fibre Channel information.",
						Elem: &schema.Resource{
							Schema: fcInformationSchema(),
						},
					},
					"iscsi_information": {
						Type:        schema.TypeList,
						Computed:    true,
						Optional:    true,
						Description: "iSCSI information.",
						Elem: &schema.Resource{
							Schema: iscsiInformationSchema(),
						},
					},
					"nvme_tcp_information": {
						Type:        schema.TypeList,
						Computed:    true,
						Optional:    true,
						Description: "NVMe-TCP information.",
						Elem: &schema.Resource{
							Schema: nvmeTcpInformationSchema(),
						},
					},
				},
			},
		},
		"port_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of ports returned.",
		},
	}
}

// ------------------- Port Info Schema -------------------
func portInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Port ID.",
		},
		"protocol": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port protocol.",
		},
		"port_wwn": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port World Wide Name.",
		},
		"port_iscsi_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port iSCSI name.",
		},
		"port_speed": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port speed.",
		},
		"actual_port_speed": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Actual port speed.",
		},
		"port_security": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Port security setting.",
		},
		"fc_information": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Fibre Channel information.",
			Elem: &schema.Resource{
				Schema: fcInformationSchema(),
			},
		},
		"iscsi_information": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "iSCSI information.",
			Elem: &schema.Resource{
				Schema: iscsiInformationSchema(),
			},
		},
		"nvme_tcp_information": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "NVMe-TCP information.",
			Elem: &schema.Resource{
				Schema: nvmeTcpInformationSchema(),
			},
		},
	}
}

// ------------------- Datasource Get Multiple Ports Schema -------------------
func datasourceAdminMultiplePortsInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of storage",
		},
		"protocol": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Filter ports by protocol. Valid values: FC, iSCSI, NVME_TCP.",
			ValidateFunc: validation.StringInSlice([]string{
				"FC",
				"iSCSI",
				"NVME_TCP",
			}, false),
		},
	}
}

func datasourceAdminMultiplePortsOutputSchema() map[string]*schema.Schema {
	return portsInfoListSchema()
}

func DatasourceAdminMultiplePortsSchema() map[string]*schema.Schema {
	schema := datasourceAdminMultiplePortsInputSchema()

	for k, v := range datasourceAdminMultiplePortsOutputSchema() {
		schema[k] = v
	}

	return schema
}

// ------------------- Datasource Get One Port Schema -------------------
func datasourceAdminOnePortInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of storage",
		},
		"port_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "ID of the port to retrieve.",
		},
	}
}

func datasourceAdminOnePortOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port_info": { // returns only one port inside a list
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Port Information",
			Elem: &schema.Resource{
				Schema: portInfoSchema(),
			},
		},
	}
}

func DatasourceAdminOnePortSchema() map[string]*schema.Schema {
	schema := datasourceAdminOnePortInputSchema()

	for k, v := range datasourceAdminOnePortOutputSchema() {
		schema[k] = v
	}

	return schema
}

// ------------------- Resource Port Update Schema -------------------

// Update-specific schemas for nested structures
func updateIPv4InformationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IPv4 address.",
		},
		"subnet_mask": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IPv4 subnet mask.",
		},
		"default_gateway": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IPv4 default gateway.",
		},
	}
}

func updateIPv6InformationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"linklocal": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IPv6 link-local setting.",
		},
		"global": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IPv6 global setting.",
		},
		"default_gateway": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IPv6 default gateway.",
		},
	}
}

func updateFCInformationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"al_pa": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Arbitrated Loop Physical Address. Hexidecimal notation 01-EF",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v := val.(string)

				if v == "" {
					return
				}

				re := regexp.MustCompile(`^[0-9A-F]{2}$`)
				if !re.MatchString(v) {
					errs = append(errs, fmt.Errorf("%q must be in hexadecimal format without 0x prefix (e.g. EF), got: %q", key, v))
				}
				return
			},
		},
		"fabric_switch_setting": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Fabric switch setting.",
		},
		"connection_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Connection type. Valid values: Point_To_Point, FC_AL.",
			ValidateFunc: validation.StringInSlice([]string{
				"Point_To_Point",
				"FC_AL",
			}, false),
		},
	}
}

func updateISCSIInformationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vlan_use": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "VLAN usage setting.",
		},
		"add_vlan_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "VLAN ID to add.",
			ValidateFunc: validation.IntBetween(1, 4094),
		},
		"delete_vlan_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "VLAN ID to delete.",
			ValidateFunc: validation.IntBetween(1, 4094),
		},
		"ip_mode": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IP mode setting. Valid values: ipv4, ipv4v6.",
			ValidateFunc: validation.StringInSlice([]string{
				"ipv4",
				"ipv4v6",
			}, false),
		},
		"ipv4_information": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "IPv4 information.",
			Elem: &schema.Resource{
				Schema: updateIPv4InformationSchema(),
			},
		},
		"ipv6_information": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "IPv6 information.",
			Elem: &schema.Resource{
				Schema: updateIPv6InformationSchema(),
			},
		},
		"tcp_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "TCP port number.",
		},
		"selective_ack": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Selective ACK setting.",
		},
		"delayed_ack": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Delayed ACK setting.",
		},
		"window_size": {
			Type:     schema.TypeString,
			Optional: true,

			Description: "TCP window size. Valid values: NUMBER_64K, NUMBER_128K, NUMBER_256K, NUMBER_512K, NUMBER_1024K.",
			ValidateFunc: validation.StringInSlice([]string{
				"NUMBER_64K",
				"NUMBER_128K",
				"NUMBER_256K",
				"NUMBER_512K",
				"NUMBER_1024K",
			}, false),
		},
		"mtu_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "MTU size. Valid values: NUMBER_1500, NUMBER_4500, NUMBER_9000.",
			ValidateFunc: validation.StringInSlice([]string{
				"NUMBER_1500",
				"NUMBER_4500",
				"NUMBER_9000",
			}, false),
		},
		"keep_alive_timer": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Keep alive timer setting.",
		},
		"isns_server_mode": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "iSNS server mode setting.",
		},
		"isns_server_ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "iSNS server IP address.",
		},
		"isns_server_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "iSNS server port number.",
		},
	}
}

func updateNVMeTCPInformationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vlan_use": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "VLAN usage setting.",
		},
		"add_vlan_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "VLAN ID to add.",
			ValidateFunc: validation.IntBetween(1, 4094),
		},
		"delete_vlan_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "VLAN ID to delete.",
			ValidateFunc: validation.IntBetween(1, 4094),
		},
		"ip_mode": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IP mode setting. Valid values: ipv4, ipv4v6.",
			ValidateFunc: validation.StringInSlice([]string{
				"ipv4",
				"ipv4v6",
			}, false),
		},
		"ipv4_information": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "IPv4 information.",
			Elem: &schema.Resource{
				Schema: updateIPv4InformationSchema(),
			},
		},
		"ipv6_information": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "IPv6 information.",
			Elem: &schema.Resource{
				Schema: updateIPv6InformationSchema(),
			},
		},
		"tcp_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "TCP port number.",
		},
		"discovery_tcp_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Discovery TCP port number.",
		},
		"selective_ack": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Selective ACK setting.",
		},
		"delayed_ack": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Delayed ACK setting.",
		},
		"window_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "TCP window size. Valid values: NUMBER_64K, NUMBER_128K, NUMBER_256K, NUMBER_512K, NUMBER_1024K, NUMBER_2048K.",
			ValidateFunc: validation.StringInSlice([]string{
				"NUMBER_64K",
				"NUMBER_128K",
				"NUMBER_256K",
				"NUMBER_512K",
				"NUMBER_1024K",
				"NUMBER_2048K",
			}, false),
		},
		"mtu_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "MTU size. Valid values: NUMBER_1500, NUMBER_4500, NUMBER_9000.",
			ValidateFunc: validation.StringInSlice([]string{
				"NUMBER_1500",
				"NUMBER_4500",
				"NUMBER_9000",
			}, false),
		},
	}
}

func resourceAdminPortInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of storage",
		},
		"port_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "ID of the port to update.",
		},
		"port_speed": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `Port speed setting. Valid values: NUMBER_0, NUMBER_1, NUMBER_4, NUMBER_8, NUMBER_10, NUMBER_16, NUMBER_25, NUMBER_32, NUMBER_64, NUMBER_100.
	- Equivalences: NUMBER_0 = Auto, NUMBER_1 = 1Gbps, NUMBER_4 = 4Gbps, NUMBER_8 = 8Gbps, NUMBER_10 = 10Gbps, NUMBER_16 = 16Gbps, NUMBER_25 = 25Gbps, NUMBER_32 = 32Gbps, NUMBER_64 = 64Gbps, NUMBER_100 = 100Gbps.
	- For iSCSI 10G (optical) ports, the speed is fixed at 10G, and any other specified value is ignored.`,
			ValidateFunc: validation.StringInSlice([]string{
				"NUMBER_0",
				"NUMBER_1",
				"NUMBER_4",
				"NUMBER_8",
				"NUMBER_10",
				"NUMBER_16",
				"NUMBER_25",
				"NUMBER_32",
				"NUMBER_64",
				"NUMBER_100",
			}, false),
		},
		"port_security": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Port security setting. This parameter cannot be specified for NVMe/TCP ports and will cause an error if provided.",
		},
		"fc_information": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "Fibre Channel information to update.",
			Elem: &schema.Resource{
				Schema: updateFCInformationSchema(),
			},
		},
		"iscsi_information": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "iSCSI information to update.",
			Elem: &schema.Resource{
				Schema: updateISCSIInformationSchema(),
			},
		},
		"nvme_tcp_information": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "NVMe-TCP information to update.",
			Elem: &schema.Resource{
				Schema: updateNVMeTCPInformationSchema(),
			},
		},
	}
}

func resourceAdminPortOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Updated port information.",
			Elem: &schema.Resource{
				Schema: portInfoSchema(),
			},
		},
	}
}

func ResourceAdminPortSchema() map[string]*schema.Schema {
	schema := resourceAdminPortInputSchema()

	for k, v := range resourceAdminPortOutputSchema() {
		schema[k] = v
	}

	return schema
}
