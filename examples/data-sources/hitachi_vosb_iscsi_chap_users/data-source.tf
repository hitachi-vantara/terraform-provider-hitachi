#
# Hitachi VOS Block iSCSI CHAP Users Data Retrieval
#
# This section defines multiple data source blocks to fetch iSCSI CHAP user information
# from a Hitachi VSP One SDS Block (VOSB) using HashiCorp Configuration Language (HCL).
#
# Each data source block in this configuration retrieves details about iSCSI CHAP users
# associated with the provided parameters. This enables you to access authentication
# information for specific iSCSI CHAP users.
#
# Customize the values of the parameters (vosb_address, target_chap_user) as needed
# to match your environment, allowing you to retrieve information about the desired iSCSI CHAP users.
#

# Retrieve iSCSI CHAP user by ID
data "hitachi_vosb_iscsi_chap_users" "chap_user_by_id" {
  vosb_address = "10.10.12.13"
  target_chap_user = "a79c1a1d-2719-4e07-b800-faf9de73d0ae" //chap user id
}

output "id_output" {
  value = data.hitachi_vosb_iscsi_chap_users.chap_user_by_id
}

# Retrieve iSCSI CHAP user by name
data "hitachi_vosb_iscsi_chap_users" "chap_user_by_name" {
  vosb_address = "10.10.12.13"
  target_chap_user = "chapusername"
}

output "name_output" {
  value = data.hitachi_vosb_iscsi_chap_users.chap_user_by_name
}

# Retrieve all iSCSI CHAP users
data "hitachi_vosb_iscsi_chap_users" "my_chap_users" {
  vosb_address = "10.10.12.13"
}

output "my_iscsi_chap_users_output" {
  value = data.hitachi_vosb_iscsi_chap_users.my_chap_users
}
