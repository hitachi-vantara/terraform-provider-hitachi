
data "hitachi_infra_volume" "volume" {
  serial  = 611032
  ldev_id = 32749
}

output "volume" {
  value = data.hitachi_infra_volume.volume
}


/*
data "hitachi_infra_volumes" "volume1" {
  serial         = 611032
  start_ldev_id  = 0
  end_ldev_id    = 32749
  undefined_ldev = false
}

output "volume1" {
  value = data.hitachi_infra_volumes.volume1
}
*/
