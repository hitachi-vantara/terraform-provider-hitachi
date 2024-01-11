package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetUcpSystems gets all UCP Systems information
func (psm *infraGwManager) GetUcpSystems() (*model.UcpSystems, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var ucpSystems model.UcpSystems

	apiSuf := "/systems"
	err := httpmethod.GetCall(psm.setting, apiSuf, &ucpSystems)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &ucpSystems, nil
}

// GetUcpSystemById gets UCP System information By Id
func (psm *infraGwManager) GetUcpSystemById(id string) (*model.UcpSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var ucpSystem model.UcpSystem

	apiSuf := fmt.Sprintf("/systems/%s", id)
	err := httpmethod.GetCall(psm.setting, apiSuf, &ucpSystem)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &ucpSystem, nil
}

// CreateUcpSystem
func (psm *infraGwManager) CreateUcpSystem(reqBody model.CreateUcpSystemParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "/systems"

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateUcpSystem - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}
