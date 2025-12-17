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
# Provide the storage system serial number and ldev_id.
#

data "hitachi_vsp_volume" "volume" {
  serial  = 12345
  ldev_id = 281
}

output "volume" {
  value = data.hitachi_vsp_volume.volume
}