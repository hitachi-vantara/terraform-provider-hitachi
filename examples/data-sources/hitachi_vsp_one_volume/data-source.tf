#
# Hitachi VSP One Volume Information Retrieval
#
# This section defines a data source block that retrieves detailed information
# about a specific volume managed by Hitachi storage systems through the
# Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_volume" data source allows you to access configuration
# and property details of an existing volume, enabling you to reference it in
# Terraform configurations and automate resource dependencies.
#
# Adjust the parameters (for example, serial or volume_id/volume_id_hex) to match your
# environment and retrieve information for the desired volume.
#

data "hitachi_vsp_one_volume" "volume" {
  serial    = var.serial_number
  volume_id = 281
  # or
  # volume_id_hex = "0X119"
}

output "volume_info" {
  value = data.hitachi_vsp_one_volume.volume.volume_info
}
