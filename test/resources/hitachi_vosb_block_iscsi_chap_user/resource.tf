resource "hitachi_vosb_block_iscsi_chap_user" "my_chap_user10" {
  vosb_block_address       = "10.76.47.55"
  target_chap_user_name   = "rahul1110"
  #target_chap_user_secret = "rahul80"
  target_chap_user_secret = "rahul80101112112"
  #initiator_chap_user_name = "rahul111"
  #initiator_chap_user_secret = "rahul8910111222"
}

output "chap_user_output2" {
  value = resource.hitachi_vosb_block_iscsi_chap_user.my_chap_user10
}

resource "hitachi_vosb_block_iscsi_chap_user" "my_chap_user9" {
  vosb_block_address       = "10.76.47.55"
  target_chap_user_name   = "rahul1198902"
  #target_chap_user_secret = "rahul80"
  target_chap_user_secret = "rahul80101112112845690"
  #initiator_chap_user_name = "rahul111"
  #initiator_chap_user_secret = "rahul8910111222"
}

output "chap_user_output" {
  value = resource.hitachi_vosb_block_iscsi_chap_user.my_chap_user9
}
