data "hitachi_vsp_iscsi_target" "iscsitarget" {
  serial              = 30078
  port_id             = "CL4-C"
  iscsi_target_number = 1
}

output "iscsitarget" {
  value = data.hitachi_vsp_iscsi_target.iscsitarget
}
data "hitachi_vsp_iscsi_targets" "iscsitargets" {
  serial   = 30078
  port_ids = ["CL5-A", "CL4-C", "CL7-B"]
}

output "iscsitargets" {
  value = data.hitachi_vsp_iscsi_targets.iscsitargets
}
