data "hitachi_vsp_iscsi_targets" "alliscsitargets" {
  serial   = 30078
  port_ids = ["CL5-A", "CL4-C", "CL7-B"]
}

output "alliscsitargets" {
  value = data.hitachi_vsp_iscsi_targets.alliscsitargets
}
