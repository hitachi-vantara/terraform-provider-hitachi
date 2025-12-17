#
# Hitachi VSP One Servers List Retrieval
#
# This section defines a data source block that retrieves a list of all servers
# managed by Hitachi storage systems through the Virtual Storage Platform One
# Block Administrator.
#
# The "hitachi_vsp_one_servers" data source allows you to access information about
# all configured servers in the storage system, enabling you to reference server
# details in Terraform configurations and automate resource dependencies.
#
# This data source is particularly useful for discovering available servers
# before creating server-specific configurations or for inventory management.
#

data "hitachi_vsp_one_servers" "all_servers" {
  serial = 12345
}

output "servers_list" {
  value = data.hitachi_vsp_one_servers.all_servers
}
