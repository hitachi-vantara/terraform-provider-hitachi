
data "hitachi_vsp_iscsi_chap_user" "my_iscsi_initiator_chap_user" {
  serial              = 30078
  port_id             = "CL4-C"
  iscsi_target_number = 1
  chap_user_type      = "initiator" # valid input value : "initiator", "target"
  chap_user_name      = "dan"
}

/*
data "hitachi_vsp_iscsi_chap_user" "my_iscsi_target_chap_user" {
   serial   = 30078
   port_id   = "CL4-C"
   iscsi_target_number   = 1
   chap_user_type = "target" # valid input value : "initiator", "target"
   chap_user_name = "rahul"
}
*/

output "my_iscsi_initiator_chap_user_output" {
  value = data.hitachi_vsp_iscsi_chap_user.my_iscsi_initiator_chap_user
}

/*
output "my_iscsi_target_chap_user_output" {
  value = data.hitachi_vsp_iscsi_chap_user.my_iscsi_target_chap_user
}
*/

data "hitachi_vsp_iscsi_chap_users" "my_iscsi_chap_users" {
  serial              = 30078
  port_id             = "CL4-C"
  iscsi_target_number = 1

}

output "my_iscsi_chap_users_output" {
  value = data.hitachi_vsp_iscsi_chap_users.my_iscsi_chap_users
}