#
# Hitachi VSP Supported Host Modes Data Retrieval
#
# This data source retrieves the supported host modes and host mode options
# from a Hitachi VSP storage system.
#

data "hitachi_vsp_supported_host_modes" "vsp_supported_host_modes" {
  serial = 12345
}

output "vsp_supported_host_modes" {
  value = data.hitachi_vsp_supported_host_modes.vsp_supported_host_modes
}
