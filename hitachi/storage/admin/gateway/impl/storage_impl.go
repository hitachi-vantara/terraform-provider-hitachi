package admin

import (
	"fmt"
	"terraform-provider-hitachi/hitachi/common/diskcache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/admin/gateway/http"
	gatewaymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

const WITH_EST_CONFIGURABLE_CAPACITIES = "withEstimatedConfigurableCapacities"

// GetStorageAdminInfo Obtains the storage admin information.
func (psm *adminStorageManager) GetStorageAdminInfo(configurable_capacities bool) (*gatewaymodel.StorageAdminInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var adminInfo gatewaymodel.StorageAdminInfo

	// read from disk cache
	key := psm.storageSetting.MgmtIP + ":StorageAdminInfo"
	found, _ := diskcache.Get(key, &adminInfo)
	if found {
		log.WriteDebug("TFDebug| Data found in disk cache: %t, Data: %+v, configurable_capacities: %t", found, adminInfo, configurable_capacities)
		return &adminInfo, nil
	}

	adminInfo = gatewaymodel.StorageAdminInfo{}
	log.WriteDebug("TFDebug| Data not found in disk cache, call API")
	apiSuf := fmt.Sprintf("objects/storage?%s=%t", WITH_EST_CONFIGURABLE_CAPACITIES, configurable_capacities)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &adminInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	log.WriteDebug("TFDebug| Data for call API: %+v", adminInfo)

	// save this to disk cache
	diskcache.Set(key, adminInfo)

	return &adminInfo, nil
}
