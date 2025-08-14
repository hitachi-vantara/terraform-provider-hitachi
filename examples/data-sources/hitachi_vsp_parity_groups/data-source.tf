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

data "hitachi_vsp_parity_groups" "myparitygroups" {
  serial = 12345
  parity_group_ids = ["1-2","1-3"]
}

output "myparitygroups" {
  value = data.hitachi_vsp_parity_groups.myparitygroups
}
