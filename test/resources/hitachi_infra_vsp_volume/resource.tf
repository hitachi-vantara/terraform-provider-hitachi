resource "hitachi_infra_vsp_volume" "newvol2" {
  serial = 40014
  name = "heyVolTest"
  capacity = "200MB"
  pool_id = 4
  parity_group_id = "1-1"
  system ="Logical-UCP-30595"
  # subscriber_id = "46519299-c43c-4c6e-a680-81dce45a3fcb"


   timeouts {
    create = "10m"
    update = "10m"
  }
}



# resource "hitachi_infra_vsp_volume" "newvol2" {
#   serial = 40014
#   name = "newVOlumeName1fgddd"
#   capacity = "100MB"
#   pool_id = 4
#   parity_group_id = "1-1"
#   system ="Logical-UCP-30595"
#   # subscriber_id = "46519299-c43c-4c6e-a680-81dce45a3fcb"
# }

# resource "hitachi_infra_vsp_volume" "newvol22" {
#   serial = 40014
#   name = "newvolumecretare11"
#   capacity = "100MB"
#   pool_id = 4
#   parity_group_id = "1-1"
#   system ="Logical-UCP-30595"

# }


# output "volumesData" {
#   value = resource.hitachi_infra_vsp_volume.newvol2
# }

