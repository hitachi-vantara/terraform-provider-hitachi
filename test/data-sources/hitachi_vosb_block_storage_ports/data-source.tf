data "hitachi_vosb_block_storage_ports" "storagePorts" {
  vosb_block_address = "10.76.47.55"
  port_name = "001-iSCSI-002"
}

output "storagePorts" {
  value = data.hitachi_vosb_block_storage_ports.storagePorts
}
