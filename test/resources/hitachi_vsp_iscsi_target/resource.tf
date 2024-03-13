resource "hitachi_vsp_iscsi_target" "myiscsi" {
  serial              = 611039 // REQUIRED
  # iscsi_target_number = 1
  iscsi_target_alias   = "NewTarget" // REQUIRED
  port_id             = "CL2-C"       // REQUIRED

  // For detail information about host_mode_options and host_mode, please look at the following link:
  // https://knowledge.hitachivantara.com/Documents/Management_Software/SVOS/9.8.6/Volume_Management_-_VSP_E_Series/Host_Attachment/14_Host_modes_and_host_mode_options
  # host_mode_options = [90]
  host_mode         = "VMware"
  #iscsi_target_name = "iqn.1995-05.com.redhat:496799ba71"
  #host_mode ="LINUX/IRIX" 
  #lun {
  #    ldev_id = 21
  #    lun_id = 42
  #}
  #initiator {
  #      initiator_name = "iqn.1995-05.com.redhat:496799ba72"
  #      initiator_nickname = "test-iqn-nickname-301"
  #}


}
