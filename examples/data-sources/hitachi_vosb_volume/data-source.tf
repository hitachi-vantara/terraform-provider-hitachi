#
# Hitachi VOS Block Volume Data Retrieval
#
# This section defines a data source block to fetch information about a specific volume
# from a Hitachi VSP One SDS Block (VOSB) using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vosb_volume" retrieves details about a volume
# associated with the provided parameters. This allows you to access configuration
# and property information for the specified volume.
#
# Customize the values of the parameters (vosb_address, volume_name) to match your
# environment, enabling you to retrieve information about the desired volume.
#

data "hitachi_vosb_volume" "vosbvolumes" {
  vosb_address = "10.10.12.13"
  volume_name  = "Mongonode3_vol4dd"
}

output "volumeoutput" {
  value = data.hitachi_vosb_volume.vosbvolumes
}
