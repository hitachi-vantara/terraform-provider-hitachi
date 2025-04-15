#
# Hitachi VOS Block Storage Drives Data Retrieval
#
# This section defines a data source block to fetch information about storage drives
# from a Hitachi VSP One SDS Block (VOSB) using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vosb_storage_drives" retrieves details about storage drives
# associated with the provided parameters. This allows you to access configuration and property
# information for the specified storage drives.
#
# Customize the values of the parameters (vosb_address) to match your
# environment, enabling you to retrieve information about the desired storage drives.
#

data "hitachi_vosb_storage_drives" "my_drives" {
  vosb_address = var.vosb_address
}

output "my_drives_output" {
  value = data.hitachi_vosb_storage_drives.my_drives
}

