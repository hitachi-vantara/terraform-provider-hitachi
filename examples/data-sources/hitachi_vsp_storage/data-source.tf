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

# # Additional details for storage system
# data "hitachi_vsp_storage" "s12345_with_details" {
#   serial = 12345

#   # - include_detail_info (Optional, default: false)
#   #   When true, the provider requests additional “detailInfoType” data for the
#   #   storage system.
#   #   Populates extra output fields that are not included in the baseline response
#   #   (for example, compression acceleration availability and detailed version
#   #   information as supported by the storage firmware/API).
#   include_detail_info = true
# }

# output "s12345_with_details" {
#   value = data.hitachi_vsp_storage.s12345_with_details
# }




