resource "hitachi_infra_storage_device" "demo_sd" {  
  serial = 715021
  management_address = "172.25.44.116"
  username = "maintenance"
  password = "raid-maintenance"
  gateway_address = "172.25.20.56"
  # out_of_band = false
  # system = "UCP-SYS2"
}
