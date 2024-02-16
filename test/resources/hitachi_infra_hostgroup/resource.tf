resource "hitachi_infra_hostgroup" "demo_fc_hg" {  
  serial = 611039
  hostgroup_number = 02  
  hostgroup_name = "demo-hostgroup-k"   
  port_id = "CL1-A"    
  host_mode_options = [54]   
  host_mode ="vmware extension"
  wwn {    
    host_wwn = "200000109B95416F"    
    wwn_nickname = "esxi-4"  
  }
  lun {    
    ldev_id = 294    
    lun     = 14  
  }
  system = "UCP-CI-11234"
}

/*
resource "hitachi_infra_hostgroup" "demo_fc_hg1" {  
  serial = 30595    
  hostgroup_number = 03  
  hostgroup_name = "demo-hostgroup1"   
  port_id = "CL7-A"    
  host_mode_options = [54]   
  host_mode ="vmware extension"
  wwn {    
    host_wwn = "200000109B95414F"    
    wwn_nickname = "esxi-2"  
  }
  lun {    
    ldev_id = 294    
    lun     = 12  
  }
}
*/




/*

resource "hitachi_vsp_hostgroup" "myhg" {
  serial           = 40014 // REQUIRED
  hostgroup_number = 36
  hostgroup_name   = "TESTING-HOSTGROUP112ki" // REQUIRED
  port_id          = "CL1-A"             // REQUIRED
  // For detail information about host_mode_options and host_mode, please look at the following link:
  // https://knowledge.hitachivantara.com/Documents/Management_Software/SVOS/9.8.6/Volume_Management_-_VSP_E_Series/Host_Attachment/14_Host_modes_and_host_mode_options
  host_mode_options = []
  host_mode         = "VMware Extension"
  #host_mode_options = [11] 
  #host_mode ="LINUX/IRIX" 


  # SET of LUN

  lun {
    # ldev_id = 25
    lun     = 25
  }
*/
  /*
  lun {
      ldev_id = 21
      lun = 13 
  }
  */

  # SET of WWN
  /*
#For Create
  wwn {
        host_wwn = "100000109b3dfbbb"
        wwn_nickname = "test-wwn1"
  }
  wwn {
        host_wwn = "100000109b3dfaaa"
        wwn_nickname = "test-wwn2"
  }
*/
  /*
# For Update
  wwn {
        host_wwn = "100000109b3dfbbb"
        wwn_nickname = "test-wwn1b"
  }
  wwn {
        host_wwn = "100000109b3dfbbc"
        wwn_nickname = "test-wwn1c"
  }
  */
#}


# terraform destroy -target hitachi_storage_hostgroup.myhg
# terraform apply -target=hitachi_storage_hostgroup.myhg

/*
First:
  lun {
      ldev_id = 25
      lun = 12 
  }
Second
  lun {
      ldev_id = 21
      lun = 13 
  }
  #  Delete 25,12 and add 21, 13
====
First:
  lun {
      ldev_id = 25
      lun = 12 
  }
Second
  lun {
      ldev_id = 25
      lun = 13
  }
  #  Delete 25, 12 and add 25, 13
===
First:
  lun {
      ldev_id = 25
      lun = 12 
  }
Second
  lun {
      ldev_id = 25
      lun = 12 
  }
  lun {
      ldev_id = 21
      lun = 13
  }
  #  Only Add 21, 13
====
First:
  lun {
      ldev_id = 25
      lun = 12 
  }
Second
  lun {
      ldev_id = 25
      lun = 12
  }
 lun {
      ldev_id = 21
      lun = 12
  }
  #  Error as Lun Id same
====
First:
  lun {
      ldev_id = 25
  }
  #  Automatic add available lun id

====
First:
  lun {
      ldev_id = 25
      lun = 12 
  }
Second
  lun {
      ldev_id = 25
  }
  #  Delete 25,12 and Automatic add available lun id
*/
