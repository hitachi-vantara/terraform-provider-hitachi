package sanstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetPools
func (psm *sanStorageManager) GetPools() (*[]sanmodel.Pool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var pools sanmodel.Pools
	apiSuf := fmt.Sprintf("objects/pools")
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &pools)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &pools.Data, nil
}
