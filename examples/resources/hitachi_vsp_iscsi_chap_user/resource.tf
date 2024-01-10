//
// Hitachi VSP iSCSI CHAP User Resource
//
// This section defines a Terraform resource block to create a Hitachi VSP iSCSI CHAP user.
// The resource "hitachi_vsp_iscsi_chap_user" represents an iSCSI CHAP user on a Hitachi
// Virtual Storage Platform (VSP) and allows you to manage its configuration using Terraform.
//
// In the storage system, a CHAP user is uniquely identified by the combination of the
// serial, port_id, iscsi_target_number, chap_user_type, and chap_user_name. In Terraform,
// the resource is uniquely identified by the resource name "my_iscsi_initiator_chap_user3".
// The values of the mandatory input fields (serial, port_id, iscsi_target_number,
// chap_user_type, and chap_user_name) cannot be changed once the resource has been created.
//
// Customize the values of the parameters (serial, port_id, iscsi_target_number,
// chap_user_type, chap_user_name, chap_user_password) to match your desired CHAP user configuration.
//

resource "hitachi_vsp_iscsi_chap_user" "my_iscsi_initiator_chap_user3" {
  serial              = 12345                
  port_id             = "CL4-C"              
  iscsi_target_number = 01                   
  chap_user_type      = "initiator"          
  chap_user_name      = "chapuser"           
  chap_user_password  = "TopSecretForMyChap" 
}
