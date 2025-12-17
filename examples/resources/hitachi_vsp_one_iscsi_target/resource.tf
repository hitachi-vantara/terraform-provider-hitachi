#
# Hitachi VSP One iSCSI Target Rename
#
# This section defines a resource block that manages the iSCSI target name
# of an existing iSCSI target on a Hitachi Virtual Storage Platform One Block system.
#
# The "hitachi_vsp_one_iscsi_target" resource does not create or delete backend
# iSCSI targets. Instead, it is used exclusively to rename an existing target
# by updating its iSCSI target name on the specified port.
#
# The "target_iscsi_name" parameter is optional. When provided, the resource
# updates the backend iSCSI target name to the specified value. When omitted,
# no rename is performed and the current target name remains unchanged.
#
# Adjust the parameters (for example, serial, server_id, and port_id) to match
# your environment and rename the desired iSCSI target.
#
# Example:

resource "hitachi_vsp_one_iscsi_target" "target_example" {
  serial            = 12345
  server_id         = 10
  port_id           = "CL1-A"
  target_iscsi_name = "iqn.sample" # Optional: specify only to rename the iSCSI target
}

output "iscsi_target_info" {
  description = "Detailed information about the renamed iSCSI target."
  value       = resource.hitachi_vsp_one_iscsi_target.target_example.iscsi_target_info
}
