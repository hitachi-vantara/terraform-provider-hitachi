resource "hitachi_infra_storage_device" "demo_sd" {  
  serial = 30595
  management_address = "172.25.47.112"
  username = "ms_vmware"
  password = "Hitachi1"
  gateway_address = "172.25.20.35"
  out_of_band = false
  #system = "UCP-CI-12035"
}
