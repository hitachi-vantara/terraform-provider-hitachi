package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	provisonerimpl "terraform-provider-hitachi/hitachi/infra_gw/provisioner/impl"
	mc "terraform-provider-hitachi/hitachi/infra_gw/reconciler/message-catalog"
)

// GetHostGroups gets host groups information
func (psm *infraGwManager) GetHostGroups(id string, port string) (*model.HostGroups, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetHostGroups(id, port)
}

func (psm *infraGwManager) GetHostGroup(storageId, port, hostGroupName string) (*model.HostGroup, error, bool) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err, false
	}

	hgs, err := provObj.GetHostGroups(storageId, port)
	if err != nil {
		log.WriteDebug("TFError| error getting GetHostGroups, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_HOSTGROUPS_FAILED), psm.setting.Address)
		return nil, err, false
	}

	var result model.HostGroup
	success := false
	if hostGroupName != "" && port != "" {
		for _, hg := range hgs.Data {
			if hg.HostGroupName == hostGroupName && hg.Port == port {
				result.Path = hgs.Path
				result.Message = hgs.Message
				result.Data = hg
				success = true
				break
			}
		}
	}
	if success {
		return &result, nil, true
	} else {
		err := fmt.Errorf("port %s and hostgroup name %s not found", port, hostGroupName)
		return nil, err, false
	}
}

// GetHostGroup gets host group information
func (psm *infraGwManager) GetHostGroupById(id string, hostGrId string) (*model.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetHostGroup(id, hostGrId)
}

// ReconcileHostGroup will reconcile and call Create/Update hostgroup
func (psm *infraGwManager) ReconcileHostGroup(storageId string, createInput *model.CreateHostGroupParam) (*model.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var reconcilerHg *model.HostGroup = &model.HostGroup{}
	// 1) If Hostroup Exisit - Update 2) Hostgroup Not Exist - Create New
	if createInput.HostGroupName != "" {
		// Get Hostgroup
		hg, err, success := psm.GetHostGroup(storageId, createInput.Port, createInput.HostGroupName)
		if err != nil {
			log.WriteDebug("TFError| error in GetHostGroup provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_HOSTGROUP_FAILED), storageId, createInput.Port, createInput.HostGroupName)
			return nil, err
		}

		if !success {
			// Hostgroup does not exist - create new
			reconcilerHg, err = psm.createHostGroup(storageId, createInput)
			if err != nil {
				log.WriteDebug("TFError| error in createHostGroup call, err: %v", err)
				return reconcilerHg, err
			}
		} else {
			// Hostgroup already exist
			reconcilerHg, err = psm.updateHostGroup(storageId, hg.Data.ResourceId, createInput)
			if err != nil {
				log.WriteDebug("TFError| error in updateHostgroup call, err: %v", err)
				return reconcilerHg, err
			}
		}
	} else {
		// Hostgroup number not given so new hostgroup will be create
		var err error = nil
		reconcilerHg, err = psm.createHostGroup(storageId, createInput)
		if err != nil {
			log.WriteDebug("TFError| error in createHostGroup call, err: %v", err)
			return reconcilerHg, err
		}
	}

	return reconcilerHg, nil
}

// CreateHostGroup .
func (psm *infraGwManager) createHostGroup(storageId string, reqBody *model.CreateHostGroupParam) (*model.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	id, err := provObj.CreateHostGroup(storageId, *reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateHostGroup call, err: %v", err)
		return nil, err
	}

	return psm.GetHostGroupById(storageId, *id)

}

// updateHostGroup .
func (psm *infraGwManager) updateHostGroup(storageId, hostGroupId string, reqBody *model.CreateHostGroupParam) (*model.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	id, err := provObj.UpdateHostGroup(storageId, hostGroupId, *reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateHostGroup call, resource id: %v err: %v", id, err)
		return nil, err
	}

	return psm.GetHostGroupById(storageId, hostGroupId)
}
