#
# Hitachi VSP One Server Path Data Source
#
# This section documents the data source block for retrieving Fibre Channel (FC)
# server path information managed by Hitachi storage systems using the
# Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_server_path" data source provides access to FC path details,
# including server serial, server ID, HBA WWN, and port ID, allowing you to
# reference path information in Terraform and automate resource dependencies.
#
# Use this data source to obtain FC path attributes for automation, validation,
# or integration tasks involving server connectivity and storage configuration.
#

# Get FC server path information
data "hitachi_vsp_one_server_path" "fc_path_info" {
  serial    = 12345
  server_id = 11
  hba_wwn   = "500104f00081b201"
  port_id   = "CL1-A"
}

output "fc_path_details" {
  description = "FC server path information"
  value       = data.hitachi_vsp_one_server_path.fc_path_info
}

# Get iSCSI server path information
data "hitachi_vsp_one_server_path" "iscsi_path_info" {
  serial      = 12345
  server_id   = 11
  iscsi_name  = "iqn.1993-08.org.debian:01"
  port_id     = "CL2-A"
}

output "iscsi_path_details" {
  description = "iSCSI server path information"
  value       = data.hitachi_vsp_one_server_path.iscsi_path_info
}