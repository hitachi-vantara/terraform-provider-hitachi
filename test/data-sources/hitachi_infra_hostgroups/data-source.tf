data "hitachi_infra_storage_devices" "storage_devices" {
  #storage_id = "storage-9a3f87a8c9dc213e8ebd02b63b97b9e8"
  serial = 611032
}

data "hitachi_infra_hostgroup" "host_group" {
  #serial = 611039
  storage_id = data.hitachi_infra_storage_devices.storage_devices.id
  port_id = "CL8-B"
  hostgroup_number = 14
  #hostgroup_name = "PODB-ESXi-218-HBA1"
}

output "host_group" {
  value = data.hitachi_infra_hostgroup.host_group
}



data "hitachi_infra_hostgroups" "host_groups" {
  serial = 611032
  #storage_id = data.hitachi_infra_storage_devices.storage_devices.id
  #port_ids = ["CL7-A", "CL7-B"]
}

output "host_groups" {
  value = data.hitachi_infra_hostgroups.host_groups
}
