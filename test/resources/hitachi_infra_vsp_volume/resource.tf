resource "hitachi_infra_vsp_volume" "newvol2" {
  serial = 40014
  name = "newVOlumeName1fg"
  capacity = "200MB"
  pool_id = 4
  parity_group_id = "1-1"
  system ="Logical-UCP-30595"
  deduplication_compression_mode = "DISABLED"
  subscriber_id = "46519299-c43c-4c6e-a680-81dce45a3fcb"
}



output "volumesData" {
  value = resource.hitachi_infra_vsp_volume.newvol2
}

