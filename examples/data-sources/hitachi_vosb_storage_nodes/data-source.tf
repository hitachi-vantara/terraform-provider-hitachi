#
# Hitachi VSP One SDS Block Storage Nodes Data Retrieval
#
# This section defines a data source block to fetch information about a storage nodes
# from a Hitachi VSP One SDS Block using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vosb_storage_nodes" retrieves details about storage nodes
# associated with the provided parameters. This allows you to access configuration and property
# information for the storage nodes.
#
# Customize the values of the parameters (vosb_address) to match your environment,
# enabling you to retrieve information about the storage nodes.
#

data "hitachi_vosb_storage_nodes" "storageNodes" {
  vosb_address = "10.10.12.13"
}

output "storageNodes" {
  value = data.hitachi_vosb_storage_nodes.storageNodes
}
