data "hitachi_vss_block_iscsi_chap_users" "my_chap_users" {
  vss_block_address = ""
}

output "my_iscsi_chap_users_output" {
  value = data.hitachi_vss_block_iscsi_chap_users.my_chap_users
}
