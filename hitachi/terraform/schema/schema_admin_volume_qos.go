package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- Datasource Volume QOS Schema -------------------
var StorageVolumeQosSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system.",
	},
	"volume_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Volume ID.",
	},
	"threshold": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Threshold settings for the volume.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_upper_iops_enabled": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Enable/disable upper IOPS setting.",
				},
				"upper_iops": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Upper IOPS value. Not displayed if is_upper_iops_enabled is false.",
				},
				"is_upper_transfer_rate_enabled": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Enable/disable upper data transfer rate setting.",
				},
				"upper_transfer_rate": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Upper data transfer rate value (MiBps). Not displayed if is_upper_transfer_rate_enabled is false.",
				},
				"is_lower_iops_enabled": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Enable/disable lower IOPS setting.",
				},
				"lower_iops": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Lower IOPS target value. Not displayed if is_lower_iops_enabled is false.",
				},
				"is_lower_transfer_rate_enabled": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Enable/disable lower data transfer rate setting.",
				},
				"lower_transfer_rate": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Lower data transfer rate target value (MiBps). Not displayed if is_lower_transfer_rate_enabled is false.",
				},
				"is_response_priority_enabled": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Enable/disable I/O processing priority.",
				},
				"response_priority": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "I/O processing priority. The smaller the value, the lower the priority; the larger the value, the higher the priority. Not displayed if is_response_priority_enabled is false.",
				},
				"target_response_time": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Response time target value (ms). Value automatically set by the storage system according to response_priority. Not displayed if is_response_priority_enabled is false.",
				},
			},
		},
	},
	"alert_setting": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Alert settings for the volume.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_upper_alert_enabled": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Enable/disable alert setting for exceeding upper IOPS or data transfer rate.",
				},
				"upper_alert_allowable_time": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Permissible duration (s) when upper IOPS or data transfer rate is exceeded. Not displayed if is_upper_alert_enabled is false.",
				},
				"is_lower_alert_enabled": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Enable/disable alert setting for not reaching lower IOPS or data transfer rate.",
				},
				"lower_alert_allowable_time": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Permissible duration (s) for not reaching lower IOPS or data transfer rate. Not displayed if is_lower_alert_enabled is false.",
				},
				"is_response_alert_enabled": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Enable/disable alert setting when response time target is not met.",
				},
				"response_alert_allowable_time": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Permissible duration (s) when response time target is not met. Not displayed if is_response_alert_enabled is false.",
				},
			},
		},
	},
	"alert_time": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Alert times for the volume.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"upper_alert_time": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "UTC timestamp of the last alert for exceeding upper IOPS/data transfer rate continuity. Not displayed if alert is not set or if no alert has occurred.",
				},
				"lower_alert_time": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "UTC timestamp of the last alert for not reaching lower IOPS/data transfer rate continuity. Not displayed if alert is not set or if no alert has occurred.",
				},
				"response_alert_time": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "UTC timestamp of the last alert for continuously not reaching response time target. Not displayed if alert is not set or if no alert has occurred.",
				},
			},
		},
	},
}

// ------------------- Resource Set Volume QOS Schema -------------------
var StorageVolumeSetQosSchema = map[string]*schema.Schema{
	"serial": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Serial number of the storage system.",
	},
	"volume_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Volume ID.",
	},
	"threshold": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Threshold settings for the volume.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_upper_iops_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enable/disable upper IOPS setting.",
				},
				"upper_iops": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Upper IOPS value. Required if is_upper_iops_enabled is true.",
				},
				"is_upper_transfer_rate_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enable/disable upper data transfer rate setting.",
				},
				"upper_transfer_rate": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Upper data transfer rate value (MiBps). Required if is_upper_transfer_rate_enabled is true.",
				},
				"is_lower_iops_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enable/disable lower IOPS setting.",
				},
				"lower_iops": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Lower IOPS target value. Required if is_lower_iops_enabled is true.",
				},
				"is_lower_transfer_rate_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enable/disable lower data transfer rate setting.",
				},
				"lower_transfer_rate": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Lower data transfer rate target value (MiBps). Required if is_lower_transfer_rate_enabled is true.",
				},
				"is_response_priority_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enable/disable I/O processing priority.",
				},
				"response_priority": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "I/O processing priority. The smaller the value, the lower the priority; the larger the value, the higher the priority. Required if is_response_priority_enabled is true.",
				},
			},
		},
	},
	"alert_setting": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Alert settings for the volume.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_upper_alert_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enable/disable alert setting for exceeding upper IOPS or data transfer rate.",
				},
				"upper_alert_allowable_time": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Permissible duration (s) when upper IOPS or data transfer rate is exceeded. Required if is_upper_alert_enabled is true.",
				},
				"is_lower_alert_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enable/disable alert setting for not reaching lower IOPS or data transfer rate.",
				},
				"lower_alert_allowable_time": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Permissible duration (s) for not reaching lower IOPS or data transfer rate. Required if is_lower_alert_enabled is true.",
				},
				"is_response_alert_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enable/disable alert setting when response time target is not met.",
				},
				"response_alert_allowable_time": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Permissible duration (s) when response time target is not met. Required if is_response_alert_enabled is true.",
				},
			},
		},
	},
}
