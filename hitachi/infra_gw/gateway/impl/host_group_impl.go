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
	err := httpmethod.GetCall(psm.setting, apiSuf, &hostGroups)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &hostGroups, nil
}

// GetHostGroup gets hostGroup information
func (psm *infraGwManager) GetHostGroup(storageId string, hostGrId string) (*model.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hostGroup model.HostGroup

	apiSuf := fmt.Sprintf("/storage/devices/%s/hostGroups/%s", storageId, hostGrId)

	err := httpmethod.GetCall(psm.setting, apiSuf, &hostGroup)
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

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateHostGroup - %s API call, err: %v", apiSuf, err)
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

	resourceId, err := httpmethod.PatchCall(psm.setting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateHostGroup - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}
