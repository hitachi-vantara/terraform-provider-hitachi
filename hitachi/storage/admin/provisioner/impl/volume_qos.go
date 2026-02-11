package admin

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/admin/gateway/impl"
	admingatewaymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/admin/provisioner/message-catalog"
	adminmodel "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetVolumeQosAdminInfo obtains the QoS settings for a specific volume by ID via the gateway layer.
func (psm *adminStorageManager) GetVolumeQosAdminInfo(volumeID int) (*adminmodel.VolumeQosResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := admingatewaymodel.StorageDeviceSettings{
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		Serial:   psm.storageSetting.Serial,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUME_QOS_ADMIN_BEGIN))
	qosInfo, err := gatewayObj.GetVolumeQosAdminInfo(volumeID)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetVolumeQosAdminInfo, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_VOLUME_QOS_ADMIN_FAILED))
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUME_QOS_ADMIN_END))

	provQosInfo := adminmodel.VolumeQosResponse{}
	err = copier.Copy(&provQosInfo, qosInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &provQosInfo, nil
}

// SetVolumeQosAdminThreshold sets the QoS threshold for a specific volume by ID via the gateway layer.
func (psm *adminStorageManager) SetVolumeQosAdminThreshold(volumeID int, threshold adminmodel.VolumeQosThreshold) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_VOLUME_QOS_ADMIN_THRESHOLD_BEGIN), volumeID)

	objStorage := admingatewaymodel.StorageDeviceSettings{
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		Serial:   psm.storageSetting.Serial,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	gwThreshold := admingatewaymodel.VolumeQosThreshold{}
	err = copier.Copy(&gwThreshold, &threshold)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from provisioner to gateway threshold, err: %v", err)
		return err
	}

	err = gatewayObj.SetVolumeQosAdminThreshold(volumeID, gwThreshold)
	if err != nil {
		log.WriteDebug("TFError| error setting volume QoS threshold, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_VOLUME_QOS_ADMIN_THRESHOLD_FAILED), volumeID, err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_VOLUME_QOS_ADMIN_THRESHOLD_END), volumeID)
	return nil
}

// SetVolumeQosAdminAlertSetting sets the QoS alert setting for a specific volume by ID via the gateway layer.
func (psm *adminStorageManager) SetVolumeQosAdminAlertSetting(volumeID int, alert adminmodel.VolumeQosAlertSetting) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_VOLUME_QOS_ADMIN_ALERT_BEGIN), volumeID)

	objStorage := admingatewaymodel.StorageDeviceSettings{
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		Serial:   psm.storageSetting.Serial,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	gwAlert := admingatewaymodel.VolumeQosAlertSetting{}
	err = copier.Copy(&gwAlert, &alert)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from provisioner to gateway alert, err: %v", err)
		return err
	}

	err = gatewayObj.SetVolumeQosAdminAlertSetting(volumeID, gwAlert)
	if err != nil {
		log.WriteDebug("TFError| error setting volume QoS alert_setting, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_VOLUME_QOS_ADMIN_ALERT_FAILED), volumeID, err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_VOLUME_QOS_ADMIN_ALERT_END), volumeID)
	return nil
}
