# data "hitachi_vsp_volume" "volume" {
#   serial  = 611039
#   ldev_id = 4322

# }



# output "volume" {
#   value = data.hitachi_vsp_volume.volume
# }

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



data "hitachi_vsp_volumes" "volume1" {
  serial         = 611039
  # start_ldev_id  = 4000
  # end_ldev_id    = 5000
  undefined_ldev = true
}

output "volume1" {
  value = data.hitachi_vsp_volumes.volume1
}
