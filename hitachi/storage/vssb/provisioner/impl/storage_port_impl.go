package vssbstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vssb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vssb/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/message-catelog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetStoragePorts gets ports information of vssb storage
func (psm *vssbStorageManager) GetStoragePorts() (*vssbmodel.StoragePorts, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_PORTS_BEGIN))
	storagePorts, err := gatewayObj.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetStoragePorts, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_PORTS_FAILED))
		return nil, err
	}

	provStoragePorts := vssbmodel.StoragePorts{}
	err = copier.Copy(&provStoragePorts, storagePorts)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_PORTS_END))

	return &provStoragePorts, nil
}

// GetStoragePort gets port information for a specific port of vssb storage
func (psm *vssbStorageManager) GetPort(portId string) (*vssbmodel.StoragePort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_BEGIN), portId)
	storagePort, err := gatewayObj.GetPort(portId)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetPort, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_FAILED), portId)
		return nil, err
	}

	provStoragePort := vssbmodel.StoragePort{}
	err = copier.Copy(&provStoragePort, storagePort)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_END), portId)

	return &provStoragePort, nil
}

// GetPortAuthSettings gets the authentication settings for the compute port for the target operation.
func (psm *vssbStorageManager) GetPortAuthSettings(portId string) (*vssbmodel.PortAuthSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_AUTH_SETTINGS_BEGIN), portId)
	portAuthSettings, err := gatewayObj.GetPortAuthSettings(portId)

	log.WriteDebug("TFError| gateway portAuthSettings: %+v", portAuthSettings)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetPortAuthSettings, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_AUTH_SETTINGS_FAILED), portId)
		return nil, err
	}

	provPortAuthSettings := vssbmodel.PortAuthSettings{}
	err = copier.Copy(&provPortAuthSettings, portAuthSettings)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteDebug("TFError| provisioner portAuthSettings: %+v", provPortAuthSettings)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_AUTH_SETTINGS_END), portId)

	return &provPortAuthSettings, nil
}
