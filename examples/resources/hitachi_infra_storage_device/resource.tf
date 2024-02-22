resource "hitachi_infra_storage_device" "demo_sd" {  
  serial = 6110340
  management_address = "172.25.44.107"
  username = "maintenance"
  password = "raid-maintenance"
  gateway_address = "172.25.20.35"
  #system = "UCP-CI-12035"

  // Increase the timeout value for create/update operations accordingly
    timeouts {
    create = "10m"
    update = "10m"
  }
}
