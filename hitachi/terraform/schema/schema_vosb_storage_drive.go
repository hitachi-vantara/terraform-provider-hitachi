package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Define the Drive struct schema
var VssbStorageDriveSchema = map[string]*schema.Schema{
	"id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the drive.",
	},
	"wwid": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The WWW ID of the drive.",
	},
	"status_summary": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Summary of the status of the drive.",
	},
	"status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The current status of the drive.",
	},
	"type_code": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The type code of the drive.",
	},
	"serial_number": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The serial number of the drive.",
	},
	"storage_node_id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the storage node where the drive is located.",
	},
	"device_file_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The device file name of the drive.",
	},
	"vendor_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The vendor name of the drive.",
	},
	"firmware_revision": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The firmware revision of the drive.",
	},
	"locator_led_status": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The locator LED status of the drive.",
	},
	"drive_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The type of the drive.",
	},
	"drive_capacity": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The capacity of the drive in GB.",
	},
}

// Define the Data source schema
var ResourceVssbStorageDriveSchema = map[string]*schema.Schema{
	// Input
	"vosb_address": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The host name or the IP address (IPv4) of VSP One SDS Block.",
	},
	"status": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"", "Offline", "Normal", "TemporaryBlockage", "Blockage"}, true),
		Description:  "Filter the drives by their status. Allowed values (case-insensitive): empty string, Offline, Normal, TemporaryBlockage, Blockage.",
	},
	// Output
	"drives": &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of drives.",
		Elem: &schema.Resource{
			Schema: VssbStorageDriveSchema,
		},
	},
}
