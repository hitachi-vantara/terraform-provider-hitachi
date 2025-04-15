//
// Hitachi VOS Block Compute Port Resource
//
// This section defines a Terraform resource block to create a Hitachi VOS Block compute port.
// The resource "hitachi_vosb_compute_port" represents a compute port on a Hitachi VSP One SDS Block
// (VOSB) using its block interface and allows you to manage its configuration
// using Terraform.
//
// Customize the values of the parameters (vosb_address, name, authentication_settings,
// target_chap_users) to match your desired compute port configuration.
//
// The "target_chap_users" parameter specifies a list of CHAP users that are associated with
// this compute port for authentication purposes.
//

resource "hitachi_vosb_compute_port" "mycomputeport" {
  vosb_address = "10.10.12.13"
  name = "001-iSCSI-002"
  authentication_settings = "CHAP"
  target_chap_users = ["test_user7", "test_user9"]
}

output "chapuser_association_with_computeport" {
  value = resource.hitachi_vosb_compute_port.mycomputeport
}
