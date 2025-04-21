data "hitachi_vsp_volume" "volume" {
  serial  = 40014
  ldev_id = 281
}

output "volume" {
  value = data.hitachi_vsp_volume.volume
}

# data "hitachi_vsp_volume" "lun2" {
#   serial  = 611039
#   ldev_id = 2
# }

# data "hitachi_vsp_volume" "lun78" {
#   serial  = 40014
#   ldev_id = 78
# }

# output "lun2" {
#   value = data.hitachi_vsp_volume.lun2
# }

# output "lun78" {
#   value = data.hitachi_vsp_volume.lun78
# }
