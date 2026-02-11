#
# Hitachi VSP One Port Management (Update)
#
# This section defines a resource block for managing ports
# in the Hitachi Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_port" resource supports updating port configurations,
# allowing you to modify port settings such as security configurations,
# speed settings, and protocol-specific parameters through Terraform.
#
# ------------------------------
# About Update Operations
# ------------------------------
# The update operation modifies an existing port's configuration.
# Port resources are identified by their port ID and are updated in-place.
#
# During update:
#   - "port_id" is automatically set from the resource ID.
#   - "serial" identifies the storage system.
#   - Various port parameters can be updated including:
#       • port_security (enable/disable port security)
#       • port_speed (modify port speed settings)
#       • Protocol-specific configurations (FC, iSCSI, NVMe-TCP)
#
# ------------------------------
# Key configuration parameters
# ------------------------------
#   - serial: Specifies the storage system's serial number.
#   - port_id: ID of the port to update.
#   - port_security: Enables or disables port security (boolean).
#   - port_speed: Sets the port speed configuration.
#   - fc_information: FC-specific port settings (for FC ports).
#   - iscsi_information: iSCSI-specific port settings (for iSCSI ports).
#   - nvme_tcp_information: NVMe-TCP-specific port settings (for NVMe-TCP ports).
#
# The resource outputs updated port information including all current
# port attributes and configurations.
#
# Adjust the parameters to match your environment and desired configuration.
#

# -------------------------------------
# Example: Update Port Security
# -------------------------------------
resource "hitachi_vsp_one_port" "port_security_update" {
  serial  = var.serial_number
  port_id = "CX-X"

  # Enable port security for the specified port
  port_security = true

  # Optional: Set port speed if needed
  # port_speed = "NUMBER_32"

  # Optional: Configure FC information (for FC ports)
  # fc_information {
  #   fabric_switch_setting = true
  #   connection_type       = "Point_To_Point"
  #   al_pa                 = "D2"
  # }

  # Optional: Configure iSCSI information (for iSCSI ports)
  # iscsi_information {
  #   ip_mode     = "ipv4"
  #   delayed_ack = true
  #   add_vlan_id = 100
  #   
  #   ipv4_information {
  #     address         = "192.168.1.100"
  #     subnet_mask     = "255.255.255.0"
  #     default_gateway = "192.168.1.1"
  #   }
  # }

  # Optional: Configure NVMe-TCP information (for NVMe-TCP ports)
  # nvme_tcp_information {
  #   ip_mode     = "ipv4"
  #   add_vlan_id = 200
  #   
  #   ipv4_information {
  #     address         = "192.168.2.100"
  #     subnet_mask     = "255.255.255.0"
  #     default_gateway = "192.168.2.1"
  #   }
  # }
}

output "port_output" {
  value       = hitachi_vsp_one_port.port_security_update.port_info
  description = "Updated port information including security settings"
}
