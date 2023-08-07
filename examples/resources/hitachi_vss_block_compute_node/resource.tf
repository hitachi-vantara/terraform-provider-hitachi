resource "hitachi_vss_block_compute_node" "mycompute" {
  vss_block_address = ""
  compute_node_name = "ComputeNode-RESTAPI123"
  os_type = "VMware"
  

  fc_connection {
    host_wwn = "60060e8107595326"

  }

  fc_connection {
    host_wwn = "90060e8107595325"
  }

}

output "computenodecreate" {
  value = resource.hitachi_vss_block_compute_node.mycompute
}
