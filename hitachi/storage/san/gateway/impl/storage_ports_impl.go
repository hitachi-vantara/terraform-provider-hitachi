package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetStoragePorts used to get storage ports information
func (psm *sanStorageManager) GetStoragePorts() (*[]sanmodel.StoragePort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storagePorts sanmodel.StoragePorts
	apiSuf := "objects/ports"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storagePorts)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in objects/ports API call, err: %v", err)
		return nil, err
	}
	return &storagePorts.Data, nil
}
