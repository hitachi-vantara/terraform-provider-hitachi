package infra_gw

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/infra_gw/gateway/impl"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetIscsiTargets gets IscsiTargets information
func (psm *infraGwManager) GetIscsiTargets(id string, port string) (*model.IscsiTargets, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.GetIscsiTargets(id, port)
}

// GetIscsiTarget gets IscsiTarget information
func (psm *infraGwManager) GetIscsiTarget(id string, iscsiTargetId string) (*model.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.GetIscsiTarget(id, iscsiTargetId)
}
