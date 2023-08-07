data "hitachi_vss_block_storage_ports" "storagePorts" {
  vss_block_address = "10.76.47.55"
  #port_id = "5f07176a-e10d-47b7-99b4-57b93806048b"
}

output "storagePorts" {
  value = data.hitachi_vss_block_storage_ports.storagePorts
}
