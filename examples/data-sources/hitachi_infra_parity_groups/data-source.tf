data "hitachi_infra_storage_devices" "storage_devices" {
  #storage_id = "storage-9a3f87a8c9dc213e8ebd02b63b97b9e8"
  serial = 611039
}

data "hitachi_infra_parity_groups" "parity_groups" {
  storage_id = data.hitachi_infra_storage_devices.storage_devices.id
  #serial = 611039
  parity_group_ids = [ "E1-1", "E1-2"]
}

output "parity_groups" {
  value = data.hitachi_infra_parity_groups.parity_groups
}
