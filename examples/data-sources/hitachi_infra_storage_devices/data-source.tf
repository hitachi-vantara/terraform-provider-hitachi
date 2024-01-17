data "hitachi_infra_storage_devices" "storage_devices" {
  #serial = 611039
}

output "storage_devices" {
  value = data.hitachi_infra_storage_devices.storage_devices
}
 