package infra_gw

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	provisonerimpl "terraform-provider-hitachi/hitachi/infra_gw/provisioner/impl"
)

// GetVolumes gets volumes information
func (psm *infraGwManager) GetVolumes(id string) (*model.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provSetting := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(provSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetVolumes(id)
}

func (psm *infraGwManager) GetVolumesFromLdevIds(id string, fromLdevId *int, toLdevId *int) (*model.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provSetting := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(provSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	defualtId := 0

	if fromLdevId == nil {
		fromLdevId = &defualtId
	}
	if toLdevId == nil {
		toLdevId = &defualtId
	}

	return provObj.GetVolumesFromLdevIds(id, *fromLdevId, *toLdevId)
}

func (psm *infraGwManager) GetVolumesByPartnerSubscriberID(id string, fromLdevId int, toLdevId int) (*model.MTVolumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provSetting := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(provSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetVolumesByPartnerSubscriberID(id, &fromLdevId, &toLdevId)
}

// ReconcileVolume will reconcile and call Create/Update/delete Volume
func (psm *infraGwManager) ReconcileVolume(storageId string, createInput *model.CreateVolumeParams, volumeID *string) (*model.VolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Get GetVolumeByName
	if createInput != nil {
		_, ok := psm.GetVolumeByName(storageId, createInput.Name)
		if !ok && volumeID == nil {
			_, err := psm.CreateVolume(storageId, createInput)
			if err != nil {
				log.WriteDebug("TFError| error in CreateVolume call, err: %v", err)
				return nil, err
			}
		} else {
			_, err := psm.UpdateVolume(storageId, volumeID, createInput)
			if err != nil {
				log.WriteDebug("TFError| error in UpdateVolume call, err: %v", err)
				return nil, err
			}

		}
	} else if createInput == nil && volumeID != nil {
		err := psm.DeleteVolume(storageId, *volumeID)
		if err != nil {
			log.WriteDebug("TFError| error in DeleteVolume call, err: %v", err)
			return nil, err
		}
		return nil, nil

	}

	volumeInfo, _ := psm.GetVolumeByName(storageId, createInput.Name)
	return volumeInfo, nil
}

func (psm *infraGwManager) CreateVolume(storageId string, reqBody *model.CreateVolumeParams) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	provSetting := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(provSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	volumeID, err := provObj.CreateVolume(storageId, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateVolume call, err: %v", err)
		return nil, err
	}

	return volumeID, nil

}

func (psm *infraGwManager) UpdateVolume(storageId string, volumeId *string, reqBody *model.CreateVolumeParams) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provSetting := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(provSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	updateVolParams := model.UpdateVolumeParams{}

	if reqBody.Name != "" {

		updateVolParams.Name = reqBody.Name
	}
	if reqBody.DeduplicationCompressionMode != "" {
		updateVolParams.DeduplicationCompressionMode = reqBody.DeduplicationCompressionMode
	}
	if reqBody.Capacity != "" {
		updateVolParams.Capacity = reqBody.Capacity
	}

	volumeID, err := provObj.UpdateVolume(storageId, *volumeId, &updateVolParams)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateVolume call, err: %v", err)
		return nil, err
	}

	return volumeID, nil

}

func (psm *infraGwManager) DeleteVolume(storageId string, volumeId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provSetting := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(provSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	err = provObj.DeleteVolume(storageId, volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteVolume call, err: %v", err)
		return err
	}

	return nil

}

func (psm *infraGwManager) GetMTVolumeByName(storageId string, volumeName string) (*model.VolumeInfo, bool) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provSetting := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(provSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, false
	}

	volumes, err := provObj.GetVolumes(storageId)
	if err != nil {
		log.WriteDebug("TFError| error in GetVolumes call, err: %v", err)
		return nil, false
	}

	var ProvVolume *model.VolumeInfo
	status := false
	for _, volume := range volumes.Data {
		if volume.Name == volumeName {
			status = true
			ProvVolume = &volume
			break
		}
	}
	return ProvVolume, status

}

func (psm *infraGwManager) GetVolumeByName(storageId string, volumeName string) (*model.VolumeInfo, bool) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provSetting := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(provSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, false
	}

	volumes, err := provObj.GetVolumes(storageId)
	if err != nil {
		log.WriteDebug("TFError| error in GetVolumes call, err: %v", err)
		return nil, false
	}

	var ProvVolume *model.VolumeInfo
	status := false
	for _, volume := range volumes.Data {
		if volume.Name == volumeName {
			status = true
			ProvVolume = &volume
			break
		}
	}
	return ProvVolume, status

}

func (psm *infraGwManager) GetVolumeByID(storageId string, volumeId string) (*model.VolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provSetting := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(provSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	volumesInfo, err := provObj.GetVolumeByID(storageId, volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in GetVolumes call, err: %v", err)
		return nil, err
	}

	return volumesInfo, nil

}
