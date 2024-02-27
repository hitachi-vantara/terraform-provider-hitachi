
# data "hitachi_vsp_volume" "volume" {
#   serial  = 611039
#   ldev_id = 1
# }


# output "volume" {
#   value = data.hitachi_vsp_volume.volume
# }



data "hitachi_vsp_volumes" "volume1" {
  serial         = 611039
  start_ldev_id  = 0
  end_ldev_id    = 10
  undefined_ldev = false
  # subscriber_id = "46519299-c43c-4c6e-a680-81dce45a3fca"
  
}

output "volume1" {
  value = data.hitachi_vsp_volumes.volume1
}

