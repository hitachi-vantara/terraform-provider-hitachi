data "hitachi_vss_block_storage_ports" "storagePorts" {
  vss_block_address = "10.10.12.13"
  port_name = "001-iSCSI-002"
}

output "storagePorts" {
  value = data.hitachi_vss_block_storage_ports.storagePorts
}
