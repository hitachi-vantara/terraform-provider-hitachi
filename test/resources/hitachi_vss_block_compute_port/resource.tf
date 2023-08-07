resource "hitachi_vss_block_compute_port" "mycomputeport" {
  vss_block_address = "10.76.47.55"
  name = "001-iSCSI-002"
  authentication_settings = "None"
  target_chap_users = []
  #target_chap_users = []
 
}

output "chapuser_association_with_computeport" {
  value = resource.hitachi_vss_block_compute_port.mycomputeport
}
