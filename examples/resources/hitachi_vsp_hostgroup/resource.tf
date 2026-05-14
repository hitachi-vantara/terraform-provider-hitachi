#
# Hitachi VSP Hostgroup Resource
#
# This section defines a Terraform resource block to create a Hitachi VSP hostgroup.
# The resource "hitachi_vsp_hostgroup" represents a hostgroup on a Hitachi Virtual
# Storage Platform (VSP) and allows you to manage its configuration using Terraform.
#
# Customize the values of the parameters (serial, hostgroup_number, hostgroup_name,
# port_id, host_mode_options, host_mode) and the nested "lun" block to match your
# desired hostgroup configuration.
#
# For detailed information about host_mode_options and host_mode, please refer to
# the official Hitachi documentation:
# https://docs.hitachivantara.com/r/en-us/svos/9.8.7/mk-97hm85026/managing-logical-volumes/configuring-hosts/host-modes-and-host-mode-options-for-host-facing-host-ports
#

################################################################################
# Example: Import (existing hostgroup)
# -----------------------------------------------------------------------------
# Use terraform import when the hostgroup already exists on storage and
# you want to bring it under Terraform management without re-creating it.
# The import reads the current state from storage and writes it to tfstate.
#
# Import ID formats:
# - <serial>/<port_id>,<hostgroup_name>
# - <serial>/<port_id>,<hostgroup_number>,<hostgroup_name>  (use when name is not unique on the port)
#
# terraform import hitachi_vsp_hostgroup.imported '12345/CL1-A,TESTING-HOSTGROUP'
# terraform import hitachi_vsp_hostgroup.imported '12345/CL1-A,23,TESTING-HOSTGROUP'

# Minimal skeleton block required for `terraform import`.
# Create this block in your .tf before running the import command.
resource "hitachi_vsp_hostgroup" "imported" {}

output "imported_hostgroup_id" {
  value = hitachi_vsp_hostgroup.imported.id
}
################################################################################


resource "hitachi_vsp_hostgroup" "myhg" {
  serial           = 12345
  hostgroup_number = 23 # optional
  hostgroup_name   = "TESTING-HOSTGROUP"
  port_id          = "CL1-A"

  # optional
  host_mode_options = [12, 32]
  host_mode         = "AIX"

  # optional
  lun {
    ldev_id = 25
    // or
    // ldev_id_hex = "0X19"
    lun = 12
  }

  # optional
  wwn {
    host_wwn     = "1200000012000001"
    wwn_nickname = "NAME-12012001"
  }
}

output "hgoutput" {
  value = resource.hitachi_vsp_hostgroup.myhg
}
