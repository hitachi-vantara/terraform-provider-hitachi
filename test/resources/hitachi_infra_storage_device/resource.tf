resource "hitachi_infra_storage_device" "demo_sd" {  
  serial = 40014
  management_address = "172.25.47.115"
  username = "ucpadmin"
  password = "Passw0rd!"
  gateway_address = "10.76.47.78"
  # out_of_band = false
  #system = "UCP-CI-12035"
}
