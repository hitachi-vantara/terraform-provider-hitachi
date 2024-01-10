#
# Hitachi VSS Block Storage Ports Data Retrieval
#
# This section defines a data source block to fetch information about a specific storage port
# from a Hitachi Virtual Storage System (VSS) using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vss_block_storage_ports" retrieves details about a storage port
# associated with the provided parameters. This allows you to access configuration and property
# information for the specified storage port.
#
# Customize the values of the parameters (vss_block_address, port_name) to match your environment,
# enabling you to retrieve information about the desired storage port.
#

data "hitachi_vss_block_storage_ports" "storagePorts" {
  vss_block_address = "10.10.12.13"
  port_name = "001-iSCSI-002"
}

output "storagePorts" {
  value = data.hitachi_vss_block_storage_ports.storagePorts
}
