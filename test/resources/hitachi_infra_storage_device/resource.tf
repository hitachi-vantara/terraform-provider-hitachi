resource "hitachi_infra_storage_device" "demo_sd" {  
  serial = 611035
  management_address = "172.25.44.107"
  username = "maintenance"
  password = "raid-maintenance"
  gateway_address = "172.25.20.35"
  #ucp_system = "UCP-CI-12035"
}
