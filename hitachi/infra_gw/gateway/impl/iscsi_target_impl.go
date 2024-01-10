package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetIscsiTargets gets IscsiTargets information
func (psm *infraGwManager) GetIscsiTargets(id string, port string) (*model.IscsiTargets, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var iscsiTargets model.IscsiTargets
	var apiSuf string

	if port == "" {
		apiSuf = fmt.Sprintf("/storage/devices/%s/iscsiTargets", id)
	} else {
		apiSuf = fmt.Sprintf("/storage/devices/%s/iscsiTargets?port=%s", id, port)
	}
	err := httpmethod.GetCall(psm.setting, apiSuf, &iscsiTargets)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &iscsiTargets, nil
}

// GetIscsiTarget gets IscsiTarget information
func (psm *infraGwManager) GetIscsiTarget(id string, iscsiTargetId string) (*model.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var iscsiTarget model.IscsiTarget

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s", id, iscsiTargetId)

	err := httpmethod.GetCall(psm.setting, apiSuf, &iscsiTarget)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &iscsiTarget, nil
}
