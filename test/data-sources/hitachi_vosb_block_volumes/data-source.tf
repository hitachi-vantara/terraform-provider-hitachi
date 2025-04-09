data "hitachi_vosb_block_volumes" "vosbvolumes" {
  vosb_block_address = ""
  #compute_node_name = "esxi-151" // Optional
  #compute_node_name = "MongoNode1" // Optional

}

output "volumeoutput" {
  value = data.hitachi_vosb_block_volumes.vosbvolumes
}
