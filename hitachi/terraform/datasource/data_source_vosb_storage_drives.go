package terraform

import (
	"context"
	"fmt"
	"strconv"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceVssbStorageDrives defines the Terraform data source for retrieving drives
func DataSourceVssbStorageDrives() *schema.Resource {
	return &schema.Resource{
		Description: "VOS Block Storage Drives: Obtains a list of drive information.",
		ReadContext: DataSourceVssbStorageDrivesRead,
		Schema:      schemaimpl.ResourceVssbStorageDriveSchema, // Use the schema we defined for the drives
	}
}

// DataSourceVssbStorageDrivesRead is the function that retrieves the drive information
func DataSourceVssbStorageDrivesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Call the GetAllStorageDrives function to get drive information
	drives, err := impl.GetAllDrives(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Convert the drive data to the appropriate format for Terraform
	driveList := []map[string]interface{}{}
	for _, drive := range *drives {
		// Convert each drive into the schema map format
		eachDrive := map[string]interface{}{
			"id":                 drive.Id,
			"wwid":               drive.WwId,
			"status_summary":     drive.StatusSummary,
			"status":             drive.Status,
			"type_code":          drive.TypeCode,
			"serial_number":      drive.SerialNumber,
			"storage_node_id":    drive.StorageNodeId,
			"device_file_name":   drive.DeviceFileName,
			"vendor_name":        drive.VendorName,
			"firmware_revision":  drive.FirmwareRevision,
			"locator_led_status": drive.LocatorLedStatus,
			"drive_type":         drive.DriveType,
			"drive_capacity":     fmt.Sprintf("%d GB", drive.DriveCapacity),
		}
		// Append to the list
		driveList = append(driveList, eachDrive)
	}

	// Set the drives output in the Terraform state
	if err := d.Set("drives", driveList); err != nil {
		return diag.FromErr(err)
	}

	// Set the ID of the data source to ensure it's unique for this operation
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	log.WriteInfo("Successfully retrieved drive information")

	return nil
}
