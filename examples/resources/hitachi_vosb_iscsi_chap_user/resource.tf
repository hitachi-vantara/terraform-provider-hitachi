//
// Hitachi VOS Block iSCSI CHAP User Resource
//
// This section defines a Terraform resource block to create a Hitachi VOS Block iSCSI CHAP user.
// The resource "hitachi_vosb_iscsi_chap_user" represents an iSCSI CHAP user on a Hitachi
// Virtual (VOSB) using its block interface and allows you to manage its configuration
// using Terraform.
//
// Customize the values of the parameters (vosb_address, target_chap_user_name,
// target_chap_user_secret) to match your desired iSCSI CHAP user configuration.
//

resource "hitachi_vosb_iscsi_chap_user" "my_chap_user" {
  vosb_address            = "10.10.12.13"
  target_chap_user_name   = "targetchapuser"
  target_chap_user_secret = "targetchapuserpasswd"
}

output "chap_user_output" {
  value = resource.hitachi_vosb_iscsi_chap_user.my_chap_user
}

