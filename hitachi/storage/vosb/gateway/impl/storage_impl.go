package vssbstorage

import (
	diskcache "terraform-provider-hitachi/hitachi/common/diskcache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/vosb/gateway/http"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
)

// GetStorageVersionInfo gets version information of vssb storage
func (psm *vssbStorageManager) GetStorageVersionInfo() (*vssbmodel.StorageVersionInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var versionInfo vssbmodel.StorageVersionInfo

	// read from disk cache first
	key := psm.storageSetting.ClusterAddress + ":StorageVersionInfo"
	found, _ := diskcache.Get(key, &versionInfo)
	if found {
		return &versionInfo, nil
	}

	apiSuf := "configuration/version"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &versionInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	// save this to disk cache
	diskcache.Set(key, versionInfo)

	return &versionInfo, nil
}

// GetHealthStatuses Obtains the health status.
func (psm *vssbStorageManager) GetHealthStatuses() (*vssbmodel.HealthStatuses, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var healthStatuses vssbmodel.HealthStatuses
	apiSuf := "objects/health-status"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &healthStatuses)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &healthStatuses, nil
}

// GetFaultDomains Obtains a list of fault domain information.
func (psm *vssbStorageManager) GetFaultDomains() (*vssbmodel.FaultDomains, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var faultDomains vssbmodel.FaultDomains
	apiSuf := "objects/fault-domains"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &faultDomains)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &faultDomains, nil
}

// GetStorageClusterInfo Obtains the storage cluster information.
func (psm *vssbStorageManager) GetStorageClusterInfo() (*vssbmodel.StorageClusterInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var clusterInfo vssbmodel.StorageClusterInfo
	apiSuf := "objects/storage"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &clusterInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &clusterInfo, nil
}

// GetDrivesInfo Obtains a list of drive information.
func (psm *vssbStorageManager) GetDrivesInfo() (*vssbmodel.Drives, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var drivesInfo vssbmodel.Drives
	apiSuf := "objects/drives"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &drivesInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &drivesInfo, nil
}
