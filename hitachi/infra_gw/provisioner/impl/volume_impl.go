package infra_gw

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/infra_gw/gateway/impl"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetVolumes gets volumes information
func (psm *infraGwManager) GetVolumes(id string) (*model.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	if psm.setting.PartnerId != nil {
		return gatewayObj.GetVolumesByPartnerSubscriberID(id)
	}

	return gatewayObj.GetVolumes(id)
}

// GetVolume by id gets volume information
func (psm *infraGwManager) GetVolumeByID(storageId string, volumeId string) (*model.VolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)

	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	volumes, err := gatewayObj.GetVolumeByID(storageId, volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in GetVolumeByID call, err: %v", err)
		return nil, err
	}

	return &volumes.Data, nil
}

func (psm *infraGwManager) GetVolumeByPartnerSubscriberID(storageId string, volumeId string) (*model.VolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)

	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	volumes, err := gatewayObj.GetVolumeByID(storageId, volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in GetVolumeByID call, err: %v", err)
		return nil, err
	}

	return &volumes.Data, nil
}

// CreateVolume created the  volume in the storage
func (psm *infraGwManager) CreateVolume(storageId string, reqBody *model.CreateVolumeParams) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	if psm.setting.PartnerId == nil {
		return gatewayObj.CreateVolume(storageId, reqBody)
	}
	return gatewayObj.CreateMTVolume(storageId, reqBody)

}

// UpdateVolume update the volume in the storage

func (psm *infraGwManager) UpdateVolume(storageId string, volumeId string, reqBody *model.UpdateVolumeParams) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	gateSetting := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.UpdateVolume(storageId, volumeId, reqBody)
}

// DeleteVolume deletes a volume from the storage
func (psm *infraGwManager) DeleteVolume(storageId string, volumeId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	gateSetting := model.InfraGwSettings(psm.setting)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	if psm.setting.SubscriberId != nil {
		return gatewayObj.DeleteMTVolume(storageId, volumeId)
	}

	return gatewayObj.DeleteVolume(storageId, volumeId)
}
