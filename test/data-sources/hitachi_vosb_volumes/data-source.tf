data "hitachi_vosb_volumes" "vosbvolumes" {
  vosb_address = ""
  #compute_node_name = "esxi-151" // Optional
  #compute_node_name = "MongoNode1" // Optional

}

output "volumeoutput" {
  value = data.hitachi_vosb_volumes.vosbvolumes
}
