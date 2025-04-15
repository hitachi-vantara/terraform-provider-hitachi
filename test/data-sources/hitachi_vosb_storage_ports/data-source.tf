data "hitachi_vosb_storage_ports" "storagePorts" {
  vosb_address = "10.76.47.55"
  port_name = "001-iSCSI-002"
}

output "storagePorts" {
  value = data.hitachi_vosb_storage_ports.storagePorts
}
