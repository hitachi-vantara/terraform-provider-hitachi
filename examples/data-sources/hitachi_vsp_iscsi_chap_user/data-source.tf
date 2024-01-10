#
# Hitachi VSP iSCSI CHAP User Data Retrieval
#
# This section defines a data source block to fetch information about a specific
# iSCSI CHAP user from a Hitachi Virtual Storage Platform (VSP) using HashiCorp
# Configuration Language (HCL).
#
# The data source block "hitachi_vsp_iscsi_chap_user" retrieves details about an
# iSCSI CHAP user associated with the provided parameters. This allows you to
# access authentication information for a specific initiator.
#
# Customize the values of the parameters (serial, port_id, iscsi_target_number,
# chap_user_type, chap_user_name) to match your environment, thereby enabling
# the retrieval of information about the desired iSCSI CHAP user.
#

data "hitachi_vsp_iscsi_chap_user" "my_iscsi_initiator_chap_user" {
  serial              = 12345
  port_id             = "CL4-C"
  iscsi_target_number = 1
  chap_user_type      = "initiator" 
  chap_user_name      = "chapuser1"

}

output "my_iscsi_initiator_chap_user_output" {
  value = data.hitachi_vsp_iscsi_chap_user.my_iscsi_initiator_chap_user
}

#
# The data source block "hitachi_vsp_iscsi_chap_users" retrieves details about
# iSCSI CHAP users associated with the provided parameters. This allows you to
# access authentication information for specific initiators on a given target.
#
# Customize the values of the parameters (serial, port_id, iscsi_target_number)
# to match your environment, thereby enabling the retrieval of information about
# the desired iSCSI CHAP users.
#
data "hitachi_vsp_iscsi_chap_users" "my_iscsi_chap_users" {
  serial              = 12345
  port_id             = "CL4-C"
  iscsi_target_number = 1

}

output "my_iscsi_chap_users_output" {
  value = data.hitachi_vsp_iscsi_chap_users.my_iscsi_chap_users
}