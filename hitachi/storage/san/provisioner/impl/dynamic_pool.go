package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetDynamicPools
func (psm *sanStorageManager) GetDynamicPools(isMainframe *bool, poolType string, detailInfoType ...string) (*[]sanmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteDebug("TFDebug| Storage Serial:%d, ManagementIP:%s\n", psm.storageSetting.Serial, psm.storageSetting.MgmtIP)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOLS_BEGIN), objStorage.Serial)
	var dynamicPools *[]sangatewaymodel.DynamicPool
	if len(detailInfoType) > 0 {
		dynamicPools, err = gatewayObj.GetDynamicPools(isMainframe, poolType, detailInfoType[0])
	} else {
		dynamicPools, err = gatewayObj.GetDynamicPools(isMainframe, poolType)
	}
	if err != nil {
		log.WriteDebug("TFError| error in GetDynamicPools gateway call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_DYNAMIC_POOLS_FAILED), objStorage.Serial)
		return nil, err
	}

	provDynamicPools := []sanmodel.DynamicPool{}
	err = copier.Copy(&provDynamicPools, dynamicPools)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOLS_END), objStorage.Serial)

	return &provDynamicPools, nil
}

// GetDynamicPoolById
func (psm *sanStorageManager) GetDynamicPoolById(poolId int) (*sanmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteDebug("TFDebug| Storage Serial:%d, ManagementIP:%s\n", psm.storageSetting.Serial, psm.storageSetting.MgmtIP)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOL_ID_BEGIN), poolId, objStorage.Serial)
	dynamicPool, err := gatewayObj.GetDynamicPoolById(poolId)
	if err != nil {
		log.WriteDebug("TFError| error in GetDynamicPoolById gateway call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_DYNAMIC_POOL_ID_FAILED), poolId, objStorage.Serial)
		return nil, err
	}

	provDynamicPool := sanmodel.DynamicPool{}
	err = copier.Copy(&provDynamicPool, dynamicPool)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOL_ID_END), poolId, objStorage.Serial)

	return &provDynamicPool, nil
}
