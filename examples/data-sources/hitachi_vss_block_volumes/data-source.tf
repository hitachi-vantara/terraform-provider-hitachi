data "hitachi_vss_block_volumes" "vssbvolumes" {
  vss_block_address = ""
  #compute_node_name = "esxi-151" // Optional
}

output "volumeoutput" {
  value = data.hitachi_vss_block_volumes.vssbvolumes
}
