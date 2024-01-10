#
# Hitachi VSS Block iSCSI Port Authentication Data Retrieval
#
# This section defines a data source block to fetch authentication information for
# a specific iSCSI port from a Hitachi Virtual Storage System (VSS) using HashiCorp
# Configuration Language (HCL).
#
# The data source block "hitachi_vss_block_iscsi_port_auth" retrieves authentication
# details for an iSCSI port associated with the provided parameters. This allows you
# to access security-related information for the specified iSCSI port.
#
# Customize the values of the parameters (vss_block_address, name) to match your
# environment, enabling you to retrieve authentication information for the desired iSCSI port.
#

data "hitachi_vss_block_iscsi_port_auth" "mycomputeport" {
  vss_block_address = "10.10.12.13"
  name = "001-iSCSI-002"
}

output "mycomputeport" {
  value = data.hitachi_vss_block_iscsi_port_auth.mycomputeport
}
