#
# Hitachi VSP iSCSI Target Data Retrieval
#
# This section defines a data source block to fetch information about a specific
# iSCSI target from a Hitachi Virtual Storage Platform (VSP) using HashiCorp
# Configuration Language (HCL).
#
#
# The data source block "hitachi_vsp_iscsi_targets" retrieves details about
# multiple iSCSI targets associated with the provided parameters. This allows
# you to access configuration and property information for specific targets.
#
# Customize the values of the parameters (serial, port_ids) to match your
# environment, enabling you to retrieve information about the desired iSCSI targets.
#

data "hitachi_vsp_iscsi_targets" "iscsitargets" {
  serial   = 12345
  port_ids = ["CL5-A", "CL4-C", "CL7-B"]
}

output "iscsitargets" {
  value = data.hitachi_vsp_iscsi_targets.iscsitargets
}
