#
# Hitachi VSP Hostgroup Data Retrieval
#
# This section defines a data source block and an output block to fetch
# and expose information about a specific hostgroup from a Hitachi Virtual
# Storage Platform (VSP) using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vsp_hostgroup" retrieves details about a
# particular hostgroup based on the provided parameters. The output block
# "hostgroup" exports the retrieved hostgroup information for further use.
#
# Adjust the values of the parameters (serial, port_id, hostgroup_number)
# according to your environment to fetch the desired hostgroup details.
#
data "hitachi_vsp_hostgroup" "hostgroup" {
  serial           = 12345
  port_id          = "CL1-A"
  hostgroup_number = 10
}

output "hostgroup" {
  value = data.hitachi_vsp_hostgroup.hostgroup
}
