resource "hitachi_vss_block_iscsi_chap_user" "my_chap_user" {
  vss_block_address       = ""
  target_chap_user_name   = "targetchapuser"
  target_chap_user_secret = "targetchapuserpasswd"
}

output "chap_user_output" {
  value = resource.hitachi_vss_block_iscsi_chap_user.my_chap_user
}
