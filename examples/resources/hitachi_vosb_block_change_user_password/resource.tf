// Hitachi VOS Block Storage Credentials Resource
//
// This section defines a Terraform resource block for changing the password of a registered storage user 
// on a Hitachi VSP One SDS Block (VOSB) using using HashiCorp Configuration Language (HCL).
//
// The resource "hitachi_vosb_block_change_user_password" allows you to change the password of a registered
// storage user on the VOSB.
//
// Customize the values of the parameters (vosb_block_address, user_id, current_password, and new_password) 
// as needed to match your environment, reflecting the appropriate settings for the registered storage user
// whose password you wish to change.

resource "hitachi_vosb_block_change_user_password" "my_user" {
  vosb_block_address = "10.10.12.13"
  user_id            = "testUser"
  current_password   = var.current_password
  new_password       = var.new_password
}

output "user_output" {
  value = resource.hitachi_vosb_block_change_user_password.my_user
}
