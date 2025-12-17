package admin

import (
	telemetry "terraform-provider-hitachi/hitachi/common/telemetry"
	manager "terraform-provider-hitachi/hitachi/storage/admin/provisioner"
	model "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"
)

// adminStorageManager contain information for storage setting
type adminStorageManager struct {
	storageSetting model.StorageDeviceSettings
}

// A private function to construct an newAdminStorageManagerEx
func newAdminStorageManagerEx(storageSetting model.StorageDeviceSettings) (*adminStorageManager, error) {
	psm := &adminStorageManager{
		storageSetting: model.StorageDeviceSettings{
			Serial:                  storageSetting.Serial,
			Username:                storageSetting.Username,
			Password:                storageSetting.Password,
			MgmtIP:                  storageSetting.MgmtIP,
			TerraformResourceMethod: telemetry.GetTerraformCallStackInfo(),
		},
	}

	return psm, nil
}

// NewEx returns a new storage Provisioner
func NewEx(storageSetting model.StorageDeviceSettings) (manager.AdminStorageManager, error) {
	return newAdminStorageManagerEx(storageSetting)
}
