package sanstorage

import (
	"fmt"
	diskcache "terraform-provider-hitachi/hitachi/common/diskcache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// TODO: fix logging, debugging, errors
// GetStorageSystemInfo used to get storage system information
func (psm *sanStorageManager) GetStorageSystemInfo(detailInfoType ...string) (*sanmodel.StorageSystemInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storageInfo sanmodel.StorageSystemInfo

	// Build API suffix with detailInfoType if provided
	apiSuf := "objects/storages/instance"
	if len(detailInfoType) > 0 && detailInfoType[0] != "" {
		apiSuf = fmt.Sprintf("objects/storages/instance?detailInfoType=%s", detailInfoType[0])
		log.WriteDebug("TFDebug| GetStorageSystemInfo: API URL constructed: %s", apiSuf)
	} else {
		log.WriteDebug("TFDebug| GetStorageSystemInfo: NO detailInfoType provided, using basic API: %s", apiSuf)
	}

	// For detailInfoType calls, don't use disk cache as the response varies
	if len(detailInfoType) == 0 || detailInfoType[0] == "" {
		// read from disk cache only for basic calls
		key := psm.storageSetting.MgmtIP + ":StorageSystemInfo"
		found, _ := diskcache.Get(key, &storageInfo)
		if found {
			return &storageInfo, nil
		}
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storageInfo)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	// save this to disk cache only for basic calls (without detailInfoType)
	if len(detailInfoType) == 0 || detailInfoType[0] == "" {
		key := psm.storageSetting.MgmtIP + ":StorageSystemInfo"
		diskcache.Set(key, storageInfo)
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
