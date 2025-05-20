resource "hitachi_vsp_iscsi_target" "myiscsi" {
  serial              = 30078 // REQUIRED
  iscsi_target_number = 1
  iscsi_target_alias   = "snewar-tgt1" // REQUIRED
  port_id             = "CL4-C"       // REQUIRED

  // For detail information about host_mode_options and host_mode, please look at the following link:
  // https://docs.hitachivantara.com/r/en-us/svos/9.8.7/mk-97hm85026/managing-logical-volumes/configuring-hosts/host-modes-and-host-mode-options-for-host-facing-host-ports
  host_mode_options = [90]
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
