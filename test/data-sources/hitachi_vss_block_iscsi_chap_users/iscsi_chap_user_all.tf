data "hitachi_vss_block_iscsi_chap_users" "all_chap_users" {
  vss_block_address = var.vssb_address
}

output "all_chap_users_output" {
  value = data.hitachi_vss_block_iscsi_chap_users.all_chap_users
}
