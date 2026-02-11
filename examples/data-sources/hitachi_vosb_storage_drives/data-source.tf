#
# Hitachi VSP One SDS Block Storage Drives Data Retrieval
#
# This section defines a data source block to fetch information about storage drives
# from a Hitachi VSP One SDS Block using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vosb_storage_drives" retrieves details about storage drives
# associated with the provided parameters. This allows you to access configuration and property
# information for the specified storage drives.
#
# Customize the values of the parameters (e.g., vosb_address) to match your
# environment, enabling you to retrieve information about the desired storage drives.
#
# Optional Input:
# - status: Filters drives by their status. Allowed values (case-insensitive): 
#   "", "Offline", "Normal", "TemporaryBlockage", "Blockage"

data "hitachi_vosb_storage_drives" "my_drives" {
  vosb_address = "10.10.12.13"
  status       = "normal"
}

output "my_drives_output" {
  value = data.hitachi_vosb_storage_drives.my_drives
}

