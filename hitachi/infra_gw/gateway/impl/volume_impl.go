package infra_gw

import (
	"encoding/json"
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetVolumes gets volumes information
func (psm *infraGwManager) GetVolumes(id string) (*model.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.Volumes

	apiSuf := fmt.Sprintf("/storage/devices/%s/volumes", id)
	log.WriteDebug(apiSuf)
	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &response, nil
}

// GetMTVolumes gets volumes information from multi-tenancy
func (psm *infraGwManager) GetMTVolumes(storageid string) (*model.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	var response model.Volumes
	headers := map[string]string{
		"partnerId":    *psm.setting.PartnerId,
		"subscriberId": *psm.setting.SubscriberId,
	}

	apiSuf := fmt.Sprintf("/storage/%s/volume", storageid)
	log.WriteDebug(apiSuf)
	err := httpmethod.GetCall(psm.setting, apiSuf, &headers, &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &response, nil
}

// GetVolumes gets volumes information
func (psm *infraGwManager) GetVolumeByID(storageId string, volumeID string) (*model.Volume, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.Volume

	psm.setting.V3API = true

	headers := map[string]string{
		"partnerId": *psm.setting.PartnerId,
	}
	if psm.setting.SubscriberId != nil {
		headers["subscriberId"] = *psm.setting.SubscriberId
	}
	apiSuf := fmt.Sprintf("/storage/devices/%s/volumes/%s", storageId, volumeID)

	err := httpmethod.GetCall(psm.setting, apiSuf, &headers, &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	jsonDataBefore, err := json.Marshal(response)
	if err != nil {
		log.WriteDebug("Error marshaling to JSON:", err)

	}

	log.WriteDebug("in gateway side in in get volume by id  >>>>>>>>>>: %s", string(jsonDataBefore))

	return &response, nil
}

// Get GetVolumeByPartnerSubscriberID
func (psm *infraGwManager) GetVolumesByPartnerSubscriberID(storageId string) (*model.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.Volumes

	apiSuf := fmt.Sprintf("/storage/%s/volume", storageId)

	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	jsonDataBefore, err := json.Marshal(response)
	if err != nil {
		log.WriteDebug("Error marshaling to JSON:", err)

	}

	log.WriteDebug("in gateway side in in get volume by id  >>>>>>>>>>: %s", string(jsonDataBefore))

	return &response, nil
}

func (psm *infraGwManager) CreateVolume(storageId string, reqBody *model.CreateVolumeParams) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/volumes", storageId)
	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in CreateVolume - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

func (psm *infraGwManager) CreateMTVolume(storageId string, reqBody *model.CreateVolumeParams) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	headers := map[string]string{
		"partnerId":    *psm.setting.PartnerId,
		"subscriberId": *psm.setting.SubscriberId,
	}

	apiSuf := fmt.Sprintf("/storage/%s/volumes", storageId)
	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, &headers)
	if err != nil {
		log.WriteDebug("TFError| error in CreateVolume - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

func (psm *infraGwManager) UpdateVolume(storageId string, volumeID string, reqBody *model.UpdateVolumeParams) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/volumes/%s", storageId, volumeID)

	resourceId, err := httpmethod.PatchCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateVolume - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

func (psm *infraGwManager) DeleteVolume(storageId string, volumeID string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/volumes/%s", storageId, volumeID)

	_, err := httpmethod.DeleteCall(psm.setting, apiSuf, nil, nil)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteVolume - %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

func (psm *infraGwManager) DeleteMTVolume(storageId string, volumeID string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	headers := map[string]string{
		"subscriberId": *psm.setting.SubscriberId,
	}
	apiSuf := fmt.Sprintf("/storage/%s/volumes/%s?isDelete=true", storageId, volumeID)

	_, err := httpmethod.DeleteCall(psm.setting, apiSuf, nil, &headers)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteVolume - %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}
