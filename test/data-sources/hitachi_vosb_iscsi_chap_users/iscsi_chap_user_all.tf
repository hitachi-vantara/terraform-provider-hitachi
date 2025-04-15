data "hitachi_vosb_iscsi_chap_users" "all_chap_users" {
  vosb_address = var.vosb_address
}

output "all_chap_users_output" {
  value = data.hitachi_vosb_iscsi_chap_users.all_chap_users
}
