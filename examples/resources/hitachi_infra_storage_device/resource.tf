resource "hitachi_infra_storage_device" "demo_sd" {  
  serial = 6110340
  management_address = "172.25.44.107"
  username = "maintenance"
  password = "raid-maintenance"
  gateway_address = "172.25.20.35"

}

