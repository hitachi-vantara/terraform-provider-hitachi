
data "hitachi_vsp_iscsi_chap_user" "my_iscsi_initiator_chap_user" {
  serial              = 30078
  port_id             = "CL4-C"
  iscsi_target_number = 1
  chap_user_type      = "initiator" 
  chap_user_name      = "chapuser1"

}

output "my_iscsi_initiator_chap_user_output" {
  value = data.hitachi_vsp_iscsi_chap_user.my_iscsi_initiator_chap_user
}

data "hitachi_vsp_iscsi_chap_users" "my_iscsi_chap_users" {
  serial              = 30078
  port_id             = "CL4-C"
  iscsi_target_number = 1

}

output "my_iscsi_chap_users_output" {
  value = data.hitachi_vsp_iscsi_chap_users.my_iscsi_chap_users
}