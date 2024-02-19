package infra_gw

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	provisonerimpl "terraform-provider-hitachi/hitachi/infra_gw/provisioner/impl"
)

// GetVolumes gets volumes information
func (psm *infraGwManager) GetPartnerAndSubscriberId(userName string) (*model.MTDetails, error) {
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
	var mtDetails model.MTDetails

	status, partnerId, subscriberId, err := provObj.GetPartnerAndSubscriberId(userName)

	if err != nil {
		log.WriteDebug("TFError| error in GetPartnerAndSubscriberId call, err: %v", err)
		return nil, err
	}

	if status {
		mtDetails.PartnerId = partnerId
		mtDetails.SubscriberId = subscriberId
	}
	return &mtDetails, nil
}
