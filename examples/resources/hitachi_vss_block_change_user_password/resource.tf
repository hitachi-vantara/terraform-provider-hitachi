// Hitachi VSS Block Storage Credentials Resource
//
// This section defines a Terraform resource block for changing the password of a registered storage user 
// on a Hitachi Virtual Storage System (VSS) using using HashiCorp Configuration Language (HCL).
//
// The resource "hitachi_vss_block_change_user_password" allows you to change the password of a registered
// storage user on the VSS.
//
// Customize the values of the parameters (vss_block_address, user_id, current_password, and new_password) 
// as needed to match your environment, reflecting the appropriate settings for the registered storage user
// whose password you wish to change.

resource "hitachi_vss_block_change_user_password" "my_user" {
  vss_block_address = var.vssb_address
  user_id           = "testUser"
  current_password  = var.current_password
  new_password      = var.new_password
}

output "user_output" {
  value = resource.hitachi_vss_block_change_user_password.my_user
}
