#
# Hitachi VSS Block iSCSI CHAP Users Data Retrieval
#
# This section defines multiple data source blocks to fetch iSCSI CHAP user information
# from a Hitachi Virtual Storage System (VSS) using HashiCorp Configuration Language (HCL).
#
# Each data source block in this configuration retrieves details about iSCSI CHAP users
# associated with the provided parameters. This enables you to access authentication
# information for specific iSCSI CHAP users.
#
# Customize the values of the parameters (vss_block_address, target_chap_user) as needed
# to match your environment, allowing you to retrieve information about the desired iSCSI CHAP users.
#

# Retrieve iSCSI CHAP user by ID
data "hitachi_vss_block_iscsi_chap_users" "chap_user_by_id" {
  vss_block_address   = var.vssb_address
  target_chap_user = "a79c1a1d-2719-4e07-b800-faf9de73d0ae" //chap user id
}

output "id_output" {
  value = data.hitachi_vss_block_iscsi_chap_users.chap_user_by_id
}

# Retrieve iSCSI CHAP user by name
data "hitachi_vss_block_iscsi_chap_users" "chap_user_by_name" {
  vss_block_address = var.vssb_address
  target_chap_user = "chapusername"
}

output "name_output" {
  value = data.hitachi_vss_block_iscsi_chap_users.chap_user_by_name
}

# Retrieve all iSCSI CHAP users
data "hitachi_vss_block_iscsi_chap_users" "my_chap_users" {
  vss_block_address = var.vssb_address
}

output "my_iscsi_chap_users_output" {
  value = data.hitachi_vss_block_iscsi_chap_users.my_chap_users
}
