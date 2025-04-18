#
# Hitachi VSP Storage System Data Retrieval
#
# This section defines a data source block to fetch information about a specific
# storage system from a Hitachi Virtual Storage Platform (VSP) using HashiCorp
# Configuration Language (HCL).
#
# The data source block "hitachi_vsp_storage" retrieves details about a storage
# system associated with the provided parameters. This allows you to access
# configuration and property information for the specified storage system.
#
# Customize the value of the parameter (serial) to match your environment,
# enabling you to retrieve information about the desired storage system.
#

data "hitachi_vsp_storage" "s12345" {
  serial = 12345
}

output "s12345" {
  value = data.hitachi_vsp_storage.s12345
}




