data "hitachi_infra_storage_devices" "storage_devices" {
  # serial = 40015
  serial = 715036
}

output "storage_devices" {
  value = data.hitachi_infra_storage_devices.storage_devices
}


