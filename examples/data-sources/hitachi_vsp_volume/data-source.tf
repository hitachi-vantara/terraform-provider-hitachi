#
# Hitachi VSP Volume Data Retrieval from VSP Direct Connect
#
# This section defines a data source block to fetch information about a specific
# volume from a Hitachi Virtual Storage Platform (VSP) using HashiCorp Configuration
# Language (HCL).
#
# The data source block "hitachi_vsp_volume" retrieves details about a volume
# associated with the provided parameters. This allows you to access configuration
# and property information for the specified volume.
#
# Customize the values of the parameters (serial, ldev_id) to match your environment,
# enabling you to retrieve information about the desired volume.
#

data "hitachi_vsp_volume" "volume" {
  serial  = 12345
  ldev_id = 281
}

output "volume" {
  value = data.hitachi_vsp_volume.volume
}
