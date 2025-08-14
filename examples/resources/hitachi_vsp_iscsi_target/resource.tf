//
// Hitachi VSP iSCSI Target Resource
//
// This section defines a Terraform resource block to create a Hitachi VSP iSCSI target.
// The resource "hitachi_vsp_iscsi_target" represents an iSCSI target on a Hitachi
// Virtual Storage Platform (VSP) and allows you to manage its configuration using Terraform.
//
// Customize the values of the parameters (serial, iscsi_target_number, iscsi_target_alias,
// port_id, host_mode_options, host_mode) to match your desired iSCSI target configuration.
// The comments provide a link for detailed information about host_mode_options and host_mode.
//

resource "hitachi_vsp_iscsi_target" "myiscsi" {
  serial              = 12345
  iscsi_target_number = 1
  iscsi_target_alias  = "snewar-tgt1" 
  port_id             = "CL4-C"  

  // For detailed information about host_mode_options and host_mode, refer to:
  // https://docs.hitachivantara.com/r/en-us/svos/9.8.7/mk-97hm85026/managing-logical-volumes/configuring-hosts/host-modes-and-host-mode-options-for-host-facing-host-ports
  host_mode_options = [90]
  host_mode         = "VMware"
}

output "iscsioutput" {
  value = resource.hitachi_vsp_iscsi_target.myiscsi
}
