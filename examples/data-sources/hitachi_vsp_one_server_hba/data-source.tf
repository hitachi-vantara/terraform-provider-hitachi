#
# Hitachi VSP One Server HBA Information Retrieval
#
# This section defines a data source block that retrieves detailed information
# about a specific server Host Bus Adapter (HBA) managed by Hitachi storage
# systems through the Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_server_hba" data source allows you to access configuration
# and property details of an existing server HBA, enabling you to reference it
# in Terraform configurations and automate resource dependencies.
#
# The data source returns comprehensive information including:
#   - Server ID and HBA WWN identification
#   - iSCSI name (for iSCSI HBAs)
#   - Associated port IDs
#   - HBA protocol type and configuration
#
# Adjust the parameters (serial, server_id, and hba_wwn) to match your environment
# and retrieve information for the desired server HBA.
#

data "hitachi_vsp_one_server_hba" "example" {
  serial    = 12345
  server_id = 23
  initiator_name   = "iqn.server107"
}

output "server_info" {
  description = "Detailed information about the server HBA"
  value       = data.hitachi_vsp_one_server_hba.example
}