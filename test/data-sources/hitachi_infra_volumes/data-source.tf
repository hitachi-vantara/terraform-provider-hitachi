
# data "hitachi_infra_volume" "volume" {
#   serial  = 40014
#   ldev_id = 550
# }

# output "volume" {
#   value = data.hitachi_infra_volume.volume
# }



data "hitachi_infra_volumes" "volume1" {
  serial         = 40014
  start_ldev_id  = 0
  end_ldev_id    = 10
  undefined_ldev = false
  subscriber_id = "46519299-c43c-4c6e-a680-81dce45a3fca"
  
}

output "volume1" {
  value = data.hitachi_infra_volumes.volume1
}

