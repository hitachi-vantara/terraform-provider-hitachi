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


resource "hitachi_vsp_hostgroup" "myhg" {
  serial            = 12345
  hostgroup_number  = 23
  hostgroup_name    = "TESTING-HOSTGROUP"
  port_id           = "CL1-A"
  host_mode_options = [12, 32]
  host_mode         = "AIX"

  # SET of LUN
  lun {
    ldev_id = 25
    lun     = 12
  }
}
