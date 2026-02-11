#
# Hitachi VSP One Server HBAs Information Retrieval
#
# This section defines a data source block that retrieves detailed information
# about all server Host Bus Adapters (HBAs) for a specific server managed by 
# Hitachi storage systems through the Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_server_hbas" data source allows you to access configuration
# and property details of all existing server HBAs for a given server, enabling 
# you to reference them in Terraform configurations and automate resource dependencies.
#
# The data source returns comprehensive information for all HBAs including:
#   - Server ID and HBA WWN identification for each HBA
#   - iSCSI names (for iSCSI HBAs)
#   - Associated port IDs for each HBA
#   - HBA protocol types and configurations
#
# Adjust the parameters (serial and server_id) to match your environment
# and retrieve information for all HBAs of the desired server.
#

data "hitachi_vsp_one_server_hbas" "example" {
  serial    = 12345
  server_id = 23
}

output "server_hbas_info" {
  description = "Detailed information about all server HBAs"
  value       = data.hitachi_vsp_one_server_hbas.example
}