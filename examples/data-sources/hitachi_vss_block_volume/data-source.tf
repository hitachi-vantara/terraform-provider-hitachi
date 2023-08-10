data "hitachi_vss_block_volume" "vssbvolumes" {
  vss_block_address = "10.10.12.13"
  volume_name       = "Mongonode3_vol4dd"

}

output "volumeoutput" {
  value = data.hitachi_vss_block_volume.vssbvolumes
}
