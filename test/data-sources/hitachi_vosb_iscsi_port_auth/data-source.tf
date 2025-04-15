data "hitachi_vosb_iscsi_port_auth" "mycomputeport" {
  vosb_address = "10.76.47.55"
  name = "001-iSCSI-002"
}

output "mycomputeport" {
  value = data.hitachi_vosb_iscsi_port_auth.mycomputeport
}
