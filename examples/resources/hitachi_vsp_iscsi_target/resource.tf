#
# Hitachi VSP iSCSI Target Resource
#
# This section defines a Terraform resource block to create a Hitachi VSP iSCSI target.
# The resource "hitachi_vsp_iscsi_target" represents an iSCSI target on a Hitachi
# Virtual Storage Platform (VSP) and allows you to manage its configuration, 
# LUN mappings, and initiators using Terraform.
#
# Customize the parameters (serial, target number, port, host mode) and the nested 
# "lun" and "initiator" blocks to match your storage networking requirements.
#
# For detailed information about host_mode_options and host_mode, refer to:
# https://docs.hitachivantara.com/r/en-us/svos/9.8.7/mk-97hm85026/managing-logical-volumes/configuring-hosts/host-modes-and-host-mode-options-for-host-facing-host-ports
#

resource "hitachi_vsp_iscsi_target" "myiscsi" {
  serial              = 12345
  iscsi_target_number = 10                         # Optional: Specific target ID
  iscsi_target_alias  = "SV-TEST-alias-10-tgt10"   # The iSCSI target name/alias
  port_id             = "CL4-E"                    # The physical port ID on the VSP

  # optional
  host_mode_options = [90]
  host_mode         = "VMware"

  # The lun block maps a Logical Device (LDEV) to a LUN ID accessible via this target
  # You can repeat this block for multiple LUN mappings
  lun {
    ldev_id = 66    # The internal storage LDEV number
    lun_id  = 5     # The LUN number seen by the host (HAVE)
  }

  # The initiator block defines the authorized host IQN allowed to access this target
  # Similar to WWN filtering in Fibre Channel
  initiator {
    initiator_name     = "iqn.1994-04.jp.co.logicl:rsd.t9v.i.00204.1v001"
    initiator_nickname = "alias-iscsi-init-name-1"
  }
}

# Output the entire resource state for verification
output "iscsioutput" {
  value = resource.hitachi_vsp_iscsi_target.myiscsi
}