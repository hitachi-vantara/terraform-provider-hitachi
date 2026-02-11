package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/provisioner/message-catalog"
)

func (psm *adminStorageManager) GetAdminServerList(queryParams gwymodel.AdminServerListParams) (*gwymodel.AdminServerListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVERS_BEGIN), psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	log.WriteDebug("TFDebug| QueryParams:%+v\n", queryParams)

	serverList, err := gatewayObj.GetAdminServerList(queryParams)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVERS_FAILED), psm.storageSetting.MgmtIP)
		return nil, err
	}

	// log.WriteDebug("TFDebug| ServerList: %+v\n", serverList)

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVERS_END), psm.storageSetting.MgmtIP)
	return serverList, nil
}

func (psm *adminStorageManager) GetAdminServerInfo(serverID int) (*gwymodel.AdminServerInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_INFO_BEGIN), serverID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	serverInfo, err := gatewayObj.GetAdminServerInfo(serverID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_INFO_FAILED), serverID, psm.storageSetting.MgmtIP)
		return nil, err
	}

	log.WriteDebug("TFDebug| ServerInfo: %+v\n", serverInfo)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_INFO_END), serverID, psm.storageSetting.MgmtIP)

	return serverInfo, nil
}

func (psm *adminStorageManager) CreateAdminServer(params gwymodel.CreateAdminServerParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_SERVER_BEGIN), params.ServerNickname, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| CreateParams:%+v\n", params)

	err = gatewayObj.CreateAdminServer(params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_SERVER_FAILED), params.ServerNickname, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_SERVER_END), params.ServerNickname, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) UpdateAdminServer(serverID int, params gwymodel.UpdateAdminServerParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_SERVER_BEGIN), serverID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| UpdateParams:%+v\n", params)

	err = gatewayObj.UpdateAdminServer(serverID, params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_SERVER_FAILED), serverID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_SERVER_END), serverID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) DeleteAdminServer(serverID int, params gwymodel.DeleteAdminServerParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SERVER_BEGIN), serverID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| DeleteParams:%+v\n", params)

	err = gatewayObj.DeleteAdminServer(serverID, params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_SERVER_FAILED), serverID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SERVER_END), serverID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) SetAdminServerPath(serverID int, params gwymodel.SetAdminServerPathParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_SERVER_PATH_BEGIN), serverID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| SetPathParams:%+v\n", params)

	err = gatewayObj.SetAdminServerPath(serverID, params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_SET_SERVER_PATH_FAILED), serverID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_SERVER_PATH_END), serverID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) DeleteAdminServerPath(serverID int, params gwymodel.DeleteAdminServerPathParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SERVER_PATH_BEGIN), serverID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| DeletePathParams:%+v\n", params)

	err = gatewayObj.DeleteAdminServerPath(serverID, params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_SERVER_PATH_FAILED), serverID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SERVER_PATH_END), serverID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) GetAdminServerPath(params gwymodel.AdminServerPathParams) (*gwymodel.AdminServerPathInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_PATH_BEGIN), params.ServerID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	log.WriteDebug("TFDebug| GetPathParams:%+v\n", params)

	result, err := gatewayObj.GetAdminServerPath(params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_PATH_FAILED), params.ServerID, psm.storageSetting.MgmtIP, err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_PATH_END), params.ServerID, psm.storageSetting.MgmtIP)
	return result, nil
}
