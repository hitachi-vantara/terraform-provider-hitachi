data "hitachi_vsp_iscsi_chap_users" "my_iscsi_all_chap_users" {
  serial              = 30078
  port_id             = "CL4-C"
  iscsi_target_number = 1

}

output "my_iscsi_all_chap_users_output" {
  value = data.hitachi_vsp_iscsi_chap_users.my_iscsi_all_chap_users
}
