package admin

import (
	"fmt"
	"net/url"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/admin/gateway/http"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func (psm *adminStorageManager) GetAdminServerList(params gwymodel.AdminServerListParams) (*gwymodel.AdminServerListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Getting admin server list from %s", psm.storageSetting.MgmtIP)

	var result gwymodel.AdminServerListResponse

	// Build query parameters
	queryParams := url.Values{}
	if params.Nickname != nil {
		queryParams.Add("nickname", *params.Nickname)
	}
	if params.HbaWwn != nil {
		queryParams.Add("hbaWwn", *params.HbaWwn)
	}
	if params.IscsiName != nil {
		queryParams.Add("iscsiName", *params.IscsiName)
	}

	apiSuf := fmt.Sprintf("objects/servers?%s", queryParams.Encode())
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &result)
	if err != nil {
		log.WriteError("Failed to get admin server list from %s: %v", psm.storageSetting.MgmtIP, err)
		return nil, err
	}

	log.WriteInfo("Successfully retrieved admin server list from %s", psm.storageSetting.MgmtIP)
	return &result, nil
}

func (psm *adminStorageManager) GetAdminServerInfo(serverID int) (*gwymodel.AdminServerInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Getting admin server info for ID %d from %s", serverID, psm.storageSetting.MgmtIP)

	var result gwymodel.AdminServerInfo

	apiSuf := fmt.Sprintf("objects/servers/%d", serverID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &result)
	if err != nil {
		log.WriteError("Failed to get admin server info for ID %d from %s: %v", serverID, psm.storageSetting.MgmtIP, err)
		return nil, err
	}

	log.WriteInfo("Successfully retrieved admin server info for ID %d from %s", serverID, psm.storageSetting.MgmtIP)
	return &result, nil
}

func (psm *adminStorageManager) CreateAdminServer(params gwymodel.CreateAdminServerParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Creating admin server '%s' on %s", params.ServerNickname, psm.storageSetting.MgmtIP)

	apiSuf := "objects/servers"
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError("Failed to create admin server '%s' on %s: %v", params.ServerNickname, psm.storageSetting.MgmtIP, err)
		return err
	}
	log.WriteInfo("Successfully created admin server '%s' on %s", params.ServerNickname, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) UpdateAdminServer(serverID int, params gwymodel.UpdateAdminServerParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Updating admin server ID %d on %s", serverID, psm.storageSetting.MgmtIP)

	apiSuf := fmt.Sprintf("objects/servers/%d", serverID)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError("Failed to update admin server ID %d on %s: %v", serverID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo("Successfully updated admin server ID %d on %s", serverID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) DeleteAdminServer(serverID int, params gwymodel.DeleteAdminServerParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Deleting admin server ID %d on %s", serverID, psm.storageSetting.MgmtIP)

	apiSuf := fmt.Sprintf("objects/servers/%d", serverID)

	// Prepare request body with keepLunConfig parameter
	var requestBody interface{}
	requestBody = map[string]interface{}{
		"keepLunConfig": params.KeepLunConfig,
	}

	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, requestBody)
	if err != nil {
		log.WriteError("Failed to delete admin server ID %d on %s: %v", serverID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo("Successfully deleted admin server ID %d on %s", serverID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) SetAdminServerPath(serverID int, params gwymodel.SetAdminServerPathParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Setting admin server path for ID %d on %s", serverID, psm.storageSetting.MgmtIP)

	apiSuf := fmt.Sprintf("objects/servers/%d/paths", serverID)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError("Failed to set admin server path for ID %d on %s: %v", serverID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo("Successfully set admin server path for ID %d on %s", serverID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) DeleteAdminServerPath(serverID int, params gwymodel.DeleteAdminServerPathParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Deleting admin server path for ID %d on %s", serverID, psm.storageSetting.MgmtIP)

	// Build the path parameter for the URL
	var pathParam string
	if params.HbaWwn != "" {
		pathParam = fmt.Sprintf("%s,%s", params.HbaWwn, params.PortId)
	} else {
		pathParam = fmt.Sprintf("%s,%s", params.IscsiName, params.PortId)
	}

	apiSuf := fmt.Sprintf("objects/servers/%d/paths/%s", serverID, pathParam)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError("Failed to delete admin server path for ID %d on %s: %v", serverID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo("Successfully deleted admin server path for ID %d on %s", serverID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) GetAdminServerPath(params gwymodel.AdminServerPathParams) (*gwymodel.AdminServerPathInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Getting admin server path for ID %d on %s", params.ServerID, psm.storageSetting.MgmtIP)

	// Build the path parameter for the URL
	var pathParam string
	if params.HbaWwn != "" {
		pathParam = fmt.Sprintf("%s,%s", params.HbaWwn, params.PortId)
	} else {
		pathParam = fmt.Sprintf("%s,%s", params.IscsiName, params.PortId)
	}

	var result gwymodel.AdminServerPathInfo

	apiSuf := fmt.Sprintf("objects/servers/%d/paths/%s", params.ServerID, pathParam)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &result)
	if err != nil {
		log.WriteError("Failed to get admin server path for ID %d on %s: %v", params.ServerID, psm.storageSetting.MgmtIP, err)
		return nil, err
	}

	log.WriteInfo("Successfully retrieved admin server path for ID %d on %s", params.ServerID, psm.storageSetting.MgmtIP)
	return &result, nil
}
