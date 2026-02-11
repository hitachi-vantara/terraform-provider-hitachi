#
# Hitachi VSP One Port Information Retrieval
#
# This section defines a data source block that retrieves detailed information
# about a specific storage port managed by Hitachi storage systems through the
# Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_port" data source allows you to access configuration
# and property details of an existing port, enabling you to reference it in
# Terraform configurations and automate resource dependencies.
#
# The data source returns comprehensive information including:
#   - Basic port details (ID, protocol, speeds, security settings)
#   - Protocol-specific information (FC, iSCSI, NVMe-TCP)
#   - Network configuration (IP addresses, VLAN, TCP settings)
#
# Adjust the parameters (serial and id) to match your environment
# and retrieve information for the desired port.
#

data "hitachi_vsp_one_port" "port" {
  serial = 12345
  port_id = "CX-X"  # Port ID to retrieve
}

output "port_info" {
  description = "Detailed information about the port"
  value       = data.hitachi_vsp_one_port.port.port_info
}