package vssbstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vssb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vssb/gateway/model"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/model"
)

// DeleteAllChapUsersFromComputePort deletes all chap users from a specific compute port of vssb storage
func (psm *vssbStorageManager) DeleteAllChapUsersFromComputePort(portId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	currentChapUsers, err := gatewayObj.GetChapUsersAllowedToAccessPort(portId)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetChapUsersAllowedToAccessPort, err: %+v", err)
		return err
	}

	for _, cu := range currentChapUsers.Data {
		err = gatewayObj.DeletePortAccessForChapUser(portId, cu.ID)

		// TODO if fails retry a couple of times before returning error
		if err != nil {
			log.WriteDebug("TFError| failed to call DeletePortAccessForChapUser, err: %+v , chap user : %+v", err, cu.TargetChapUserName)
			return err
		}
	}

	return nil
}

// UpdatePortAuthSettings updates the authmode of a specific compute port of vssb storage
func (psm *vssbStorageManager) UpdatePortAuthSettings(portId string, reqBody *vssbmodel.PortAuthSettings) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	createReq := vssbgatewaymodel.PortAuthSettings{
		AuthMode:            reqBody.AuthMode,
		IsDiscoveryChapAuth: reqBody.IsDiscoveryChapAuth,
		IsMutualChapAuth:    reqBody.IsMutualChapAuth,
	}

	err = gatewayObj.UpdatePortAuthSettings(portId, &createReq)
	if err != nil {
		log.WriteDebug("TFError| failed to call UpdatePortAuthSettings, err: %+v", err)
		return err
	}

	log.WriteInfo("UpdatePortAuthSettings Successful")

	return nil
}

// AddChapUsersToComputePort add chap users to a specific compute port of vssb storage
func (psm *vssbStorageManager) AddChapUsersToComputePort(portId string, chapUserIds []string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	for _, cuId := range chapUserIds {
		req := vssbgatewaymodel.ChapUserIdReq{
			ChapUserId: cuId,
		}
		err = gatewayObj.AllowChapUserToAccessPort(portId, &req)

		// TODO if fails retry a couple of times before returning error
		if err != nil {
			log.WriteDebug("TFError| failed to call AllowChapUserToAccessPort, err: %+v , chap user Id: %+v", err, cuId)
			return err
		}
	}

	return nil
}
