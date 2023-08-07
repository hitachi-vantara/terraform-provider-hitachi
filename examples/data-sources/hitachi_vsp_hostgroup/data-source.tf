
data "hitachi_vsp_hostgroup" "hostgroup" {
  serial           = 40014
  port_id          = "CL1-A"
  hostgroup_number = 10
}

output "hostgroup" {
  value = data.hitachi_vsp_hostgroup.hostgroup
}

data "hitachi_vsp_hostgroups" "hostgroups" {
  serial   = 40014
  port_ids = ["CL7-C", "CL7-A", "CL8-B", "CL9-C"]
}

output "hostgroups" {
  value = data.hitachi_vsp_hostgroups.hostgroups
}
