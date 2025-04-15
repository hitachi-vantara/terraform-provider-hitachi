
/*
data "hitachi_vosb_iscsi_chap_users" "chap_user_by_id" {
   vosb_address   = ""
   target_chap_user = "a79c1a1d-2719-4e07-b800-faf9de73d0ae"
}

output "id_output" {
  value = data.hitachi_vosb_iscsi_chap_users.chap_user_by_id
}

*/

data "hitachi_vosb_iscsi_chap_users" "chap_user_by_name" {
  vosb_address = var.vosb_address
  #target_chap_user = "rahul70"
}

output "name_output" {
  value = data.hitachi_vosb_iscsi_chap_users.chap_user_by_name
}

data "hitachi_vosb_iscsi_chap_users" "my_chap_users" {
  vosb_address = var.vosb_address
}

output "my_iscsi_chap_users_output" {
  value = data.hitachi_vosb_iscsi_chap_users.my_chap_users
}
