#
# Hitachi VSS Block Storage Drives Data Retrieval
#
# This section defines a data source block to fetch information about storage drives
# from a Hitachi Virtual Storage System (VSS) using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vosb_block_storage_drives" retrieves details about storage drives
# associated with the provided parameters. This allows you to access configuration and property
# information for the specified storage drives.
#
# Customize the values of the parameters (vosb_block_address) to match your
# environment, enabling you to retrieve information about the desired storage drives.
#

data "hitachi_vosb_block_storage_drives" "my_drives" {
  vosb_block_address = var.vosb_block_address
}

output "my_drives_output" {
  value = data.hitachi_vosb_block_storage_drives.my_drives
}

