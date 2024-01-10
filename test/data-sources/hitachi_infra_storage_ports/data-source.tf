/*
data "hitachi_infra_storage_devices" "storage_devices" {
  serial = 611039
}
*/

data "hitachi_infra_storage_ports" "storage_ports" {
  serial = 611039
  #storage_id = data.hitachi_infra_storage_devices.storage_devices.id
  port_id = "CL8-B"
}

output "storage_ports" {
  value = data.hitachi_infra_storage_ports.storage_ports
}
