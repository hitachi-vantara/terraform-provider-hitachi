
resource "hitachi_vsp_hostgroup" "myhg" {
  serial           = 40014 // REQUIRED
  hostgroup_number = 23
  hostgroup_name   = "TESTING-HOSTGROUP" // REQUIRED
  port_id          = "CL1-A"             // REQUIRED
  // For detail information about host_mode_options and host_mode, please look at the following link:
  // https://knowledge.hitachivantara.com/Documents/Management_Software/SVOS/9.8.6/Volume_Management_-_VSP_E_Series/Host_Attachment/14_Host_modes_and_host_mode_options
  host_mode_options = [12, 32]
  host_mode         = "AIX"
 
  # SET of LUN
  lun {
    ldev_id = 25
    lun     = 12
  }

}