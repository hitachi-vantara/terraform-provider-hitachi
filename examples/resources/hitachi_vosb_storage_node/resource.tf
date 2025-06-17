// Hitachi VOS Block Storage Node Resource
//
// This section defines a Terraform resource block for managing storage node
// on a Hitachi VSP One SDS Block (VOSB) using using HashiCorp Configuration Language (HCL).
//
// The resource "hitachi_vosb_storage_node" represents the storage node on a Hitachi VSP One SDS Block
// (VOSB) using its block interface and allows you to manage its configuration
// using Terraform.
//
// Customize the values of the parameters (vosb_address, node_name, configuration_file, setup_user_password) 
// as needed to match your desired storage node configuration.
//
// - Set "node_name" to the name of the storage node to be added to the cluster.
// - Set "configuration_file" to be used to add the storage node.
// - Set "setup_user_password" to be used to add the storage node.
//
// Parameters:
// - node_name: name of the storage node to be added to the cluster.
// - configuration_file: configuration file to be used to add the storage node.
// - setup_user_password: password to use to log into the storage node to be added.


resource "hitachi_vosb_storage_node" "storageNode" {
  vosb_address = var.vosb_address
  node_name = "SDSB-NODE6"
  configuration_file = "/root/configuration.csv"
  setup_user_password = "password"
}

output "node_output" {
  value = resource.hitachi_vosb_storage_node.storageNode
}
