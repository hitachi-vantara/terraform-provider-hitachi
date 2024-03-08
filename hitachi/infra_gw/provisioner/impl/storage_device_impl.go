package infra_gw

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/infra_gw/gateway/impl"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetStorageDevices gets storage devices information
func (psm *infraGwManager) GetStorageDevices() (*model.StorageDevices, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.GetStorageDevices()
}

// GetMTStorageDevices gets storage devices information
func (psm *infraGwManager) GetMTStorageDevices() (*[]model.MTStorageDevice, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.GetMTStorageDevices()
}

// GetStorageDevices gets storage device information
func (psm *infraGwManager) GetStorageDevice(storageId string) (*model.StorageDevice, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.GetStorageDevice(storageId)
}

// AddStorageDevice adds a storage device
func (psm *infraGwManager) AddStorageDevice(reqBody model.CreateStorageDeviceParam) (msg *string, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.AddStorageDevice(reqBody)
}

// AddMTStorageDevice adds a storage device
func (psm *infraGwManager) AddMTStorageDevice(reqBody model.CreateMTStorageDeviceParam) (msg *string, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.AddMTStorageDevice(reqBody)
}

// UpdateStorageDevice  updates a storage device
func (psm *infraGwManager) UpdateStorageDevice(storageId string, reqBody model.PatchStorageDeviceParam) (msg *string, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.UpdateStorageDevice(storageId, reqBody)
}

// DeleteStorageDevice deletes a storage device from ucp inventory
func (psm *infraGwManager) DeleteStorageDevice(storageId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	return gatewayObj.DeleteStorageDevice(storageId)
}

// DeleteStorageDeviceFromUcp deletes a storage device from a ucp system
func (psm *infraGwManager) DeleteStorageDeviceFromUcp(ucpId, storageId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	return gatewayObj.DeleteStorageDeviceFromUcp(ucpId, storageId)
}
