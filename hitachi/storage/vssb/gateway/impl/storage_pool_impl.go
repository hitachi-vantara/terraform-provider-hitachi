package vssbstorage

import (
	"fmt"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/vssb/gateway/http"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/gateway/model"
)

// GetAllStoragePools gets all storage pool details
func (psm *vssbStorageManager) GetAllStoragePools() (*vssbmodel.StoragePools, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storagePools vssbmodel.StoragePools
	apiSuf := "objects/pools"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storagePools)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storagePools, nil
}

// GetStoragePoolsByPoolNames gets storage pools by pool names
func (psm *vssbStorageManager) GetStoragePoolsByPoolNames(poolNames []string) (*vssbmodel.StoragePools, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storagePools vssbmodel.StoragePools
	names := strings.Join(poolNames, ",")
	apiSuf := fmt.Sprintf("objects/pools?names=%s", names)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storagePools)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storagePools, nil
}
