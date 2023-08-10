package vssbstorage

import (
	"fmt"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/message-catelog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetStoragePorts gets ports information of vssb storage
func (psm *vssbStorageManager) GetStoragePorts() (*vssbmodel.StoragePorts, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_PORTS_BEGIN))
	provStoragePorts, err := provObj.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragePorts provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_PORTS_FAILED))
		return nil, err
	}

	// Converting Prov to Reconciler
	reconStoragePorts := vssbmodel.StoragePorts{}
	err = copier.Copy(&reconStoragePorts, provStoragePorts)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_PORTS_END))

	return &reconStoragePorts, nil
}

// GetPort gets ports information of vssb storage with port auth settings
func (psm *vssbStorageManager) GetPort(portName string) (*vssbmodel.StoragePort, *vssbmodel.PortAuthSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_BEGIN), portName)
	allStoragePorts, err := psm.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragePorts provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_FAILED), portName)
		return nil, nil, err
	}

	portId := ""
	for _, port := range allStoragePorts.Data {
		if strings.EqualFold(portName, port.Nickname) {
			portId = port.ID
		}
	}
	if portId == "" {
		log.WriteDebug("TFError| error getting specified port id from list of all ports")
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_FAILED), portName)
		return nil, nil, fmt.Errorf("The specified port name could not be found.")
	}
	provPort, err := provObj.GetPort(portId)
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragePorts provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_FAILED), portId)
		return nil, nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_AUTH_SETTINGS_BEGIN), portId)
	provPortAuthSettings, err := provObj.GetPortAuthSettings(portId)
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragePorts provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_AUTH_SETTINGS_FAILED), portId)
		return nil, nil, err
	}

	// Converting Prov to Reconciler
	reconPort := vssbmodel.StoragePort{}
	err = copier.Copy(&reconPort, provPort)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_END), portId)

	// Converting Prov to Reconciler
	reconPortAuthSetting := vssbmodel.PortAuthSettings{}
	err = copier.Copy(&reconPortAuthSetting, provPortAuthSettings)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, nil, err
	}
	log.WriteDebug("TFError|  reconPortAuthSetting : %v", reconPortAuthSetting)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_AUTH_SETTINGS_END), portId)

	return &reconPort, &reconPortAuthSetting, nil
}
