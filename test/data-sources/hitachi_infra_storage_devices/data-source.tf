data "hitachi_infra_storage_devices" "storage_devices" {
  serial = 40015
}

output "storage_devices" {
  value = data.hitachi_infra_storage_devices.storage_devices
}


