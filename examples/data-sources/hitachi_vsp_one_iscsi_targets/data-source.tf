#
# Hitachi VSP One iSCSI Targets Information Retrieval
#
# This section defines a data source block that retrieves detailed information
# about all iSCSI targets associated with a specific server on a Hitachi Virtual
# Storage Platform One Block system.
#
# The "hitachi_vsp_one_iscsi_targets" data source allows you to list and inspect
# configuration details of multiple iSCSI targets, such as their port IDs,
# target iSCSI names.
#
# Adjust the parameters (for example, serial and server_id) to match your
# environment and retrieve the list of iSCSI targets for the desired server.
#
# Example:

data "hitachi_vsp_one_iscsi_targets" "targets_example" {
  serial    = 12345
  server_id = 10
}

output "iscsi_targets_info" {
  description = "Detailed information for all iSCSI targets of the specified server."
  value       = data.hitachi_vsp_one_iscsi_targets.targets_example.iscsi_targets_info
}

output "iscsi_targets_count" {
  description = "Number of iSCSI targets found for the specified server."
  value       = data.hitachi_vsp_one_iscsi_targets.targets_example.iscsi_targets_count
}
