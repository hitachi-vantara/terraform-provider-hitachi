#
# Hitachi VSP One Volume–Server Connections Management
#
# This section defines a resource block that creates and manages connections
# between volumes and servers in Hitachi storage systems through the
# Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_volume_server_connection" resource allows you to attach
# one or more volumes to one or more servers, enabling automated configuration
# and management of volume–server mappings.
#
# All input parameters (`serial`, `volume_ids`, and `server_ids`) are required.
# The `serial` parameter identifies the target storage system, while the
# `volume_ids` and `server_ids` lists define which volumes should be connected
# to which servers.
#
# Adjust these parameters to reflect your environment and use Terraform to
# create, update, or remove the desired volume–server connections.
#

resource "hitachi_vsp_one_volume_server_connection" "connections" {
  serial     = var.serial_number
  server_ids = [9, 10]
  volume_ids = [6078, 6079]
  # or
  # volume_id_hexs = ["0X17BE", "0X17BF"]

}

output "connections" {
  value = resource.hitachi_vsp_one_volume_server_connection.connections.connections_info
}
