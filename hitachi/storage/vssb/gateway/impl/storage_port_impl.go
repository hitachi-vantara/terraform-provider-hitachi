package vssbstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/vssb/gateway/http"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/gateway/model"
)

// GetStoragePorts gets ports information of vssb storage
func (psm *vssbStorageManager) GetStoragePorts() (*vssbmodel.StoragePorts, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storagePorts vssbmodel.StoragePorts
	apiSuf := "objects/ports"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storagePorts)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storagePorts, nil
}

// GetStoragePorts gets port information for a specific port of vssb storage
func (psm *vssbStorageManager) GetPort(portId string) (*vssbmodel.StoragePort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storagePort vssbmodel.StoragePort
	apiSuf := fmt.Sprintf("objects/ports/%s", portId)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storagePort)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storagePort, nil
}

// GetPortAuthSettings gets the authentication settings for the compute port for the target operation.
func (psm *vssbStorageManager) GetPortAuthSettings(portId string) (*vssbmodel.PortAuthSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var pas vssbmodel.PortAuthSettings
	apiSuf := fmt.Sprintf("objects/port-auth-settings/%s", portId)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &pas)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &pas, nil
}
