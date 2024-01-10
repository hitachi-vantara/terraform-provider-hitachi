package sanstorage

import (
	spmanager "terraform-provider-hitachi/hitachi/storage/san/provisioner"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
)

// sanStorageManager contain information for storage setting
type sanStorageManager struct {
	storageSetting sanmodel.StorageDeviceSettings
}

// A private function to construct an newSanStorageManagerEx
func newSanStorageManagerEx(storageSetting sanmodel.StorageDeviceSettings) (*sanStorageManager, error) {

	psm := &sanStorageManager{
		storageSetting: sanmodel.StorageDeviceSettings{
			Serial:   storageSetting.Serial,
			Username: storageSetting.Username,
			Password: storageSetting.Password,
			MgmtIP:   storageSetting.MgmtIP,
		},
	}

	return psm, nil
}

// NewEx returns a new storage Provisioner
func NewEx(storageSetting sanmodel.StorageDeviceSettings) (spmanager.SanStorageManager, error) {
	return newSanStorageManagerEx(storageSetting)
}
