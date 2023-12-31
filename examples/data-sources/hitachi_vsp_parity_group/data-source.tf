data "hitachi_vsp_parity_groups" "myparitygroup" {
  serial = 12345
  parity_group_ids = ["1-2","1-3"]
}

output "myparitygroup" {
  value = data.hitachi_vsp_parity_groups.myparitygroup
}
