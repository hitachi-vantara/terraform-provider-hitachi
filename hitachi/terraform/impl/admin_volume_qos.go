package terraform

import (
	"fmt"
	"strconv"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	// utils "terraform-provider-hitachi/hitachi/common/utils"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	provimpladmin "terraform-provider-hitachi/hitachi/storage/admin/provisioner/impl"
	provmodeladmin "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"

	terrcommon "terraform-provider-hitachi/hitachi/terraform/common"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetAdminVolumeQos(d *schema.ResourceData) (*terraformmodel.VolumeQosResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	finalLdev, err := terrcommon.ExtractLdevFields(d, "volume_id", "volume_id_hex")
	if err != nil {
		return nil, err
	}
	if finalLdev == nil {
		return nil, fmt.Errorf("either volume_id or volume_id_hex must be specified")
	}
	volumeID := *finalLdev

	storageSetting, err := cache.GetAdminSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := provmodeladmin.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	provObj, err := provimpladmin.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUME_QOS_ADMIN_BEGIN), setting.MgmtIP)
	reconVolumeQos, err := provObj.GetVolumeQosAdminInfo(volumeID)
	if err != nil {
		log.WriteDebug("TFError| error getting volume QoS, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_VOLUME_QOS_ADMIN_FAILED), setting.MgmtIP)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformVolumeQos := terraformmodel.VolumeQosResponse{}
	err = copier.Copy(&terraformVolumeQos, reconVolumeQos)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_VOLUME_QOS_ADMIN_END), setting.MgmtIP)

	return &terraformVolumeQos, nil
}

