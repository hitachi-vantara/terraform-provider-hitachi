resource "hitachi_vss_block_iscsi_chap_user" "my_chap_user2" {
  vss_block_address       = "10.76.47.55"
  target_chap_user_name   = "rahul112"
  target_chap_user_secret = "rahul8010111211"
  #initiator_chap_user_name = "rahul111"
  #initiator_chap_user_secret = "rahul8910111222"
}

output "chap_user_output" {
  value = resource.hitachi_vss_block_iscsi_chap_user.my_chap_user2
}
