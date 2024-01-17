data "hitachi_infra_storage_devices" "storage_devices" {
  serial = 611039
}

data "hitachi_infra_chap_users" "chap_users" {
  serial = 611039
  #storage_id = data.hitachi_infra_storage_devices.storage_devices.id
  port_id = "CL4-C"
  iscsi_target_number = 251
}

output "chap_users" {
  value = data.hitachi_infra_chap_users.chap_users

}

/*
data "hitachi_infra_iscsi_target" "iscsi_target" {
  serial = 611039
  #storage_id = data.hitachi_infra_storage_devices.storage_devices.id
  port_id = "CL4-C"
  #iscsi_target_number = 234
  iscsi_name = "Auto-iSCSITarget-394414"
}

output "iscsi_target" {
  value = data.hitachi_infra_iscsi_target.iscsi_target
}
*/
