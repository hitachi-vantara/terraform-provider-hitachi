package sanstorage

import (
	"fmt"
	// "time"
	// "encoding/json"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	// "terraform-provider-hitachi/hitachi/common/utils"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

// TODO: fix logging, debugging, errors

func GetStorageSystemInfo(storageSetting sanmodel.StorageDeviceSettings) (*sanmodel.StorageSystemInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storageInfo sanmodel.StorageSystemInfo
	apiSuf := fmt.Sprintf("objects/storages/instance")
	err := GetCall(storageSetting, apiSuf, &storageInfo)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	return &storageInfo, nil
}

// only for VSP 5000 series
func GetStorageSystemSummary(storageSetting sanmodel.StorageDeviceSettings) (*sanmodel.StorageSystemSummary, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storageSystemSummary sanmodel.StorageSystemSummary
	apiSuf := fmt.Sprintf("objects/storages/instance")
	err := GetCall(storageSetting, apiSuf, &storageSystemSummary)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	return &storageSystemSummary, nil
}

func GetStorageCapacity(storageSetting sanmodel.StorageDeviceSettings) (*sanmodel.StorageCapacity, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storageCapacity sanmodel.StorageCapacity
	apiSuf := fmt.Sprintf("objects/total-capacities/instance")
	err := GetCall(storageSetting, apiSuf, &storageCapacity)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	return &storageCapacity, nil
}
