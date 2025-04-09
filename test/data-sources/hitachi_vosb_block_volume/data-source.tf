data "hitachi_vosb_block_volume" "vosbvolumes" {
  vosb_block_address = ""
  volume_name       = "Mongonode3_vol4dd"

}

output "volumeoutput" {
  value = data.hitachi_vosb_block_volume.vosbvolumes
}
