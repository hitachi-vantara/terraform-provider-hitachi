resource "hitachi_infra_storage_device" "demo_sd" {  
  serial = 611039
  management_address = "172.25.44.107"
  username = "maintenance"
  password = "raid-maintenance"
  # gateway_address = "172.25.22.81"
  out_of_band = true
  # system = "UCP-SYS2"
 
}
