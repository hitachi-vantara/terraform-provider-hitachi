
# resource "hitachi_vsp_iscsi_target" "myiscsi" {
#   serial              = 30078 // REQUIRED
#   iscsi_target_number = 5
#   iscsi_target_alias   = "TESTING_REST_API_CL4-C5" // REQUIRED
#   port_id             = "CL4-C"                   // REQUIRED
#   // For detail information about host_mode_options and host_mode, please look at the following link:
#   // https://docs.hitachivantara.com/r/en-us/svos/9.8.7/mk-97hm85026/managing-logical-volumes/configuring-hosts/host-modes-and-host-mode-options-for-host-facing-host-ports
#   host_mode_options = [12, 32]
#   #host_mode_options = [11,12] 
#   host_mode = "AIX"
#   #host_mode ="LINUX/IRIX" 


#   lun {
#     ldev_id = 16
#     lun_id  = 12
#   }

#   lun {
#     ldev_id = 20
#     lun_id  = 13
#   }


#   initiator {
#     initiator_name     = "iqn.1994-05.com.redhat:496799baaa"
#     initiator_nickname = "test-iqn-nickname-1"
#   }
#   initiator {
#     initiator_name     = "iqn.1994-05.com.redhat:496799bbbb"
#     initiator_nickname = "test-iqn-nickname-2"
#   }

# }

