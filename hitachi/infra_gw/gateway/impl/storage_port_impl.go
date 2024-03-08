package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetStoragePorts gets ports information
func (psm *infraGwManager) GetStoragePorts(id string) (*model.StoragePorts, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storagePorts model.StoragePorts

	apiSuf := fmt.Sprintf("/storage/devices/%s/ports", id)
	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &storagePorts)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storagePorts, nil
}

// GetStoragePortsByPartnerIdOrSubscriberId  gets all StoragePorts by subscriberId/PartnerId.
func (psm *infraGwManager) GetStoragePortsByPartnerIdOrSubscriberId(id string) (*model.MTPorts, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	headers := map[string]string{
		"partnerId": *psm.setting.PartnerId,
	}

	if psm.setting.SubscriberId != nil {
		headers["subscriberId"] = *psm.setting.SubscriberId
	}

	var storagePorts model.MTPorts

	apiSuf := fmt.Sprintf("/storage/%s/ports", id)
	err := httpmethod.GetCall(psm.setting, apiSuf, &headers, &storagePorts)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storagePorts, nil
}

/*

// GetStoragePorts gets port information for a specific port of vssb storage
func (psm *infraGwManager) GetPort(portId string) (*model.StoragePort, error) {
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
*/
