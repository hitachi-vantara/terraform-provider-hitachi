//
// Hitachi VOS Block Compute Node Resource
//
// This section defines a Terraform resource block to create a Hitachi VOS Block compute node.
// The resource "hitachi_vosb_compute_node" represents a compute node on a Hitachi VSP One SDS Block
// (VOSB) using its block interface and allows you to manage its configuration
// using Terraform.
//
// Customize the values of the parameters (vosb_address, compute_node_name, os_type),
// and the nested "fc_connection" blocks to match your desired compute node configuration.
//
// The "fc_connection" blocks define Fibre Channel connections for the compute node,
// including the host WWNs (World Wide Names).
//

resource "hitachi_vosb_compute_node" "mycompute" {
  vosb_address = "10.10.12.13"
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
  value = resource.hitachi_vosb_compute_node.mycompute
}

