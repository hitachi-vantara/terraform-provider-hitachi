#
# Hitachi VSP One Server HBA Management (Add / Remove HBA Information)
#
# This example demonstrates how to manage Host Bus Adapters (HBAs) for servers
# in the Hitachi Virtual Storage Platform One Block Administrator using Terraform.
#
# The "hitachi_vsp_one_server_hba" resource allows you to add or remove HBA information
# (such as Fibre Channel WWN or iSCSI initiator name) for a given server. This enables
# you to automate the registration and management of server connectivity in your storage
# environment.
#
# ------------------------------
# Resource Arguments
# ------------------------------
# - serial: Serial number of the storage system (Required)
# - server_id: The ID of the server to which the HBA will be added (Required)
# - hbas: List of HBAs to add to the server (Required, Min: 1)
#   - hba_wwn: HBA World Wide Name (Optional)
#   - iscsi_name: iSCSI initiator name (Optional)
#
# ------------------------------
# Resource Attributes
# ------------------------------
# - id: The ID of this resource
# - server_hba_count: Total number of server HBAs
# - server_hba_info: List of server HBA information returned from API
#
# Note: Either hba_wwn or iscsi_name must be specified for each HBA block.
#

# Create a standard FC server (prerequisite)
resource "hitachi_vsp_one_server" "fc_server" {
  serial          = 12345
  server_nickname = "fc-server-001-temporary6"
  protocol        = "FC"
  os_type         = "Linux"
  is_reserved     = false
}

# Add HBA to the FC server
resource "hitachi_vsp_one_server_hba" "fc_server_hba" {
  serial    = 12345
  server_id = hitachi_vsp_one_server.fc_server.data[0].server_id
  
  hbas {
    hba_wwn = "500143802426b1c0"
  }
}

output "fc_server_hba_info" {
  description = "FC server HBA information"
  value = hitachi_vsp_one_server_hba.fc_server_hba.server_hba_info
}

# Create an iSCSI server (prerequisite)
resource "hitachi_vsp_one_server" "iscsi_server" {
  serial          = 12345
  server_nickname = "iscsi-server-001-temporary6"
  protocol        = "iSCSI"
  os_type         = "Windows"
  is_reserved     = false
}

# Add iSCSI initiator to the iSCSI server
resource "hitachi_vsp_one_server_hba" "iscsi_server_hba" {
  serial    = 12345
  server_id = hitachi_vsp_one_server.iscsi_server.data[0].server_id
  
  hbas {
    iscsi_name = "iqn.1991-05.com.microsoft:server01.example.com"
  }
}

output "iscsi_server_hba_info" {
  description = "iSCSI server HBA information"
  value = hitachi_vsp_one_server_hba.iscsi_server_hba.server_hba_info
}