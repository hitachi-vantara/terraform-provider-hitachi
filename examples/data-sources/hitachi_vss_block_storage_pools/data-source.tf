#
# Hitachi VSS Block Storage Pools Data Retrieval
#
# This section defines a data source block to fetch information about specific storage pools
# from a Hitachi Virtual Storage System (VSS) using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vss_block_storage_pools" retrieves details about storage pools
# associated with the provided parameters. This allows you to access configuration and property
# information for the specified storage pools.
#
# Customize the values of the parameters (vssb_address, storage_pool_names) to match your
# environment, enabling you to retrieve information about the desired storage pools.
#

data "hitachi_vss_block_storage_pools" "pool" {
  vss_block_address  = var.vssb_address
  storage_pool_names = ["SP01"]
}

output "pool" {
  value = data.hitachi_vss_block_storage_pools.pool
}
