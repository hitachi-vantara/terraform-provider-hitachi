data "hitachi_vsp_volumes" "volume1" {
  serial         = 40014
  start_ldev_id  = 280
  end_ldev_id    = 285
  undefined_ldev = false
}

output "volume1" {
  value = data.hitachi_vsp_volumes.volume1
}
