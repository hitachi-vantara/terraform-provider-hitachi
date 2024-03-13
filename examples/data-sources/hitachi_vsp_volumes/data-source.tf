#
# Hitachi VSP Volume Data Retrieval from VSP Direct Connect and Gateway connections
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

  // Mandatory parameters for both the provider
  serial  = 12345
  ldev_id = 281

  // Optional parameters for the hitachi_infrastructure_gateway_provider

  subscriber_id = ""
}

output "volume" {
  value = data.hitachi_vsp_volume.volume
}

#
# The data source block "hitachi_vsp_volumes" retrieves details about volumes within a
# specified range of logical device IDs (LDEVs). This allows you to access configuration
# and property information for the specified volumes.
#
# Customize the values of the parameters (serial, start_ldev_id, end_ldev_id, undefined_ldev)
# to match your environment. By doing so, you can retrieve information about the desired range
# of volumes while indicating whether undefined LDEVs should be included or not.
#

data "hitachi_vsp_volumes" "volume1" {
  
  // Mandatory parameters for both the providers
  serial         = 12345

  // optional parameters for both the providers
  start_ldev_id  = 280
  end_ldev_id    = 285
  undefined_ldev = false

// Optional parameters for the hitachi_infrastructure_gateway_provider

  subscriber_id = ""

}

output "volume1" {
  value = data.hitachi_vsp_volumes.volume1
}
