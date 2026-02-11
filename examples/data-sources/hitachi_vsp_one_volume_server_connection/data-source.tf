#
# Hitachi VSP One Volume–Server Connection Retrieval
#
# This section defines a data source block that retrieves detailed information
# about a single volume–server connection managed by Hitachi storage systems
# through the Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_volume_server_connection" data source allows you to obtain
# configuration and mapping details for a specific connection between a volume
# and a server, providing insight into how the volume is associated with that
# server.
#
# All input parameters (`serial`, `server_id`, and `volume_id/volume_id_hex`) are required and
# must match the identifiers of the target storage system, server, and volume
# you want to query.
#
# Adjust these parameters to reflect your environment and retrieve detailed
# information about the desired volume–server connection.
#

data "hitachi_vsp_one_volume_server_connection" "connection" {
  serial    = var.serial_number
  server_id = 9
  volume_id = 6079
  # or
  # volume_id_hex = "0X17BF"
}

output "connection" {
  value = data.hitachi_vsp_one_volume_server_connection.connection.connection_info
}
