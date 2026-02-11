package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetDynamicPools
func (psm *sanStorageManager) GetDynamicPools(isMainframe *bool, poolType string, detailInfoType ...string) (*[]sanmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOLS_BEGIN), objStorage.Serial)
	var provDynamicPools *[]provisonermodel.DynamicPool
	if len(detailInfoType) > 0 {
		provDynamicPools, err = provObj.GetDynamicPools(isMainframe, poolType, detailInfoType[0])
	} else {
		provDynamicPools, err = provObj.GetDynamicPools(isMainframe, poolType)
	}
	if err != nil {
		log.WriteDebug("TFError| error in GetDynamicPools provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_DYNAMIC_POOLS_FAILED), objStorage.Serial)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconDynamicPools := []sanmodel.DynamicPool{}
	err = copier.Copy(&reconDynamicPools, provDynamicPools)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOLS_END), objStorage.Serial)
	return &reconDynamicPools, nil
}

// GetDynamicPoolById
func (psm *sanStorageManager) GetDynamicPoolById(poolId int) (*sanmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOL_ID_BEGIN), poolId, objStorage.Serial)
	provDynamicPool, err := provObj.GetDynamicPoolById(poolId)
	if err != nil {
		log.WriteDebug("TFError| error in GetDynamicPoolById provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_DYNAMIC_POOL_ID_FAILED), poolId, objStorage.Serial)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconDynamicPool := sanmodel.DynamicPool{}
	err = copier.Copy(&reconDynamicPool, provDynamicPool)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOL_ID_END), poolId, objStorage.Serial)
	return &reconDynamicPool, nil
}
