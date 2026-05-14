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
////////////////////////////////////////////////////////////////////////////////
// Example: Import (existing iSCSI CHAP user)
// -----------------------------------------------------------------------------
// Use terraform import when the CHAP user already exists on storage and
// you want to bring it under Terraform management without re-creating it.
// The import reads the current state from storage and writes it to tfstate.
//
// Import ID format: <vosb_address>/<target_chap_user_name>
//
// terraform import hitachi_vosb_iscsi_chap_user.imported 10.10.12.13/targetchapuser

// Minimal skeleton block required for `terraform import`.
// Create this block in your .tf before running the import command.
resource "hitachi_vosb_iscsi_chap_user" "imported" {}

output "imported_vosb_iscsi_chap_user_id" {
  value = hitachi_vosb_iscsi_chap_user.imported.id
}
////////////////////////////////////////////////////////////////////////////////

resource "hitachi_vosb_iscsi_chap_user" "my_chap_user" {
  vosb_address            = "10.10.12.13"
  target_chap_user_name   = "targetchapuser"
  target_chap_user_secret = var.target_chap_user_secret
}

output "chap_user_output" {
  # Explicitly specify 'chap_users' since it does not contain sensitive data.
  value = resource.hitachi_vosb_iscsi_chap_user.my_chap_user.chap_users

  # If you don't explicitly list output parameters, Terraform will display all inputs and outputs by default.
  # Since some input fields contain sensitive data, they must be marked as sensitive to avoid exposing them.
  # value     = resource.hitachi_vosb_iscsi_chap_user.my_chap_user
  # sensitive = true
}
