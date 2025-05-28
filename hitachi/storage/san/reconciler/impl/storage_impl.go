package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetStorageSystemInfo used to get Storage systme information
func (psm *sanStorageManager) GetStorageSystemInfo() (*sanmodel.StorageSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_BEGIN), psm.storageSetting.MgmtIP)
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
	provStorageSysInfo, err := provObj.GetStorageSystemInfo()
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_SYSTEM_FAILED), objStorage.MgmtIP)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilertStorageSystem := sanmodel.StorageSystem{}
	err = copier.Copy(&reconcilertStorageSystem, provStorageSysInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_END), objStorage.MgmtIP)
	return &reconcilertStorageSystem, nil
}

// GetStorageSystem used to get Storage systme information ( Storage info + Capacity)
func (psm *sanStorageManager) GetStorageSystem() (*sanmodel.StorageSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_BEGIN), psm.storageSetting.MgmtIP)
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
	provStorageSysInfo, err := provObj.GetStorageSystem()
	if err != nil {
		log.WriteInfo(mc.GetMessage(mc.ERR_GET_STORAGE_SYSTEM_FAILED), objStorage.MgmtIP)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilertStorageSystem := sanmodel.StorageSystem{}
	err = copier.Copy(&reconcilertStorageSystem, provStorageSysInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_END), objStorage.MgmtIP)
	return &reconcilertStorageSystem, nil
}
