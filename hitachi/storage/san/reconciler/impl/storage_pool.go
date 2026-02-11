package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetPools
func (psm *sanStorageManager) GetPools() (*[]sanmodel.Pool, error) {
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
	provPools, err := provObj.GetPools()
	if err != nil {
		log.WriteDebug("TFError| error in GetPools provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_DYNAMIC_POOLS_FAILED), objStorage.Serial)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconPools := []sanmodel.Pool{}
	err = copier.Copy(&reconPools, provPools)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOLS_END), objStorage.Serial)
	return &reconPools, nil
}
