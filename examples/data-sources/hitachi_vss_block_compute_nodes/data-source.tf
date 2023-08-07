data "hitachi_vss_block_compute_nodes" "computenodes" {
  vss_block_address = ""
  compute_node_name = "ComputeNode-RESTAPI123" // Optional
}

output "nodeoutput" {
  value = data.hitachi_vss_block_compute_nodes.computenodes
}
