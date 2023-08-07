data "hitachi_vss_block_iscsi_port_auth" "mycomputeport" {
  vss_block_address = "10.76.47.55"
  name = "001-iSCSI-002"
  #authentication_settings = "CHAP"
  #target_chap_users = ["rahul", "rahul50"]
 
}

output "mycomputeport" {
  value = data.hitachi_vss_block_iscsi_port_auth.mycomputeport
}
