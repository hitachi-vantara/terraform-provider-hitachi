#
# Hitachi VSP One Volume QoS Information Retrieval
#
# This section defines a data source block that retrieves the QoS (Quality of Service)
# settings for a specific volume managed by Hitachi storage systems through the
# Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_volume_qos" data source allows you to access QoS configuration
# details of an existing volume, enabling you to reference these settings in
# Terraform configurations and automate resource management based on QoS policies.
#
# Adjust the parameters (such as serial or volume_id) to match your environment
# and retrieve QoS information for the desired volume.
#

data "hitachi_vsp_one_volume_qos" "volume_qos" {
  serial = 12345
  volume_id = 678
}


output "volume_qos_info" {
  value = data.hitachi_vsp_one_volume_qos.volume_qos
}