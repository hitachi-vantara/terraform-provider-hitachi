package infra_gw

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/infra_gw/gateway/impl"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	mc "terraform-provider-hitachi/hitachi/infra_gw/reconciler/message-catalog"
)

// GetVolumes gets volumes information
func (psm *infraGwManager) GetVolumes(id string) (*model.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_INFRA_GET_ALL_VOLUMES_BEGIN))

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_GET_ALL_VOLUMES_FAILED))

		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.GetVolumes(id)
}

// GetVolumesFromLdevIds gets volumes information
func (psm *infraGwManager) GetVolumesFromLdevIds(id string, fromLdevId int, toLdevId int) (*model.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_INFRA_GET_VOLUME_BEGIN), id)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_GET_VOLUME_FAILED), id)

		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return gatewayObj.GetVolumesFromLdevIds(id, fromLdevId, toLdevId)
}

func (psm *infraGwManager) GetVolumesByPartnerSubscriberID(storageId string, fromLdevId *int, toLdevId *int) (*model.MTVolumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_INFRA_GET_VOLUME_BEGIN), storageId)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_GET_VOLUME_FAILED), storageId)

		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	if fromLdevId == nil {
		temp := 0
		fromLdevId = &temp
	}

	if toLdevId == nil {
		temp := 65000
		toLdevId = &temp
	}
	return gatewayObj.GetVolumesDetailsByPartnerSubscriberID(storageId, *fromLdevId, *toLdevId)
}

// GetVolume by id gets volume information
func (psm *infraGwManager) GetVolumeByID(storageId string, volumeId string) (*model.VolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	gateSetting := model.InfraGwSettings(psm.setting)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_INFRA_GET_VOLUME_BEGIN), volumeId)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)

	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	volumes, err := gatewayObj.GetVolumeByID(storageId, volumeId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_GET_VOLUME_FAILED), volumeId)

		log.WriteDebug("TFError| error in GetVolumeByID call, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_INFRA_GET_VOLUME_END), volumeId)

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
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_INFRA_CREATE_VOLUME_BEGIN))

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		log.WriteInfo(mc.GetMessage(mc.ERR_GET_INFRA_CREATE_VOLUME_FAILED))

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
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_INFRA_UPDATE_VOLUME_BEGIN), storageId)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_UPDATE_VOLUME_FAILED), storageId)

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
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_INFRA_DELETE_VOLUME_BEGIN), storageId)

	gatewayObj, err := gatewayimpl.NewEx(gateSetting)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_DELETE_VOLUME_FAILED), storageId)

		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	if psm.setting.SubscriberId != nil {
		return gatewayObj.DeleteMTVolume(storageId, volumeId)
	}

	return gatewayObj.DeleteVolume(storageId, volumeId)
}
