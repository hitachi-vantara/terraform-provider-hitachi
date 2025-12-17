#
# Hitachi VSP One Server Path Management (Create / Read / Update / Delete)
#
# This section defines resource blocks for creating and managing server paths
# in the Hitachi Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_server_path" resource supports full lifecycle management (CRUD)
# of server paths — allowing you to provision new paths, modify existing ones, and
# remove them when no longer needed. It provides detailed control over HBA WWN or
# iSCSI name and port assignments, enabling consistent automation of path provisioning
# through Terraform.
#
# This resource is designed to work in conjunction with the "hitachi_vsp_one_server"
# resource, providing granular control over path management separate from server
# management for better modularity and state management.
#
# ------------------------------
# About Create Operations
# ------------------------------
# The create operation provisions a new server path in the storage system.
# During create:
#   - Either "hba_wwn" or "iscsi_name" must be specified, but not both.
#   - "port_ids" specifies the list of ports to associate with the path.
#   - The path is created on the specified existing server.
#   - The resource ID and all attributes are returned.
#
# ------------------------------
# About Update Operations
# ------------------------------
# The update operation allows modification of path properties:
#   - "port_ids" can be updated to add or remove port associations.
#   - Removed ports will be automatically deleted from the path.
#   - Added ports will be automatically associated with the path.
#   - Smart diffing ensures only necessary changes are applied.
#
# ------------------------------
# About Delete Operations
# ------------------------------
# The delete operation removes the server path from the storage system.
# All port associations will be cleaned up automatically.
#
# ------------------------------
# Path Types and Validation
# ------------------------------
# FC Paths:
#   - Use "hba_wwn" to specify the Fibre Channel HBA World Wide Name.
#   - Format: "XXXXXXXXXXXXXXXX" (16 hex characters).
#
# iSCSI Paths:
#   - Use "iscsi_name" to specify the iSCSI Qualified Name.
#   - Format: "iqn.yyyy-mm.domain:identifier" (standard IQN format).
#
# Either hba_wwn OR iscsi_name must be specified, but never both.
#

# Create an FC server path
resource "hitachi_vsp_one_server_path" "fc_path" {
  serial    = 12345
  server_id = 11
  hba_wwn   = "500104f00081b201"
  port_ids  = ["CL1-A", "CL1-B"]
}

output "fc_path_info" {
  description = "FC server path information"
  value       = hitachi_vsp_one_server_path.fc_path
}

# Create an iSCSI server path
resource "hitachi_vsp_one_server_path" "iscsi_path" {
  serial     = 12345
  server_id  = 12
  iscsi_name = "iqn.1991-05.com.example:server01"
  port_ids   = ["CL2-A", "CL2-B"]
}

output "iscsi_path_info" {
  description = "iSCSI server path information"
  value       = hitachi_vsp_one_server_path.iscsi_path
}

# Example of path update - this will add CL3-A and remove CL1-B
resource "hitachi_vsp_one_server_path" "fc_path_updated" {
  serial    = 12345
  server_id = 11
  hba_wwn   = "500104f00081b202"          # Different HBA WWN for different path
  port_ids  = ["CL1-A", "CL3-A"]          # Port CL1-B will be removed, CL3-A added
}

output "updated_path_info" {
  description = "Updated FC server path information"
  value       = hitachi_vsp_one_server_path.fc_path_updated
}
