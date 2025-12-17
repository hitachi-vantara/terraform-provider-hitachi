#
# Hitachi VSP One Server Information Retrieval
#
# This section defines a data source block that retrieves detailed information
# about a specific server managed by Hitachi storage systems through the
# Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_server" data source allows you to access configuration
# and property details of an existing server, enabling you to reference server
# information in Terraform configurations and automate resource dependencies.
#
# This data source is useful for retrieving server details such as nickname,
# OS type, and port configurations for further automation or validation.
#

data "hitachi_vsp_one_server" "server_info" {
  serial    = 12345
  server_id = 11
}

output "server_details" {
  value = data.hitachi_vsp_one_server.server_info
}

data "hitachi_vsp_one_server" "server_path_info" {
  serial    = 12345
  server_id = 11
}

# To list paths associated with the server
output "server_paths_info" {
  value = data.hitachi_vsp_one_server.server_path_info.data[0].paths
}