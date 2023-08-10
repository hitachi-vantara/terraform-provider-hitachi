resource "hitachi_vss_block_compute_port" "mycomputeport" {
  vss_block_address = "10.10.12.13"
  name = "001-iSCSI-002"
  authentication_settings = "CHAP"
  target_chap_users = ["test_user7", "test_user9"]
 
}

output "chapuser_association_with_computeport" {
  value = resource.hitachi_vss_block_compute_port.mycomputeport
}
