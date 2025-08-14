#
# Hitachi VSP One SDS Block Compute Nodes Data Retrieval
#
# This section defines a data source block to fetch information about a specific
# compute node from a Hitachi VSP One SDS Block using HashiCorp
# Configuration Language (HCL).
#
# The data source block "hitachi_vosb_compute_nodes" retrieves details
# about a compute node associated with the provided parameters. This allows
# you to access configuration and property information for the specified compute node.
#
# Customize the values of the parameters (vosb_address, compute_node_name)
# to match your environment, enabling you to retrieve information about the desired
# compute node.
#

data "hitachi_vosb_compute_nodes" "computenodes" {
  vosb_address      = "10.10.12.13"
  compute_node_name = "ComputeNode-RESTAPI123"
}

output "nodeoutput" {
  value = data.hitachi_vosb_compute_nodes.computenodes
}
