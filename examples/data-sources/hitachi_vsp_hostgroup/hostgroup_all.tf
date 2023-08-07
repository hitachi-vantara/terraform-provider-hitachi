data "hitachi_vsp_hostgroups" "hostgroups" {
  serial   = 40014
  port_ids = ["CL7-C", "CL7-A", "CL8-B", "CL9-C"]
}

output "hostgroups" {
  value = data.hitachi_vsp_hostgroups.hostgroups
}
