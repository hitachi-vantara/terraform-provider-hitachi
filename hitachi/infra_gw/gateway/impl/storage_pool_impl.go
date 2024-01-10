package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetStoragePools gets storage pools information
func (psm *infraGwManager) GetStoragePools(id string) (*model.StoragePools, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storagePools model.StoragePools

	apiSuf := fmt.Sprintf("/storage/devices/%s/pools", id)
	err := httpmethod.GetCall(psm.setting, apiSuf, &storagePools)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storagePools, nil
}

// GetStoragePool gets storage pool information
func (psm *infraGwManager) GetStoragePool(id, poolId string) (*model.StoragePool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storagePool model.StoragePool

	apiSuf := fmt.Sprintf("/storage/devices/%s/pools/%s", id, poolId)
	err := httpmethod.GetCall(psm.setting, apiSuf, &storagePool)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storagePool, nil
}
