data "hitachi_vss_block_volumes" "vssbvolumes" {
  vss_block_address = "10.10.12.13"
}

output "volumeoutput" {
  value = data.hitachi_vss_block_volumes.vssbvolumes
}
