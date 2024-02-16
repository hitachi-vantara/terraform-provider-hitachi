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
	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &storageDevices)
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
	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &storageDevice)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storageDevice, nil
}

// AddStorageDevice adds storage device to a ucp system
func (psm *infraGwManager) AddStorageDevice(reqBody model.CreateStorageDeviceParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "/storage/devices"
	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, &reqBody, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return resourceId, nil
}

// AddMTStorageDevice adds storage device to a ucp system multi-tenancy
func (psm *infraGwManager) AddMTStorageDevice(reqBody model.CreateStorageDeviceParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	psm.setting.V3API = true

	headers := map[string]string{
		"partnerId": *psm.setting.PartnerId,
	}

	apiSuf := "/storage/devices"
	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, &reqBody, &headers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return resourceId, nil
}

func (psm *infraGwManager) AddStorageDeviceToPartner(reqBody *model.StorageDeviceToPartnerReq) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "/storage"
	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, &reqBody, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return resourceId, nil
}

// UpdateStorageDevice updates storage device to a ucp system
func (psm *infraGwManager) UpdateStorageDevice(storageId string, reqBody model.PatchStorageDeviceParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s", storageId)
	resourceId, err := httpmethod.PatchCall(psm.setting, apiSuf, &reqBody, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return resourceId, nil
}
