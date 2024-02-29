data "hitachi_vss_block_volume" "vssbvolumes" {
  vss_block_address = "172.25.58.151"
  volume_name       = "Mongonode3_vol4dd"

}

output "volumeoutput" {
  value = data.hitachi_vss_block_volume.vssbvolumes
}


data "hitachi_vsp_volumes" "volume1" {
  serial         = 810046
  start_ldev_id  = 0
  end_ldev_id    = 10
  undefined_ldev = false
  # subscriber_id = "46519299-c43c-4c6e-a680-81dce45a3fca"
  
}

output "volume1" {
  value = data.hitachi_vsp_volumes.volume1
}
