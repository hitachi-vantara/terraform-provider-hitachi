#
# Hitachi VSP Parity Groups Data Retrieval
#
# This section defines a data source block to fetch information about specific
# parity groups from a Hitachi Virtual Storage Platform (VSP) using HashiCorp
# Configuration Language (HCL).
#
# The data source block "hitachi_vsp_parity_groups" retrieves details about
# parity groups associated with the provided parameters. This enables you to
# access configuration and property information for the specified parity groups.
#
# Customize the values of the parameters (serial, parity_group_ids) to match
# your environment, allowing you to retrieve information about the desired
# parity groups.
#
/*
data "hitachi_vsp_parity_groups" "myparitygroup" {
  serial = 12345
  parity_group_ids = ["1-2","1-3"]
}

output "myparitygroup" {
  value = data.hitachi_vsp_parity_groups.myparitygroup
}

data "hitachi_infra_storage_devices" "storage_devices" {
  #storage_id = "storage-9a3f87a8c9dc213e8ebd02b63b97b9e8"
  serial = 611039
}
*/
data "hitachi_vsp_parity_groups" "parity_groups" {
  #storage_id = data.hitachi_infra_storage_devices.storage_devices.id
  serial = 611039
  parity_group_ids = [ "E1-1", "E1-2"]
}

output "parity_groups" {
  value = data.hitachi_vsp_parity_groups.parity_groups
}

