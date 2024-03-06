package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	provisonerimpl "terraform-provider-hitachi/hitachi/infra_gw/provisioner/impl"
	mc "terraform-provider-hitachi/hitachi/infra_gw/reconciler/message-catalog"
)

// GetIscsiTargets gets IscsiTargets information
func (psm *infraGwManager) GetIscsiTargets(id string, port string) (*model.IscsiTargets, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetIscsiTargets(id, port)
}

func (psm *infraGwManager) GetIscsiTarget(storageId, port, iscsiName string) (*model.IscsiTarget, error, bool) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err, false
	}

	iscsiTargets, err := provObj.GetIscsiTargets(storageId, port)
	if err != nil {
		log.WriteDebug("TFError| error getting GetHostGroups, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_ISCSI_TARGETS_FAILED), psm.setting.Address)
		return nil, err, false
	}

	var result model.IscsiTarget
	success := false
	if iscsiName != "" && port != "" {
		for _, iscsi := range iscsiTargets.Data {
			if iscsi.ISCSIName == iscsiName && iscsi.PortId == port {
				result.Path = iscsiTargets.Path
				result.Message = iscsiTargets.Message
				result.Data = iscsi
				success = true
				break
			}
		}
	}
	if success {
		return &result, nil, true
	} else {
		err := fmt.Errorf("port %s and hostgroup name %s not found", port, iscsiName)
		return nil, err, false
	}
}

// GetIscsiTarget gets IscsiTarget information
func (psm *infraGwManager) GetIscsiTargetById(id string, iscsiTargetId string) (*model.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetIscsiTarget(id, iscsiTargetId)
}

// ReconcileIscsiTarget will reconcile and call Create/Update iscsi target
func (psm *infraGwManager) ReconcileIscsiTarget(storageId string, createInput *model.CreateIscsiTargetParam) (*model.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var reconcilerIscsiTarget *model.IscsiTarget = &model.IscsiTarget{}

	if createInput.IscsiName == "" {
		// Hostgroup name not given so throw err
		err := fmt.Errorf("%s", "iscsi_target_alias Name empty")
		return reconcilerIscsiTarget, err
	}

	// Get Hostgroup
	iscsiTarget, err, success := psm.GetIscsiTarget(storageId, createInput.Port, createInput.IscsiName)
	if err != nil {
		log.WriteDebug("TFError| error in GetHostGroup provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_ISCSI_TARGET_FAILED), storageId, createInput.Port, createInput.IscsiName)
		return nil, err
	}

	// 1) If Iscsi Target Exists - Update 2) Iscsi Target does not Exist - Create New
	if !success {
		// Hostgroup does not exist - create new
		reconcilerIscsiTarget, err = psm.createIscsiTarget(storageId, createInput)
		if err != nil {
			log.WriteDebug("TFError| error in createHostGroup call, err: %v", err)
			return nil, err
		}
	} else {
		// Hostgroup already exist
		reconcilerIscsiTarget, err = psm.updateIscsiTarget(storageId, iscsiTarget.Data.ResourceId, createInput)
		if err != nil {
			log.WriteDebug("TFError| error in updateHostgroup call, err: %v", err)
			return nil, err
		}
	}

	return reconcilerIscsiTarget, nil
}

// updateIscsiTarget .
func (psm *infraGwManager) createIscsiTarget(storageId string, reqBody *model.CreateIscsiTargetParam) (*model.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	id, err := provObj.CreateIscsiTarget(storageId, *reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateHostGroup call, err: %v", err)
		return nil, err
	}

	return psm.GetIscsiTargetById(storageId, *id)

}

// updateIscsiTarget .
func (psm *infraGwManager) updateIscsiTarget(storageId, hostGroupId string, reqBody *model.CreateIscsiTargetParam) (*model.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	id, err := provObj.UpdateIscsiTarget(storageId, hostGroupId, *reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateHostGroup call, resource id: %v err: %v", id, err)
		return nil, err
	}

	return psm.GetIscsiTargetById(storageId, hostGroupId)
}
