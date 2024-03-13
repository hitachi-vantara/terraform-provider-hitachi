package infra_gw

import (
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetHostGroups gets hostGroups information
func (psm *infraGwManager) GetHostGroups(storageId string, port string) (*model.HostGroups, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hostGroups model.HostGroups
	var apiSuf string

	if port == "" {
		apiSuf = fmt.Sprintf("/storage/devices/%s/hostGroups", storageId)
	} else {
		apiSuf = fmt.Sprintf("/storage/devices/%s/hostGroups?port=%s", storageId, port)
	}
	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &hostGroups)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &hostGroups, nil
}

// Get GetHostGroupsByPartnerIdOrSubscriberID
func (psm *infraGwManager) GetHostGroupsByPartnerIdOrSubscriberID(storageId string) (*model.MTHostGroups, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.MTHostGroups

	psm.setting.V3API = true

	headers := map[string]string{
		"partnerId": *psm.setting.PartnerId,
	}
	if psm.setting.SubscriberId != nil {
		headers["subscriberId"] = *psm.setting.SubscriberId
	}

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups", storageId)

	err := httpmethod.GetCall(psm.setting, apiSuf, &headers, &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &response, nil
}

// GetHostGroup gets hostGroup information
func (psm *infraGwManager) GetHostGroup(storageId string, hostGrId string) (*model.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hostGroup model.HostGroup

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups/%s", storageId, hostGrId)

	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &hostGroup)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &hostGroup, nil
}

// CreateHostGroup .
func (psm *infraGwManager) CreateHostGroup(storageId string, reqBody model.CreateHostGroupParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups", storageId)

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in CreateHostGroup - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// CreateMTHostGroup .
func (psm *infraGwManager) CreateMTHostGroup(storageId string, reqBody model.CreateHostGroupParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	headers := map[string]string{
		"partnerId":    *psm.setting.PartnerId,
		"subscriberId": *psm.setting.SubscriberId,
	}

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups", storageId)

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, &headers)
	if err != nil {
		log.WriteDebug("TFError| error in CreateMTHostGroup - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// UpdateHostGroup .
func (psm *infraGwManager) UpdateHostGroup(storageId, hostGroupId string, reqBody model.PatcheHostGroupParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups/%s", storageId, hostGroupId)

	resourceId, err := httpmethod.PatchCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateHostGroup - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// AddVolumesToHostGroup .
func (psm *infraGwManager) AddVolumesToHostGroup(storageId, hostGroupId string, reqBody model.AddVolumesToHostGroupParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups/%s/volumes", storageId, hostGroupId)

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in AddVolumesToHostGroup - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// AddVolumesToHostGroupToSubscriber .
func (psm *infraGwManager) AddVolumesToHostGroupToSubscriber(storageId, hostGroupId string, reqBody model.AddVolumesToHostGroupParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	headers := map[string]string{
		"subscriberId": *psm.setting.SubscriberId,
	}

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups/%s/volumes", storageId, hostGroupId)

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, &headers)
	if err != nil {
		log.WriteDebug("TFError| error in AddVolumesToHostGroupToSubscriber - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// DeleteVolumesFromHostGroup .
func (psm *infraGwManager) DeleteVolumesFromHostGroup(storageId, hostGroupId string, reqBody model.DeleteVolumesToHostGroupParam) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups/%s/volumes", storageId, hostGroupId)

	_, err := httpmethod.DeleteCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteVolumesFromHostGroup - %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// DeleteVolumesFromHostGroupFromSubscriber .
func (psm *infraGwManager) DeleteVolumesFromHostGroupFromSubscriber(storageId, hostGroupId string, reqBody model.DeleteVolumesToHostGroupParam) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	headers := map[string]string{
		"subscriberId": *psm.setting.SubscriberId,
	}

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups/%s/volumes", storageId, hostGroupId)

	_, err := httpmethod.DeleteCall(psm.setting, apiSuf, reqBody, &headers)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteVolumesFromHostGroupFromSubscriber - %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// DeleteHostGroup .
func (psm *infraGwManager) DeleteHostGroup(storageId, hostGroupId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups/%s", storageId, hostGroupId)

	_, err := httpmethod.DeleteCall(psm.setting, apiSuf, nil, nil)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteHostGroup - %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// DeleteMTHostGroup .
func (psm *infraGwManager) DeleteMTHostGroup(storageId, hostGroupId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	headers := map[string]string{
		"subscriberId": *psm.setting.SubscriberId,
	}

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups/%s?isDelete=true", storageId, hostGroupId)

	_, err := httpmethod.DeleteCall(psm.setting, apiSuf, nil, &headers)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteMTHostGroup - %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// Get GetHostGroupsDetailsByPartnerIdOrSubscriberID
func (psm *infraGwManager) GetHostGroupsDetailsByPartnerIdOrSubscriberID(storageId string) (*model.MTHostGroups, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.MTHostGroups

	psm.setting.V3API = true

	headers := map[string]string{
		"partnerId": *psm.setting.PartnerId,
	}
	if psm.setting.SubscriberId != nil {
		headers["subscriberId"] = *psm.setting.SubscriberId
	}

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups/details", storageId)

	err := httpmethod.GetCall(psm.setting, apiSuf, &headers, &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &response, nil
}
