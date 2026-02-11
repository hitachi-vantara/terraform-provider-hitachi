package admin

import (
	"fmt"
	"net/url"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/admin/gateway/http"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func (psm *adminStorageManager) GetAdminPoolList(params gwymodel.AdminPoolListParams) (*gwymodel.AdminPoolListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Getting admin pool list on %s", psm.storageSetting.MgmtIP)

	// Build query parameters
	queryParams := url.Values{}
	if params.Name != nil {
		queryParams.Set("name", *params.Name)
	}
	if params.Status != nil {
		queryParams.Set("status", *params.Status)
	}
	if params.ConfigStatus != nil {
		queryParams.Set("configStatus", *params.ConfigStatus)
	}

	apiSuf := "objects/pools"
	if len(queryParams) > 0 {
		apiSuf += "?" + queryParams.Encode()
	}

	var result gwymodel.AdminPoolListResponse
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &result)
	if err != nil {
		log.WriteError("Failed to get admin pool list on %s: %v", psm.storageSetting.MgmtIP, err)
		return nil, err
	}

	log.WriteInfo("Successfully retrieved admin pool list on %s", psm.storageSetting.MgmtIP)
	return &result, nil
}

func (psm *adminStorageManager) GetAdminPoolInfo(poolID int) (*gwymodel.AdminPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Getting admin pool info for ID %d on %s", poolID, psm.storageSetting.MgmtIP)

	var result gwymodel.AdminPool
	apiSuf := fmt.Sprintf("objects/pools/%d", poolID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &result)
	if err != nil {
		log.WriteError("Failed to get admin pool info for ID %d on %s: %v", poolID, psm.storageSetting.MgmtIP, err)
		return nil, err
	}

	if result.SavingEffects.TotalEfficiency >= 9223372036854775806 {
		result.SavingEffects.TotalEfficiency = 0
	}

	log.WriteInfo("Successfully retrieved admin pool info for ID %d on %s", poolID, psm.storageSetting.MgmtIP)
	return &result, nil
}

func (psm *adminStorageManager) CreateAdminPool(params gwymodel.CreateAdminPoolParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Creating admin pool %s on %s", params.Name, psm.storageSetting.MgmtIP)

	apiSuf := "objects/pools"
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError("Failed to create admin pool %s on %s: %v", params.Name, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo("Successfully created admin pool %s on %s", params.Name, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) UpdateAdminPool(poolID int, params gwymodel.UpdateAdminPoolParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Updating admin pool ID %d on %s", poolID, psm.storageSetting.MgmtIP)

	apiSuf := fmt.Sprintf("objects/pools/%d", poolID)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError("Failed to update admin pool ID %d on %s: %v", poolID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo("Successfully updated admin pool ID %d on %s", poolID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) ExpandAdminPool(poolID int, params gwymodel.ExpandAdminPoolParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Expanding admin pool ID %d on %s", poolID, psm.storageSetting.MgmtIP)

	apiSuf := fmt.Sprintf("objects/pools/%d/actions/expand/invoke", poolID)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError("Failed to expand admin pool ID %d on %s: %v", poolID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo("Successfully expanded admin pool ID %d on %s", poolID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) DeleteAdminPool(poolID int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Deleting admin pool ID %d on %s", poolID, psm.storageSetting.MgmtIP)

	apiSuf := fmt.Sprintf("objects/pools/%d", poolID)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError("Failed to delete admin pool ID %d on %s: %v", poolID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo("Successfully deleted admin pool ID %d on %s", poolID, psm.storageSetting.MgmtIP)
	return nil
}
