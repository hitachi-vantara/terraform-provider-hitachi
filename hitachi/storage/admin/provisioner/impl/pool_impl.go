package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/provisioner/message-catalog"
)

func (psm *adminStorageManager) GetAdminPoolList(params gwymodel.AdminPoolListParams) (*gwymodel.AdminPoolListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_POOLS_BEGIN), psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	log.WriteDebug("TFDebug| QueryParams:%+v\n", params)

	poolList, err := gatewayObj.GetAdminPoolList(params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_POOLS_FAILED), psm.storageSetting.MgmtIP)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_POOLS_END), psm.storageSetting.MgmtIP)
	return poolList, nil
}

func (psm *adminStorageManager) GetAdminPoolInfo(poolID int) (*gwymodel.AdminPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_POOL_INFO_BEGIN), poolID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	poolInfo, err := gatewayObj.GetAdminPoolInfo(poolID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_POOL_INFO_FAILED), poolID, psm.storageSetting.MgmtIP)
		return nil, err
	}

	log.WriteDebug("TFDebug| PoolInfo: %+v\n", poolInfo)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_POOL_INFO_END), poolID, psm.storageSetting.MgmtIP)

	return poolInfo, nil
}

func (psm *adminStorageManager) CreateAdminPool(params gwymodel.CreateAdminPoolParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_POOL_BEGIN), params.Name, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| CreateParams:%+v\n", params)

	err = gatewayObj.CreateAdminPool(params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_POOL_FAILED), params.Name, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_POOL_END), params.Name, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) UpdateAdminPool(poolID int, params gwymodel.UpdateAdminPoolParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_POOL_BEGIN), poolID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| UpdateParams:%+v\n", params)

	err = gatewayObj.UpdateAdminPool(poolID, params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_POOL_FAILED), poolID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_POOL_END), poolID, psm.storageSetting.MgmtIP)
	return nil
}

func (psm *adminStorageManager) ExpandAdminPool(poolID int, params gwymodel.ExpandAdminPoolParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Expanding admin pool ID %d on %s", poolID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| ExpandParams:%+v\n", params)

	err = gatewayObj.ExpandAdminPool(poolID, params)
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

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_POOL_BEGIN), poolID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	err = gatewayObj.DeleteAdminPool(poolID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_POOL_FAILED), poolID, psm.storageSetting.MgmtIP, err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_POOL_END), poolID, psm.storageSetting.MgmtIP)
	return nil
}
