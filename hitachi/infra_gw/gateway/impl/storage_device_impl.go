package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetStorageDevices gets storage devices information
func (psm *infraGwManager) GetStorageDevices() (*model.StorageDevices, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storageDevices model.StorageDevices

	apiSuf := "/storage/devices"
	err := httpmethod.GetCall(psm.setting, apiSuf, &storageDevices)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storageDevices, nil
}

// GetStorageDevices gets storage device information
func (psm *infraGwManager) GetStorageDevice(storageId string) (*model.StorageDevice, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storageDevice model.StorageDevice

	apiSuf := fmt.Sprintf("/storage/devices/%s", storageId)
	err := httpmethod.GetCall(psm.setting, apiSuf, &storageDevice)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storageDevice, nil
}

// AddStorageDevice adds storage device to a ucp system
func (psm *infraGwManager) AddStorageDevice(storageId string, reqBody model.CreateStorageDeviceParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "/storage/devices"
	ret, err := httpmethod.PostCall(psm.setting, apiSuf, &reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return ret, nil
}

// UpdateStorageDevice updates storage device to a ucp system
func (psm *infraGwManager) UpdateStorageDevice(storageId string, reqBody model.PatchStorageDeviceParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s", storageId)
	ret, err := httpmethod.PatchCall(psm.setting, apiSuf, &reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return ret, nil
}
