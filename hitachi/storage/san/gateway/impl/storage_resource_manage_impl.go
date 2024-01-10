package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// LockResources used to lock all resources on storage device
func (psm *sanStorageManager) LockResources(reqBody sanmodel.LockResourcesReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "services/resource-group-service/actions/lock/invoke"
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in LockResources - %s, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// UnlockResources used to unlock all resources on storage device
// FIX ME: Currently it's not working as it requires same session token id who has locked the resource
func (psm *sanStorageManager) UnlockResources() error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "services/resource-group-service/actions/unlock/invoke"
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteDebug("TFError| error in UnlockResources - %s, err: %v", apiSuf, err)
		return err
	}

	return nil
}
