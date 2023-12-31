package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// TODO: fix logging, debugging, errors
// GetStorageSystemInfo used to get storage system information
func (psm *sanStorageManager) GetStorageSystemInfo() (*sanmodel.StorageSystemInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storageInfo sanmodel.StorageSystemInfo
	apiSuf := "objects/storages/instance"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storageInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in objects/storages/instance API call, err: %v", err)
		return nil, err
	}
	return &storageInfo, nil
}

// GetStorageCapacity get storage capacity information
func (psm *sanStorageManager) GetStorageCapacity() (*sanmodel.StorageCapacity, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storageCapacity sanmodel.StorageCapacity
	apiSuf := "objects/total-capacities/instance"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storageCapacity)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	return &storageCapacity, nil
}
