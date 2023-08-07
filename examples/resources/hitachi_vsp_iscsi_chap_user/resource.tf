// In Storage System, chap user is uniquely identified by serial, port_id, iscsi_target_number, chap_user_type and chap_user_name
// In Terraform, chap user is uniquely identified by resource name "my_iscsi_initiator_chap_user"
// The value of the mandatory input fields (serial, port_id, iscsi_target_number, chap_user_type and chap_user_name) 
// cann't be changed one the resource has been created 

resource "hitachi_vsp_iscsi_chap_user" "my_iscsi_initiator_chap_user3" {
  serial              = 30078                # mandatory input
  port_id             = "CL4-C"              # mandatory input
  iscsi_target_number = 01                   # mandatory input
  chap_user_type      = "initiator"          # mandatory input. valid input value : "initiator", "target"
  chap_user_name      = "chapuser"             # mandatory input
  chap_user_password  = "TopSecretForMyChap" # optional input
 

}

