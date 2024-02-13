resource "hitachi_infra_storage_device" "demo_sd" {  
  serial = 40014
  management_address = "172.25.47.115"
  username = "ms_vmware"
  password = "Hitachi1"
  gateway_address = "172.25.20.35"
  #ucp_system = "UCP-CI-12035"
}
