
# data "hitachi_infra_volume" "volume" {
#   serial  = 40014
#   ldev_id = 0
# }

# output "volume" {
#   value = data.hitachi_infra_volume.volume
# }



data "hitachi_infra_volumes" "volume1" {
  serial         = 40014
  start_ldev_id  = 0
  end_ldev_id    = 10
  undefined_ldev = false
}

output "volume1" {
  value = data.hitachi_infra_volumes.volume1
}

