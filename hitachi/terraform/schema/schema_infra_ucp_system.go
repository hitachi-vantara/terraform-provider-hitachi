package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var InfraUcpSystemSchema = map[string]*schema.Schema{
	"resource_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource  ID",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of UCP System",
	},
	"resource_state": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Resource State",
	},
	"compute_devices": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "Compute Devices",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"resource_id": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Resource Id",
				},
				"bmc_address": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "BMC Address",
				},
				"bmc_firmware_version": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "BMC Firmware Version",
				},
				"bios_version": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "BIOS Version",
				},
				"resource_state": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Resource State",
				},
				"model": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Model",
				},
				"serial": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Serial Number",
				},
				"is_management": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Is Management",
				},
				"health_status": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Health Status",
				},
				"gateway_address": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Gateway Address",
				},
			},
		},
	},
	"storage_devices": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "Storage Devices",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"serial_number": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Serial Number",
				},
				"resource_id": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Resource Id",
				},
				"address": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Address",
				},
				"model": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Model",
				},
				"microcode_version": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Microcode Version",
				},
				"resource_state": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Resource State",
				},
				"health_state": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Health State",
				},
				"ucp_systems": &schema.Schema{
					Type:        schema.TypeList,
					Computed:    true,
					Description: "List of UCP Systems",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"svp_ip": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "SVP IP",
				},
				"gateway_address": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Gateway Address",
				},
			},
		},
	},
	"ethernet_switches": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "Ethernet Switches",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"resource_id": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Resource Id",
				},
				"address": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Address",
				},
				"name": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the Ethernet Switch",
				},
				"model": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Model",
				},
				"serial_number": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Serial Number",
				},
				"firmware_version": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Firmware Version",
				},
				"resource_state": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Resource State",
				},
				"health_status": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Health Status",
				},
				"gateway_address": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Gateway Address",
				},
				"is_management": {
					Computed:    true,
					Type:        schema.TypeBool,
					Description: "Is Management",
				},
			},
		},
	},
	"fibre_channel_switches": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "Fibre Channel Switches",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"resource_id": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Resource Id",
				},
				"address": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Address",
				},
				"model": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Model",
				},
				"serial_number": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Serial Number",
				},
				"firmware_version": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Firmware Version",
				},
				"resource_state": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Resource State",
				},
				"health_state": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Health State",
				},
				"switch_name": &schema.Schema{
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the Fiber Channel Switch",
				},
				"gateway_address": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Gateway Address",
				},
			},
		},
	},
	"serial_number": {
		Computed:    true,
		Type:        schema.TypeString,
		Description: "Serial Number",
	},
	"gateway_address": {
		Computed:    true,
		Type:        schema.TypeString,
		Description: "Gateway Address",
	},
	"model": {
		Computed:    true,
		Type:        schema.TypeString,
		Description: "Model",
	},
	"vcenter": {
		Computed:    true,
		Type:        schema.TypeString,
		Description: "vcenter",
	},
	"zone": {
		Computed:    true,
		Type:        schema.TypeString,
		Description: "Zone",
	},
	"vcenter_resource_id": {
		Computed:    true,
		Type:        schema.TypeString,
		Description: "vcenter Resource Id",
	},
	"region": {
		Computed:    true,
		Type:        schema.TypeString,
		Description: "Region",
	},
	"geo_information": &schema.Schema{
		Computed:    true,
		Optional:    true,
		Description: "Geo Information",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"geo_location": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Geo Location",
				},
				"country": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Country",
				},
				"latitude": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Latitude",
				},
				"longitude": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Longitude",
				},
				"zipcode": {
					Computed:    true,
					Type:        schema.TypeString,
					Description: "Zipcode",
				},
			},
		},
	},
	"workload_type": {
		Computed:    true,
		Type:        schema.TypeString,
		Description: "Workload Type",
	},
	"result_status": {
		Computed:    true,
		Type:        schema.TypeString,
		Description: "Result Status",
	},
	"result_message": {
		Computed:    true,
		Type:        schema.TypeString,
		Description: "Result Message",
	},
	"plugin_registered": {
		Computed:    true,
		Type:        schema.TypeBool,
		Description: "Plugin Registered",
	},
	"linked": {
		Computed:    true,
		Type:        schema.TypeBool,
		Description: "Linked",
	},
}

var DataInfraUcpSystemSchema = map[string]*schema.Schema{
	"serial_number": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Serial Number of the UCP System",
	},

	"name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Name of the UCP System",
	},

	// output
	"ucp_systems": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "This is ucp systems output",
		Elem: &schema.Resource{
			Schema: InfraUcpSystemSchema,
		},
	},
}
