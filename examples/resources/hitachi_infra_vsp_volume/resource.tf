resource "hitachi_infra_vsp_volume" "newvol2" {
  serial = 40014
  name = "newVOlumeName1f"
  capacity = "200MB"
  pool_id = 4
  parity_group_id = "1-1"
  system ="Logical-UCP-30595"
  deduplication_compression_mode = "DISABLED"
}



output "volumesData" {
  value = resource.hitachi_infra_vsp_volume.newvol2
}

