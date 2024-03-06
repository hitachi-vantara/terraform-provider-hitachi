package infra_gw

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	provisonerimpl "terraform-provider-hitachi/hitachi/infra_gw/provisioner/impl"
)

// GetStoragePorts gets ports information
func (psm *infraGwManager) GetStoragePorts(id string) (*model.StoragePorts, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetStoragePorts(id)
}

// GetStoragePortsByPartnerIdOrSubscriberId gets ports information
func (psm *infraGwManager) GetStoragePortsByPartnerIdOrSubscriberId(id string) (*model.MTPorts, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetStoragePortsByPartnerIdOrSubscriberId(id)
}
