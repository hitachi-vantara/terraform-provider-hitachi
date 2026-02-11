#
# Manage QoS settings for a Hitachi VSP One Volume
#
# This example demonstrates how to configure QoS (Quality of Service) settings
# for a specific volume in a Hitachi storage system using the
# "hitachi_vsp_one_volume_qos" resource.
#
# Adjust the parameters (such as serial, volume_id, threshold and alert_setting) to match your environment.
#

resource "hitachi_vsp_one_volume_qos" "example" {
  serial    = 12345
  volume_id = 678
  # or
  # volume_id_hex = "0X2A6"
  threshold {
    is_upper_iops_enabled          = true
    upper_iops                     = 3000
    is_upper_transfer_rate_enabled = true
    upper_transfer_rate            = 200
    is_lower_iops_enabled          = true
    lower_iops                     = 1000
    is_lower_transfer_rate_enabled = true
    lower_transfer_rate            = 50
    is_response_priority_enabled   = true
    response_priority              = 3
  }
  alert_setting {
    is_upper_alert_enabled        = true
    upper_alert_allowable_time    = 60
    is_lower_alert_enabled        = false
    lower_alert_allowable_time    = 0
    is_response_alert_enabled     = true
    response_alert_allowable_time = 30
  }
}

output "volume_qos_info" {
  value = hitachi_vsp_one_volume_qos.example
}
