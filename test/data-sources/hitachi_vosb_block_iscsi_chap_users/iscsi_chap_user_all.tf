data "hitachi_vosb_block_iscsi_chap_users" "all_chap_users" {
  vosb_block_address = var.vosb_block_address
}

output "all_chap_users_output" {
  value = data.hitachi_vosb_block_iscsi_chap_users.all_chap_users
}
