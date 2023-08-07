resource "hitachi_vsp_iscsi_target" "myiscsi" {
  serial              = 30078 // REQUIRED
  iscsi_target_number = 1
  iscsi_target_alias   = "snewar-tgt1" // REQUIRED
  port_id             = "CL4-C"       // REQUIRED

  // For detail information about host_mode_options and host_mode, please look at the following link:
  // https://knowledge.hitachivantara.com/Documents/Management_Software/SVOS/9.8.6/Volume_Management_-_VSP_E_Series/Host_Attachment/14_Host_modes_and_host_mode_options
  host_mode_options = [90]
  host_mode         = "VMware"



}
