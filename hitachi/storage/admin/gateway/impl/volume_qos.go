package admin

import (
	"fmt"
	"strconv"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/admin/gateway/http"
	gatewaymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// GetVolumeQosAdminInfo obtains the QoS settings for a specific volume by ID.
func (psm *adminStorageManager) GetVolumeQosAdminInfo(volumeID int) (*gatewaymodel.VolumeQosResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var qosInfo gatewaymodel.VolumeQosResponse

	log.WriteDebug("TFDebug| Data not found in disk cache, call API")

	apiSuf := fmt.Sprintf("objects/volumes/%s/qos-setting", strconv.Itoa(volumeID))
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &qosInfo)
	log.WriteDebug("TFDebug| Data for call API: %+v", qosInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &qosInfo, nil
}

// SetVolumeQosAdminThreshold updates only the threshold settings for a volume.
func (psm *adminStorageManager) SetVolumeQosAdminThreshold(volumeID int, threshold gatewaymodel.VolumeQosThreshold) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("TFDebug| Updating QoS threshold for volume ID: %d with settings: %+v", volumeID, threshold)

	apiSuf := fmt.Sprintf("objects/volumes/%s/qos-setting", strconv.Itoa(volumeID))
	payload := map[string]interface{}{"threshold": threshold}
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, payload)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	log.WriteDebug("TFDebug| Successfully updated QoS threshold for volume ID: %d", volumeID)
	return nil
}

// SetVolumeQosAdminAlertSetting updates only the alert setting for a volume.
func (psm *adminStorageManager) SetVolumeQosAdminAlertSetting(volumeID int, alertSetting gatewaymodel.VolumeQosAlertSetting) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("TFDebug| Updating QoS alert setting for volume ID: %d with settings: %+v", volumeID, alertSetting)

	apiSuf := fmt.Sprintf("objects/volumes/%s/qos-setting", strconv.Itoa(volumeID))
	payload := map[string]interface{}{"alertSetting": alertSetting}
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, payload)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	log.WriteDebug("TFDebug| Successfully updated QoS alert setting for volume ID: %d", volumeID)
	return nil
}
