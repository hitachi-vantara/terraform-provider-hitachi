package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetStoragePorts used to get Storage ports
func (psm *sanStorageManager) GetStoragePorts() (*[]sanmodel.StoragePort, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_BEGIN), objStorage.Serial)
	provStoragePorts, err := provObj.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragePorts provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_PORTS_FAILED), objStorage.Serial)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconStoragePorts := []sanmodel.StoragePort{}
	err = copier.Copy(&reconStoragePorts, provStoragePorts)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_END), objStorage.Serial)
	return &reconStoragePorts, nil
}

// GetStoragePortByPortId used to get storage port information for given portId
func (psm *sanStorageManager) GetStoragePortByPortId(portId string) (*sanmodel.StoragePort, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_PORTID_BEGIN), portId, objStorage.Serial)
	provStoragePort, err := provObj.GetStoragePortByPortId(portId)
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragePorts provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_PORTS_PORTID_FAILED), portId, objStorage.Serial)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconStoragePort := sanmodel.StoragePort{}
	err = copier.Copy(&reconStoragePort, provStoragePort)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_PORTID_END), portId, objStorage.Serial)
	return &reconStoragePort, nil
}
