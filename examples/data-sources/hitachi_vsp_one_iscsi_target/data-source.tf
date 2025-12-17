#
# Hitachi VSP One iSCSI Target Information Retrieval
#
# This section defines a data source block that retrieves detailed information
# about a specific iSCSI target managed by Hitachi storage systems through the
# Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_iscsi_target" data source enables you to access configuration
# and property details of an existing iSCSI target, such as its port ID and target
# iSCSI name.
#
# Adjust the parameters (for example, serial, server_id, and port_id) to match your
# environment and retrieve information for the desired iSCSI target.
#
# Example:

data "hitachi_vsp_one_iscsi_target" "target_example" {
  serial    = 12345
  server_id = 10
  port_id   = "CL1-A"
}

output "target_name" {
  value = data.hitachi_vsp_one_iscsi_target.target_example.iscsi_target_info
}
