package sanstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetDynamicPools
func (psm *sanStorageManager) GetDynamicPools() (*[]sanmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var dynamicPools sanmodel.DynamicPools
	apiSuf := fmt.Sprintf("objects/pools?poolType=DP")
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &dynamicPools)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &dynamicPools.Data, nil
}

// GetDynamicPool by id
func (psm *sanStorageManager) GetDynamicPoolById(poolId int) (*sanmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var dynamicPool sanmodel.DynamicPool
	apiSuf := fmt.Sprintf("objects/pools/%d", poolId)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &dynamicPool)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &dynamicPool, nil
}
