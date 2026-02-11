package sanstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetStoragePorts used to get storage ports information
func (psm *sanStorageManager) GetStoragePorts() (*[]sanmodel.StoragePort, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_BEGIN), objStorage.Serial)
	storagePorts, err := gatewayObj.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragePorts gateway call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_PORTS_FAILED), objStorage.Serial)
		return nil, err
	}

	provStoragePorts := []sanmodel.StoragePort{}
	err = copier.Copy(&provStoragePorts, storagePorts)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_END), objStorage.Serial)

	return &provStoragePorts, nil
}

// GetStoragePortByPortId used to get storage port information for given portId
func (psm *sanStorageManager) GetStoragePortByPortId(portId string) (*sanmodel.StoragePort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storagePorts, err := psm.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragePorts gateway call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_PORTS_FAILED), psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_PORTID_BEGIN), portId, psm.storageSetting.Serial)

	for _, sp := range *storagePorts {
		if sp.PortId == portId {
			log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_PORTID_END), portId, psm.storageSetting.Serial)
			return &sp, nil
		}
	}

	log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_PORTS_PORTID_FAILED), portId, psm.storageSetting.Serial)

	return nil, fmt.Errorf("The portId %s is not present on storage serial %d.", portId, psm.storageSetting.Serial)
}
