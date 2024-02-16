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
	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &iscsiTargets)
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

	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &iscsiTarget)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &iscsiTarget, nil
}

// CreateIscsiTarget .
func (psm *infraGwManager) CreateIscsiTarget(storageId string, reqBody model.CreateIscsiTargetParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets", storageId)

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in CreateIscsiTarget - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// AddVolumesToIscsiTarget .
func (psm *infraGwManager) AddVolumesToIscsiTarget(storageId, iscsiTargetId string, reqBody model.AddVolumesToIscsiTargetParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s/volumes", storageId, iscsiTargetId)

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in AddVolumesToIscsiTarget - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// RemoveVolumesFromIscsiTarget .
func (psm *infraGwManager) RemoveVolumesFromIscsiTarget(storageId, iscsiTargetId string, reqBody model.RemoveVolumesFromIscsiTargetParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s/volumes", storageId, iscsiTargetId)

	resourceId, err := httpmethod.DeleteCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in RemoveVolumesFromIscsiTarget - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// AddIqnInitiatorsToIscsiTarget .
func (psm *infraGwManager) AddIqnInitiatorsToIscsiTarget(storageId, iscsiTargetId string, reqBody model.AddIqnInitiatorsToIscsiTargetParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s/iqns", storageId, iscsiTargetId)

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in AddIqnInitiatorsToIscsiTarget - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// RemoveIqnInitiatorsFromIscsiTarget .
func (psm *infraGwManager) RemoveIqnInitiatorsFromIscsiTarget(storageId, iscsiTargetId string, reqBody model.RemoveIqnInitiatorsFromIscsiTargetParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s/volumes", storageId, iscsiTargetId)

	resourceId, err := httpmethod.DeleteCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in RemoveIqnInitiatorsFromIscsiTarget - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// UpdateHostMode .
func (psm *infraGwManager) UpdateHostMode(storageId, iscsiTargetId string, reqBody model.UpdateHostModeParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s/iqns", storageId, iscsiTargetId)

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateHostMode - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// SetChapUser .
func (psm *infraGwManager) SetChapUser(storageId, iscsiTargetId string, reqBody model.SetChapUserParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s/chapUser", storageId, iscsiTargetId)

	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in SetChapUser - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// UpdateChapUser .
func (psm *infraGwManager) UpdateChapUser(storageId, iscsiTargetId string, reqBody model.UpdateChapUserParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s/chapUser", storageId, iscsiTargetId)

	resourceId, err := httpmethod.PatchCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateChapUser - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// UpdateTargetIqnInIscsiTarget .
func (psm *infraGwManager) UpdateTargetIqnInIscsiTarget(storageId, iscsiTargetId string, reqBody model.UpdateTargetIqnInIscsiTargetParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s/chapUser", storageId, iscsiTargetId)

	resourceId, err := httpmethod.PatchCall(psm.setting, apiSuf, reqBody, nil)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateTargetIqnInIscsiTarget - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// DeleteIscsiTarget
func (psm *infraGwManager) DeleteIscsiTarget(storageId, iscsiTargetId string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s", storageId, iscsiTargetId)

	resourceId, err := httpmethod.DeleteCall(psm.setting, apiSuf, nil, nil)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteIscsiTarget - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}

// RemoveChapUser
func (psm *infraGwManager) RemoveChapUser(storageId, iscsiTargetId, chapUserId string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("/storage/devices/%s/iscsiTargets/%s/chapUsers/%s", storageId, iscsiTargetId, chapUserId)

	resourceId, err := httpmethod.DeleteCall(psm.setting, apiSuf, nil, nil)
	if err != nil {
		log.WriteDebug("TFError| error in RemoveChapUser - %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return resourceId, nil
}
