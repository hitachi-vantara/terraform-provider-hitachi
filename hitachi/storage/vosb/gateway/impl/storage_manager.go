package vssbstorage

import (
	telemetry "terraform-provider-hitachi/hitachi/common/telemetry"
	spmanager "terraform-provider-hitachi/hitachi/storage/vosb/gateway"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
)

// vssbStorageManager contain information for storage setting
type vssbStorageManager struct {
	storageSetting vssbmodel.StorageDeviceSettings
}

// A private function to construct an newVssbStorageManagerEx
func newVssbStorageManagerEx(storageSetting vssbmodel.StorageDeviceSettings) (*vssbStorageManager, error) {
	psm := &vssbStorageManager{
		storageSetting: vssbmodel.StorageDeviceSettings{
			Username:                storageSetting.Username,
			Password:                storageSetting.Password,
			ClusterAddress:          storageSetting.ClusterAddress,
			TerraformResourceMethod: telemetry.GetTerraformCallStackInfo(),
		},
	}
	return psm, nil
}

// NewEx returns a new storage Provisioner
func NewEx(storageSetting vssbmodel.StorageDeviceSettings) (spmanager.VssbStorageManager, error) {
	return newVssbStorageManagerEx(storageSetting)
}
