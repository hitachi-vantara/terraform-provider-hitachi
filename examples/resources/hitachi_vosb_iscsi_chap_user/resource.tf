//
// Hitachi VSP One SDS Block iSCSI CHAP User Resource
//
// This section defines a Terraform resource block to create a Hitachi VSP One SDS Block iSCSI CHAP user.
// The resource "hitachi_vosb_iscsi_chap_user" represents an iSCSI CHAP user on a Hitachi VSP One SDS Block
// using its block interface and allows you to manage its configuration using Terraform.
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

output "chap_user_output" {
  # Explicitly specify 'chap_users' since it does not contain sensitive data.
  value = resource.hitachi_vosb_iscsi_chap_user.my_chap_user.chap_users

  # If you don't explicitly list output parameters, Terraform will display all inputs and outputs by default.
  # Since some input fields contain sensitive data, they must be marked as sensitive to avoid exposing them.
  # value     = resource.hitachi_vosb_iscsi_chap_user.my_chap_user
  # sensitive = true
}
