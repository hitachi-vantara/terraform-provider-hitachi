data "hitachi_infra_storage_devices" "storage_devices" {
  serial = 611032
}

/*
data "hitachi_infra_iscsi_targets" "iscsi_targets" {
  serial = 611039
  #storage_id = data.hitachi_infra_storage_devices.storage_devices.id
  port_ids = ["CL4-C"]
}

output "iscsi_targets" {
  value = data.hitachi_infra_iscsi_targets.iscsi_targets
}
*/


data "hitachi_infra_iscsi_target" "iscsi_target" {
  serial = 611032
  #storage_id = data.hitachi_infra_storage_devices.storage_devices.id
  port_id = "CL4-C"
  #iscsi_target_number = 234
  #iscsi_name = "Auto-iSCSITarget-394414"
}

output "iscsi_target" {
  value = data.hitachi_infra_iscsi_target.iscsi_target
}

