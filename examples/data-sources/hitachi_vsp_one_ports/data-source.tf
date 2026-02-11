#
# Hitachi VSP One Ports Retrieval (Filtered List)
#
# This section defines a data source block that retrieves a list of existing
# storage ports from the Hitachi Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_ports" data source supports filtering by protocol,
# allowing you to narrow results to specific port types (FC, iSCSI, NVME_TCP).
# This is useful for auditing, reporting, or referencing existing ports in
# Terraform configurations.
#
# You can specify optional filtering parameters such as:
#   - protocol to target specific port types (FC, ISCSI, NVME_TCP)
#
# Adjust the parameters to fit your environment and query requirements.
#

##############
# Example without protocol filter for all ports
data "hitachi_vsp_one_ports" "ports" {
  serial = 12345
}

output "all_ports_info" {
  description = "Information about all ports"
  value       = data.hitachi_vsp_one_ports.ports.ports_info
}

output "all_ports_count" {
  description = "Total number of ports returned"
  value       = data.hitachi_vsp_one_ports.ports.port_count
}

##############
# Example with protocol filter for FC ports only
data "hitachi_vsp_one_ports" "fc_ports" {
  serial   = 12345
  protocol = "FC"
}

output "fc_ports_info" {
  description = "Information about FC ports only"
  value       = data.hitachi_vsp_one_ports.fc_ports.ports_info
}

output "fc_ports_count" {
  description = "Number of FC ports"
  value       = data.hitachi_vsp_one_ports.fc_ports.port_count
}

##############
# Example with protocol filter for iSCSI ports only
data "hitachi_vsp_one_ports" "iscsi_ports" {
  serial   = 12345
  protocol = "iSCSI"
}

output "iscsi_ports_info" {
  description = "Information about iSCSI ports only"
  value       = data.hitachi_vsp_one_ports.iscsi_ports.ports_info
}

output "iscsi_ports_count" {
  description = "Number of iSCSI ports"
  value       = data.hitachi_vsp_one_ports.iscsi_ports.port_count
}
