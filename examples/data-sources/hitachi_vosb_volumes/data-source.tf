#
# Hitachi VSP One SDS Block Volumes Data Retrieval
#
# This section defines a data source block to fetch information about multiple volumes
# from a Hitachi VSP One SDS Block using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vosb_volumes" retrieves details about volumes
# associated with the provided parameters. This allows you to access configuration
# and property information for the specified volumes.
#
# Customize the value of the parameter (vosb_address) to match your environment,
# enabling you to retrieve information about the desired volumes.
#

data "hitachi_vosb_volumes" "vosbvolumes" {
  vosb_address = "10.10.12.13"
}

output "volumes_output" {
  value = data.hitachi_vosb_volumes.vosbvolumes
}
