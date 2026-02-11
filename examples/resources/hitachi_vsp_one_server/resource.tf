#
# Hitachi VSP One Server Management (Create / Read / Update / Delete)
#
# This section defines a resource block for creating and managing servers
# in the Hitachi Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_server" resource supports full lifecycle management (CRUD)
# of servers — allowing you to provision new servers, modify existing ones, and
# remove them when no longer needed. It provides detailed control over server
# configuration, OS type, and path assignments, enabling consistent automation
# of server provisioning through Terraform.
#
# ------------------------------
# About Create Operations
# ------------------------------
# The create operation provisions a new server in the storage system.
# During create:
#   - "server_id" must not be specified (it will be auto-assigned).
#   - "nickname", "os_type", and protocol can be specified.
#   - "host_groups" is optional.
#   - The created server ID and all attributes are returned.
#
# ------------------------------
# About Update Operations
# ------------------------------
# The update operation allows modification of server properties:
#   - "nickname" can be updated.
#   - "os_type" can be changed if supported.
#   - "host_groups" is optional.
#
# ------------------------------
# About Delete Operations
# ------------------------------
# The delete operation removes the server from the storage system.
# All associated configurations will be cleaned up.
#

resource "hitachi_vsp_one_server" "fc_server" {
  serial          = "12345"
  server_nickname = "MyFCServer"
  os_type         = "Linux"
  protocol        = "FC"
  is_reserved     = false

  # Either host_group_name or host_group_id must be specified, but not both.
  host_groups {
    port_id         = "CL1-A"
    host_group_name = "testhg1"
  }
  # host_groups {
  #   port_id       = "CL3-A"
  #   host_group_id = 252
  # }
}

output "server_info" {
  value = hitachi_vsp_one_server.fc_server.data
}

resource "hitachi_vsp_one_server" "isci_server" {
  serial          = "12345"
  server_nickname = "MyISCSIServer"
  os_type         = "Windows"
  protocol        = "iSCSI"
  is_reserved     = false

  # Either host_group_name or host_group_id must be specified, but not both.
  host_groups {
    port_id         = "CL2-B"
    host_group_name = "testhg2"
  }
  # host_groups {
  #   port_id       = "CL3-A"
  #   host_group_id = 253
  # }
}

output "iscsi_server_info" {
  value = hitachi_vsp_one_server.isci_server.data
}