func SetAdminVolumeQosThreshold(d *schema.ResourceData, serial string, volumeID int, qosSettings terraformmodel.VolumeQosThreshold) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageSetting, err := cache.GetAdminSettingsFromCache(serial)
	if err != nil {
		return err
	}

	setting := provmodeladmin.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	provObj, err := provimpladmin.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_VOLUME_QOS_ADMIN_BEGIN), volumeID)

	// Converting terraform to reconciler
	provQosThresholdSettings := provmodeladmin.VolumeQosThreshold{}
	err = copier.Copy(&provQosThresholdSettings, qosSettings)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from terraform to reconciler structure, err: %v", err)
		return err
	}

	// Set threshold if present
	if provQosThresholdSettings == (provmodeladmin.VolumeQosThreshold{}) {
		log.WriteInfo(mc.GetMessage(mc.INFO_SET_VOLUME_QOS_ADMIN_END), volumeID)
		return nil
	}

	err = provObj.SetVolumeQosAdminThreshold(volumeID, provQosThresholdSettings)
	if err != nil {
		log.WriteDebug("TFError| error setting volume QoS threshold, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_VOLUME_QOS_ADMIN_FAILED), volumeID, err)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_VOLUME_QOS_ADMIN_END), volumeID)

	return nil
}

// SetAdminVolumeQosAlertSetting sets only the alert_setting for a volume
func SetAdminVolumeQosAlertSetting(d *schema.ResourceData, serial string, volumeID int, qosSettings terraformmodel.VolumeQosAlertSetting) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageSetting, err := cache.GetAdminSettingsFromCache(serial)
	if err != nil {
		return err
	}

	setting := provmodeladmin.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	provObj, err := provimpladmin.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_VOLUME_QOS_ADMIN_BEGIN), volumeID)

	// Converting terraform to reconciler
	provQosAlertSettings := provmodeladmin.VolumeQosAlertSetting{}
	err = copier.Copy(&provQosAlertSettings, qosSettings)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from terraform to reconciler structure, err: %v", err)
		return err
	}
	// Set alert_setting if present
	if provQosAlertSettings == (provmodeladmin.VolumeQosAlertSetting{}) {
		log.WriteInfo(mc.GetMessage(mc.INFO_SET_VOLUME_QOS_ADMIN_END), volumeID)
		return nil
	}
	err = provObj.SetVolumeQosAdminAlertSetting(volumeID, provQosAlertSettings)
	if err != nil {
		log.WriteDebug("TFError| error setting volume QoS alert_setting, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_VOLUME_QOS_ADMIN_FAILED), volumeID, err)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_VOLUME_QOS_ADMIN_END), volumeID)

	return nil
}

func ResourceAdminVolumeQosRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	finalLdev, err := terrcommon.ExtractLdevFields(d, "volume_id", "volume_id_hex")
	if err != nil {
		return diag.FromErr(err)
	}
	if finalLdev == nil {
		return diag.FromErr(fmt.Errorf("either volume_id or volume_id_hex must be specified"))
	}
	volumeID := *finalLdev

	log.WriteDebug("TFDebug| ResourceAdminVolumeQosRead: serial=%d, volumeID=%d", serial, volumeID)

	// Call the API to get the QoS info for the volume
	resp, err := GetAdminVolumeQos(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.WriteDebug("GetAdminVolumeQos response: %+v", resp)

	// Set the ID for the resource
	d.SetId(fmt.Sprintf("%d-%d", serial, resp.VolumeId))

	// Set qos values

	// these are inputs, they can't be set in update.
	// these schema fields cannot have Computed: true because they are inputs,
	// otherwise terraform would think both are set at the same time in next run, and fail
	// better solution: overhaul the schema to separate inputs and outputs.
	// d.Set("volume_id", volumeID)
	// d.Set("volume_id_hex", utils.IntToHexString(volumeID))

	d.Set("threshold", []interface{}{
		map[string]interface{}{
			"is_upper_iops_enabled":          resp.Threshold.IsUpperIopsEnabled,
			"upper_iops":                     resp.Threshold.UpperIops,
			"is_upper_transfer_rate_enabled": resp.Threshold.IsUpperTransferRateEnabled,
			"upper_transfer_rate":            resp.Threshold.UpperTransferRate,
			"is_lower_iops_enabled":          resp.Threshold.IsLowerIopsEnabled,
			"lower_iops":                     resp.Threshold.LowerIops,
			"is_lower_transfer_rate_enabled": resp.Threshold.IsLowerTransferRateEnabled,
			"lower_transfer_rate":            resp.Threshold.LowerTransferRate,
			"is_response_priority_enabled":   resp.Threshold.IsResponsePriorityEnabled,
			"response_priority":              resp.Threshold.ResponsePriority,
			"target_response_time":           resp.Threshold.TargetResponseTime,
		},
	})
	d.Set("alert_setting", []interface{}{
		map[string]interface{}{
			"is_upper_alert_enabled":        resp.AlertSetting.IsUpperAlertEnabled,
			"upper_alert_allowable_time":    resp.AlertSetting.UpperAlertAllowableTime,
			"is_lower_alert_enabled":        resp.AlertSetting.IsLowerAlertEnabled,
			"lower_alert_allowable_time":    resp.AlertSetting.LowerAlertAllowableTime,
			"is_response_alert_enabled":     resp.AlertSetting.IsResponseAlertEnabled,
			"response_alert_allowable_time": resp.AlertSetting.ResponseAlertAllowableTime,
		},
	})
	d.Set("alert_time", []interface{}{
		map[string]interface{}{
			"upper_alert_time":    resp.AlertTime.UpperAlertTime,
			"lower_alert_time":    resp.AlertTime.LowerAlertTime,
			"response_alert_time": resp.AlertTime.ResponseAlertTime,
		},
	})

	return nil
}

func ResourceAdminVolumeQosCreate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	finalLdev, err := terrcommon.ExtractLdevFields(d, "volume_id", "volume_id_hex")
	if err != nil {
		return diag.FromErr(err)
	}
	if finalLdev == nil {
		return diag.FromErr(fmt.Errorf("either volume_id or volume_id_hex must be specified"))
	}
	volumeID := *finalLdev

	log.WriteDebug("TFDebug| ResourceAdminVolumeQosCreate: serial=%d, volumeID=%d", serial, volumeID)

	thresholdList := d.Get("threshold").([]interface{})
	alertSettingList := d.Get("alert_setting").([]interface{})

	var qosSettingsWithThreshold terraformmodel.VolumeQosThreshold
	var qosSettingsWithAlert terraformmodel.VolumeQosAlertSetting

	if len(thresholdList) > 0 && thresholdList[0] != nil {
		th := thresholdList[0].(map[string]interface{})
		qosSettingsWithThreshold = terraformmodel.VolumeQosThreshold{
			IsUpperIopsEnabled:         th["is_upper_iops_enabled"].(bool),
			UpperIops:                  th["upper_iops"].(int),
			IsUpperTransferRateEnabled: th["is_upper_transfer_rate_enabled"].(bool),
			UpperTransferRate:          th["upper_transfer_rate"].(int),
			IsLowerIopsEnabled:         th["is_lower_iops_enabled"].(bool),
			LowerIops:                  th["lower_iops"].(int),
			IsLowerTransferRateEnabled: th["is_lower_transfer_rate_enabled"].(bool),
			LowerTransferRate:          th["lower_transfer_rate"].(int),
			IsResponsePriorityEnabled:  th["is_response_priority_enabled"].(bool),
			ResponsePriority:           th["response_priority"].(int),
		}
	}

	if len(alertSettingList) > 0 && alertSettingList[0] != nil {
		al := alertSettingList[0].(map[string]interface{})
		qosSettingsWithAlert = terraformmodel.VolumeQosAlertSetting{
			IsUpperAlertEnabled:        al["is_upper_alert_enabled"].(bool),
			UpperAlertAllowableTime:    al["upper_alert_allowable_time"].(int),
			IsLowerAlertEnabled:        al["is_lower_alert_enabled"].(bool),
			LowerAlertAllowableTime:    al["lower_alert_allowable_time"].(int),
			IsResponseAlertEnabled:     al["is_response_alert_enabled"].(bool),
			ResponseAlertAllowableTime: al["response_alert_allowable_time"].(int),
		}
	}

	// Set threshold and alert_setting if present
	err = SetAdminVolumeQosThreshold(d, strconv.Itoa(serial), volumeID, qosSettingsWithThreshold)
	if err != nil {
		return diag.FromErr(err)
	}
	log.WriteDebug("TFDebug| ResourceAdminVolumeQosCreate: successfully updated QoS Threshold for volumeID=%d with settings=%v", volumeID, qosSettingsWithThreshold)

	err = SetAdminVolumeQosAlertSetting(d, strconv.Itoa(serial), volumeID, qosSettingsWithAlert)
	if err != nil {
		return diag.FromErr(err)
	}
	log.WriteDebug("TFDebug| ResourceAdminVolumeQosCreate: successfully updated QoS Alert for volumeID=%d with settings=%v", volumeID, qosSettingsWithAlert)

	id := strconv.Itoa(serial) + "-" + strconv.Itoa(volumeID)
	d.SetId(id)
	return ResourceAdminVolumeQosRead(d)
}

