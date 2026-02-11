package terraform

import (
	"context"
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageVolumeQos() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage Volume QoS: It returns the QoS settings for a specific storage volume.",
		ReadContext: dataSourceStorageVolumeQosRead,
		Schema:      schemaimpl.StorageVolumeQosSchema,
	}
}

func dataSourceStorageVolumeQosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	volumeID := d.Get("volume_id").(int)

	// Call the API to get the QoS info for the volume
	resp, err := impl.GetAdminVolumeQos(d, volumeID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.WriteInfo("Successfully retrieved QoS info for volume ID: %d and response: %+v", volumeID, resp)

	// Set the ID for the resource
	d.SetId(fmt.Sprintf("%d-%d", serial, resp.VolumeId))

	// Set qos values
	d.Set("volume_id", resp.VolumeId)
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
