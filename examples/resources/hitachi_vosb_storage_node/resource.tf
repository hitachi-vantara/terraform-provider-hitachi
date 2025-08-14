// Hitachi VSP One SDS Block Storage Node Resource
//
// This section defines a Terraform resource block for managing storage node
// on a Hitachi VSP One SDS Block using HashiCorp Configuration Language (HCL).
//
// The resource "hitachi_vosb_storage_node" represents the storage node on a Hitachi VSP One SDS Block
// using its block interface and allows you to manage its configuration
// using Terraform.
//

// Expected Cloud Provider Behavior
// The `expected_cloud_provider` (String) parameter specifies the expected cloud provider type. Valid values: "google", "azure", "baremetal".
//	- Used to validate combinations of inputs based on the deployment environment.
//	- If set to "google" or "azure", specific parameters may be required for certain operations.
//	- If set to "baremetal" (default), other cloud-specific inputs are ignored.
//	- Note: The actual cloud provider is determined by the VSP One SDS Block system at the "vosb_address" endpoint.
//	If there's a mismatch, the request still proceeds and behaves according to the actual environment.

/////////////////////////////// Azure /////////////////////////////////
// Customize the values of the parameters (vosb_address, exported_configuration_file, expected_cloud_provider) 
// as needed to match your desired Azure storage node configuration.
//
// - Set "exported_configuration_file" to be used to add the storage node.
// - Set "expected_cloud_provider" to be used to add the storage node.
//
// Parameters:
// - exported_configuration_file: configuration file to be used to add the storage node.
// - expected_cloud_provider: cloud platform.
//
// Example:
// resource "hitachi_vosb_storage_node" "storageNode" {
//   vosb_address = var.vosb_address
//   exported_configuration_file = "/tmp/configuration.tar.gz"
//   expected_cloud_provider = "azure"
// }

/////////////////////////// Google Cloud Platfor (GCP) /////////////////////////////
// Customize the values of the parameters (vosb_address, exported_configuration_file, expected_cloud_provider) 
// as needed to match your desired GCP storage node configuration.
//
// - Set "exported_configuration_file" to be used to add the storage node.
// - Set "expected_cloud_provider" to be used to add the storage node.
//
// Parameters:
// - exported_configuration_file: configuration file to be used to add the storage node.
// - expected_cloud_provider: cloud platform.
//
// Example:
// resource "hitachi_vosb_storage_node" "storageNode" {
//   vosb_address = var.vosb_address
//   exported_configuration_file = "/tmp/configuration.csv"
//   expected_cloud_provider = "google"
// }

//////////////////////////////// Baremetal /////////////////////////////////
// Customize the values of the parameters (vosb_address, configuration_file, setup_user_password) 
// as needed to match your desired baremetal storage node configuration.
//
// - Set "configuration_file" to be used to add the storage node.
// - Set "setup_user_password" to be used to add the storage node.
//
// Parameters:
// - configuration_file: configuration file to be used to add the storage node.
// - setup_user_password: password to use to log into the storage node to be added.

resource "hitachi_vosb_storage_node" "storageNode" {
  vosb_address = var.vosb_address
  configuration_file = "/tmp/configuration.csv"
  setup_user_password = "password"
}

output "node_output" {
  value = resource.hitachi_vosb_storage_node.storageNode
}
