resource "hitachi_vsp_volume" "newvol2" {
  storage_id = "storage-e51aa8e9806a70a036a77fec150d1407"
  # serial = 715036
  name = "NewVOlTest111ggggg"
  size_gb = 1 
  pool_id = 0
  # system ="Logical-UCP-611"
  # subscriber_id = "46519299-c43c-4c6e-a680-81dce45a3fcb"

}


# resource "hitachi_vsp_volume" "newvol23" {
#   serial = 611039
#   name = "myvol222"
#   size_gb = 1 
#   pool_id = 0
#   system ="Logical-UCP-611039"
#   # subscriber_id = "46519299-c43c-4c6e-a680-81dce45a3fcb"

# }



output "volumesData" {
  value = resource.hitachi_vsp_volume.newvol2
}


