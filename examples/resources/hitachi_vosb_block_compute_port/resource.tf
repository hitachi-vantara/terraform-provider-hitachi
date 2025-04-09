//
// Hitachi VSS Block Compute Port Resource
//
// This section defines a Terraform resource block to create a Hitachi VSS block compute port.
// The resource "hitachi_vosb_block_compute_port" represents a compute port on a Hitachi Virtual
// Storage System (VSS) using its block interface and allows you to manage its configuration
// using Terraform.
//
// Customize the values of the parameters (vosb_block_address, name, authentication_settings,
// target_chap_users) to match your desired compute port configuration.
//
// The "target_chap_users" parameter specifies a list of CHAP users that are associated with
// this compute port for authentication purposes.
//

resource "hitachi_vosb_block_compute_port" "mycomputeport" {
  vosb_block_address = "10.10.12.13"
  name = "001-iSCSI-002"
  authentication_settings = "CHAP"
  target_chap_users = ["test_user7", "test_user9"]
}

output "chapuser_association_with_computeport" {
  value = resource.hitachi_vosb_block_compute_port.mycomputeport
}
