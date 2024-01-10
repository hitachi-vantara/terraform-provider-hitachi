#
# Hitachi VSS Block Volume Data Retrieval
#
# This section defines a data source block to fetch information about a specific volume
# from a Hitachi Virtual Storage System (VSS) using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vss_block_volume" retrieves details about a volume
# associated with the provided parameters. This allows you to access configuration
# and property information for the specified volume.
#
# Customize the values of the parameters (vss_block_address, volume_name) to match your
# environment, enabling you to retrieve information about the desired volume.
#

data "hitachi_vss_block_volume" "vssbvolumes" {
  vss_block_address = "10.10.12.13"
  volume_name       = "Mongonode3_vol4dd"
}

output "volumeoutput" {
  value = data.hitachi_vss_block_volume.vssbvolumes
}
