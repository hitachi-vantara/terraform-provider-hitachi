package terraform

import (
	"errors"
	"strings"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetAllDrives fetches all drives without any arguments
func GetAllDrives(d *schema.ResourceData) (*[]terraformmodel.Drive, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_address").(string)

	// Fetch settings from cache based on VOSB block address
	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return nil, errors.New("failed to retrieve VOSB settings from cache")
	}

	// Initialize storage device settings
	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	// Create a reconcilers object using the settings
	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteError("Error initializing reconcilers object", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DRIVES_BEGIN))

	// Retrieve all drives from the reconciler (this is the key action)
	var canonicalStatusValues = map[string]string{
		"offline":           "Offline",
		"normal":            "Normal",
		"blockage":          "Blockage",
		"temporaryblockage": "TemporaryBlockage",
		"":                  "",
	}

	rawStatus := d.Get("status").(string)
	status := canonicalStatusValues[strings.ToLower(rawStatus)]

	driveInfo, err := reconObj.GetDrivesInfo(status)
	if err != nil {
		log.WriteError("Error fetching drives from reconciler", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_DRIVES_FAILED))
		return nil, err
	}

	// Map the drive information into the required Terraform model (Drive)
	terraformDrives := []terraformmodel.Drive{}
	for _, drive := range driveInfo.Data {
		terraformDrives = append(terraformDrives, terraformmodel.Drive{
			Id:               strings.TrimSpace(drive.Id),
			WwId:             strings.TrimSpace(drive.WwId),
			StatusSummary:    strings.TrimSpace(drive.StatusSummary),
			Status:           strings.TrimSpace(drive.Status),
			TypeCode:         strings.TrimSpace(drive.TypeCode),
			SerialNumber:     strings.TrimSpace(drive.SerialNumber),
			StorageNodeId:    strings.TrimSpace(drive.StorageNodeId),
			DeviceFileName:   strings.TrimSpace(drive.DeviceFileName),
			VendorName:       strings.TrimSpace(drive.VendorName),
			FirmwareRevision: strings.TrimSpace(drive.FirmwareRevision),
			LocatorLedStatus: strings.TrimSpace(drive.LocatorLedStatus),
			DriveType:        strings.TrimSpace(drive.DriveType),
			DriveCapacity:    drive.DriveCapacity,
		})
	}

	log.WriteInfo("Successfully retrieved drives")
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DRIVES_END))

	return &terraformDrives, nil
}
