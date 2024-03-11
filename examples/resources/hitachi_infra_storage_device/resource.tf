resource "hitachi_infra_storage_device" "demo_sd" {  
  serial = 40015
  management_address = "172.25.47.116"
  username = "maintenance"
  password = "raid-maintenance"
  gateway_address = "172.25.20.35"

}

