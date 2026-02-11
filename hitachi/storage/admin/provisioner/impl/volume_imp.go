package admin

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymanager "terraform-provider-hitachi/hitachi/storage/admin/gateway"
	gwyimpl "terraform-provider-hitachi/hitachi/storage/admin/gateway/impl"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/provisioner/message-catalog"
)

// thin pass-through function to call gateway
// no separate model (uses gateway), but has message catalog logging
func (psm *adminStorageManager) GetVolumes(queryParams gwymodel.GetVolumeParams) (*gwymodel.VolumeInfoList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUMES_BEGIN), psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	log.WriteDebug("TFDebug| QueryParams:%+v\n", queryParams)

	volumeInfoList, err := gatewayObj.GetVolumes(queryParams)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_VOLUMES_FAILED), psm.storageSetting.Serial)
		return nil, err
	}

	// log.WriteDebug("TFDebug| VolumeInfoList: %+v\n", volumeInfoList)

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUMES_END), psm.storageSetting.Serial)
	return volumeInfoList, nil
}

// thin pass-through function to call gateway via GetVolumes
func (psm *adminStorageManager) GetVolumeByID(volumeID int) (*gwymodel.VolumeInfoByID, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUME_BY_ID_BEGIN), volumeID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return nil, err
	}

	volInfo, err := gatewayObj.GetVolumeByID(volumeID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_VOLUME_BY_ID_FAILED), volumeID, psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteDebug("TFDebug| VolumeInfoByID: %+v\n", volInfo)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUME_BY_ID_END), volumeID, psm.storageSetting.Serial)

	return volInfo, nil
}

// thin pass-through function to call gateway via CreateVolume
func (psm *adminStorageManager) CreateVolume(params gwymodel.CreateVolumeParams) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_VOLUMES_BEGIN), *params.Number, params.PoolID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return "", err
	}

	log.WriteDebug("TFDebug| Params:%+v\n", params)

	volumeIDs, err := gatewayObj.CreateVolume(params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUMES_FAILED), *params.Number, params.PoolID, psm.storageSetting.Serial)
		return "", err
	}

	log.WriteDebug("TFDebug| Created volume ID: %v\n", volumeIDs)
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_VOLUMES_END), *params.Number, params.PoolID, psm.storageSetting.Serial)

	return volumeIDs, nil
}

// thin pass-through function to call gateway via DeleteVolume
func (psm *adminStorageManager) DeleteVolume(volumeID int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_VOLUME_BEGIN), volumeID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	err = gatewayObj.DeleteVolume(volumeID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_VOLUME_FAILED), volumeID, psm.storageSetting.Serial)
		return err
	}

	log.WriteDebug("TFDebug| Deleted volume ID: %d\n", volumeID)
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_VOLUME_END), volumeID, psm.storageSetting.Serial)

	return nil
}

func (psm *adminStorageManager) ExpandVolume(volumeID int, params gwymodel.ExpandVolumeParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_VOLUME_BEGIN), volumeID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| Params:%+v\n", params)

	err = gatewayObj.ExpandVolume(volumeID, params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_EXPAND_VOLUME_FAILED), volumeID, psm.storageSetting.Serial)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_EXPAND_VOLUME_END), volumeID, psm.storageSetting.Serial)
	return nil
}

func (psm *adminStorageManager) UpdateVolumeNickname(volumeID int, params gwymodel.UpdateVolumeNicknameParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_VOLUME_NICKNAME_BEGIN), volumeID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| Params:%+v\n", params)

	err = gatewayObj.UpdateVolumeNickname(volumeID, params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_VOLUME_NICKNAME_FAILED), volumeID, psm.storageSetting.Serial)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_VOLUME_NICKNAME_END), volumeID, psm.storageSetting.Serial)
	return nil
}

func (psm *adminStorageManager) UpdateVolumeReductionSettings(volumeID int, params gwymodel.UpdateVolumeReductionParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_VOLUME_REDUCTION_BEGIN), volumeID, psm.storageSetting.Serial)

	gatewayObj, err := psm.getGatewayManager()
	if err != nil {
		return err
	}

	log.WriteDebug("TFDebug| Params: SavingSetting=%+v, compressionAcceleration=%+v\n", params.SavingSetting, params.CompressionAcceleration)

	err = gatewayObj.UpdateVolumeReductionSettings(volumeID, params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_VOLUME_REDUCTION_FAILED), volumeID, psm.storageSetting.Serial)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_VOLUME_REDUCTION_END), volumeID, psm.storageSetting.Serial)
	return nil
}

// ------------------- Helpers -------------------
func (psm *adminStorageManager) getGatewayManager() (gwymanager.AdminStorageManager, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	setting := gwymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gwyObj, err := gwyimpl.NewEx(setting)
	if err != nil {
		log.WriteError("failed to get gateway manager: %v", err)
		return nil, fmt.Errorf("failed to get gateway manager: %w", err)
	}

	log.WriteDebug("TFDebug| Storage Serial:%v, ManagementIP:%v\n", psm.storageSetting.Serial, psm.storageSetting.MgmtIP)
	return gwyObj, nil
}
