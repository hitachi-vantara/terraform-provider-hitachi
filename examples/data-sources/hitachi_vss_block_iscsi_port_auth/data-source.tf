data "hitachi_vss_block_iscsi_port_auth" "mycomputeport" {
  vss_block_address = "10.10.12.13"
  name = "001-iSCSI-002"
}

output "mycomputeport" {
  value = data.hitachi_vss_block_iscsi_port_auth.mycomputeport
}
