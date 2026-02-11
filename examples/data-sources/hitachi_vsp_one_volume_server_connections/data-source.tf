#
# Hitachi VSP One Volume–Server Connections Retrieval
#
# This section defines a data source block that retrieves detailed information
# about volume–server connections managed by Hitachi storage systems through the
# Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_volume_server_connections" data source allows you to query
# connection details between volumes and servers, providing visibility into how
# specific volumes are mapped to one or more servers.
#
# Adjust the parameters (`serial`, `server_id` or `server_nickname`, and optionally
# `start_volume_id/start_volume_id_hex` and `requested_count`) to match your environment and retrieve
# the desired set of volume–server connection records.
#
# Either the `server_id` or `server_nickname` parameter must be specified to
# identify the target server. If both are provided, an error will occur.
#

data "hitachi_vsp_one_volume_server_connections" "connections" {
  serial    = var.serial_number
  server_id = 9
  # server_nickname = "my_server"
  start_volume_id = 6700
  # or
  # start_volume_id_hex = "0X1A0"
  requested_count = 10
}

output "connections" {
  value = data.hitachi_vsp_one_volume_server_connections.connections.connections_info
}
output "connections_count" {
  value = data.hitachi_vsp_one_volume_server_connections.connections.connections_count
}
output "connections_total" {
  value = data.hitachi_vsp_one_volume_server_connections.connections.total_count
}
