data "hitachi_vsp_volume" "volume" {
  serial  = 611039
  ldev_id = 1995

}

output "volume" {
  value = data.hitachi_vsp_volume.volume
}

# data "hitachi_vsp_volume" "lun2" {
#   serial  = 40015
#   ldev_id = 2
# }

# # data "hitachi_vsp_volume" "lun78" {
# #   serial  = 40014
# #   ldev_id = 78
# # }

# output "lun2" {
#   value = data.hitachi_vsp_volume.lun2
# }

# # output "lun78" {
# #   value = data.hitachi_vsp_volume.lun78
# # }



# # data "hitachi_vsp_volumes" "volume1" {
# #   serial         = 40014
# #   start_ldev_id  = 280
# #   end_ldev_id    = 285
# #   undefined_ldev = false
# # }

# # output "volume1" {
# #   value = data.hitachi_vsp_volumes.volume1
# # }
