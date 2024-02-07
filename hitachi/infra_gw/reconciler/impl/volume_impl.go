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

	return provObj.GetVolumes(id)
}

// ReconcileVolume will reconcile and call Create/Update/delete Volume
func (psm *infraGwManager) ReconcileVolume(storageId string, createInput *model.CreateVolumeParams, volumeID *string) (*model.VolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Get GetVolumeByName
	if createInput != nil {
		volumeInfo, ok := psm.GetVolumeByName(storageId, createInput.Name)
		if !ok {
			volId, err := psm.CreateVolume(storageId, createInput)
			if err != nil {
				log.WriteDebug("TFError| error in CreateVolume call, err: %v", err)
				return nil, err
			}
			volumeID = volId
		} else {
			_, err := psm.UpdateVolume(storageId, volumeInfo.ResourceId, createInput)
			if err != nil {
				log.WriteDebug("TFError| error in UpdateVolume call, err: %v", err)
				return nil, err
			}
			volumeID = &volumeInfo.ResourceId

		}
	} else if createInput == nil && volumeID != nil {
		err := psm.DeleteVolume(storageId, *volumeID)
		if err != nil {
			log.WriteDebug("TFError| error in DeleteVolume call, err: %v", err)
			return nil, err
		}
		return nil, nil

	}

	provVolInfo, err := psm.GetVolumeByID(storageId, *volumeID)
	if err != nil {
		log.WriteDebug("TFError| error in GetVolumeByID call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("Volume here >>>>>>>>>>>>>>>>>>>> %s", provVolInfo)

	return provVolInfo, nil
}

func (psm *infraGwManager) CreateVolume(storageId string, reqBody *model.CreateVolumeParams) (*string, error) {
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

	volumeID, err := provObj.CreateVolume(storageId, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in CreateVolume call, err: %v", err)
		return nil, err
	}

	return volumeID, nil

}

func (psm *infraGwManager) UpdateVolume(storageId string, volumeId string, reqBody *model.CreateVolumeParams) (*string, error) {
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
	updateVolParams := model.UpdateVolumeParams{Name: reqBody.Name, Capacity: reqBody.Capacity, DeduplicationCompressionMode: reqBody.DeduplicationCompressionMode}

	volumeID, err := provObj.UpdateVolume(storageId, volumeId, &updateVolParams)
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

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
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

func (psm *infraGwManager) GetVolumeByName(storageId string, volumeName string) (*model.VolumeInfo, bool) {
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

	volumesInfo, err := provObj.GetVolumeByID(storageId, volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in GetVolumes call, err: %v", err)
		return nil, err
	}

	return volumesInfo, nil

}
