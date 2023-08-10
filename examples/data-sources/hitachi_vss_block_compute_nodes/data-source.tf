data "hitachi_vss_block_compute_nodes" "computenodes" {
  vss_block_address = "10.10.12.13"
  compute_node_name = "ComputeNode-RESTAPI123"
}

output "nodeoutput" {
  value = data.hitachi_vss_block_compute_nodes.computenodes
}
