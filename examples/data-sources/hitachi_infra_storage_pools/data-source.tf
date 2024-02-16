
data "hitachi_infra_storage_devices" "storage_devices" {
  #storage_id = "storage-9a3f87a8c9dc213e8ebd02b63b97b9e8"
  serial = 611039
}

data "hitachi_infra_storage_pools" "storage_pools" {
  #storage_id = data.hitachi_infra_storage_devices.storage_devices.id
  serial = 611039
  #pool_name = "AutoPoolTag-291246"
  pool_id = 51
}

output "storage_pools" {
  value = data.hitachi_infra_storage_pools.storage_pools
}
