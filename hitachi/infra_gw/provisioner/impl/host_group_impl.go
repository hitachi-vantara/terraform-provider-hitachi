package infra_gw

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/infra_gw/gateway/impl"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetHostGroups gets host groups information
func (psm *infraGwManager) GetHostGroups(storageId string, port string) (*model.HostGroups, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.GetHostGroups(storageId, port)
}

// GetHostGroupsByPartnerIdOrSubscriberID gets host groups information
func (psm *infraGwManager) GetHostGroupsByPartnerIdOrSubscriberID(storageId string) (*model.MTHostGroups, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.GetHostGroupsByPartnerIdOrSubscriberID(storageId)
}

// GetHostGroup gets host group information
func (psm *infraGwManager) GetHostGroup(storageId string, hostGrId string) (*model.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.GetHostGroup(storageId, hostGrId)
}

// CreateHostGroup .
func (psm *infraGwManager) CreateHostGroup(storageId string, reqBody model.CreateHostGroupParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.CreateHostGroup(storageId, reqBody)
}

// CreateMTHostGroup .
func (psm *infraGwManager) CreateMTHostGroup(storageId string, reqBody model.CreateHostGroupParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.CreateMTHostGroup(storageId, reqBody)
}

// UpdateHostGroup .
func (psm *infraGwManager) UpdateHostGroup(storageId, hostGroupId string, reqBody model.CreateHostGroupParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	updateReqBody := model.PatcheHostGroupParam{
		HostMode:        reqBody.HostMode,
		HostModeOptions: reqBody.HostModeOptions,
	}

	return gatewayObj.UpdateHostGroup(storageId, hostGroupId, updateReqBody)
}

// DeleteHostGroup .
func (psm *infraGwManager) DeleteHostGroup(storageId, hgId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	return gatewayObj.DeleteHostGroup(storageId, hgId)
}

// DeleteHostGroup .
func (psm *infraGwManager) DeleteMTHostGroup(storageId, hgId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	return gatewayObj.DeleteMTHostGroup(storageId, hgId)
}