func ResourceAdminVolumeQosUpdate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	finalLdev, err := terrcommon.ExtractLdevFields(d, "volume_id", "volume_id_hex")
	if err != nil {
		return diag.FromErr(err)
	}
	if finalLdev == nil {
		return diag.FromErr(fmt.Errorf("either volume_id or volume_id_hex must be specified"))
	}
	volumeID := *finalLdev

	thresholdList := d.Get("threshold").([]interface{})
	alertSettingList := d.Get("alert_setting").([]interface{})

	thresholdChanged := d.HasChange("threshold")
	alertChanged := d.HasChange("alert_setting")

	if thresholdChanged && len(thresholdList) > 0 && thresholdList[0] != nil {
		th := thresholdList[0].(map[string]interface{})
		qosSettingsWithThreshold := terraformmodel.VolumeQosThreshold{
			IsUpperIopsEnabled:         th["is_upper_iops_enabled"].(bool),
			UpperIops:                  th["upper_iops"].(int),
			IsUpperTransferRateEnabled: th["is_upper_transfer_rate_enabled"].(bool),
			UpperTransferRate:          th["upper_transfer_rate"].(int),
			IsLowerIopsEnabled:         th["is_lower_iops_enabled"].(bool),
			LowerIops:                  th["lower_iops"].(int),
			IsLowerTransferRateEnabled: th["is_lower_transfer_rate_enabled"].(bool),
			LowerTransferRate:          th["lower_transfer_rate"].(int),
			IsResponsePriorityEnabled:  th["is_response_priority_enabled"].(bool),
			ResponsePriority:           th["response_priority"].(int),
		}
		err := SetAdminVolumeQosThreshold(d, strconv.Itoa(serial), volumeID, qosSettingsWithThreshold)
		if err != nil {
			return diag.FromErr(err)
		}
		log.WriteDebug("TFDebug| ResourceAdminVolumeQosUpdate: updated QoS Threshold for volumeID=%d with settings=%v", volumeID, qosSettingsWithThreshold)
		// Set the updated values in state
		d.Set("threshold", []interface{}{th})
	} else {
		log.WriteDebug("TFDebug| ResourceAdminVolumeQosUpdate: threshold not changed, skipping update for volumeID=%d", volumeID)
	}

	if alertChanged && len(alertSettingList) > 0 && alertSettingList[0] != nil {
		al := alertSettingList[0].(map[string]interface{})
		qosSettingsWithAlert := terraformmodel.VolumeQosAlertSetting{
			IsUpperAlertEnabled:        al["is_upper_alert_enabled"].(bool),
			UpperAlertAllowableTime:    al["upper_alert_allowable_time"].(int),
			IsLowerAlertEnabled:        al["is_lower_alert_enabled"].(bool),
			LowerAlertAllowableTime:    al["lower_alert_allowable_time"].(int),
			IsResponseAlertEnabled:     al["is_response_alert_enabled"].(bool),
			ResponseAlertAllowableTime: al["response_alert_allowable_time"].(int),
		}
		err := SetAdminVolumeQosAlertSetting(d, strconv.Itoa(serial), volumeID, qosSettingsWithAlert)
		if err != nil {
			return diag.FromErr(err)
		}
		log.WriteDebug("TFDebug| ResourceAdminVolumeQosUpdate: updated QoS Alert for volumeID=%d with settings=%v", volumeID, qosSettingsWithAlert)
		// Set the updated values in state
		d.Set("alert_setting", []interface{}{al})
	} else {
		log.WriteDebug("TFDebug| ResourceAdminVolumeQosUpdate: alert_setting not changed, skipping update for volumeID=%d", volumeID)
	}

	d.SetId(fmt.Sprintf("%d-%d", serial, volumeID))
	return ResourceAdminVolumeQosRead(d)
}
