#
# Hitachi VSP One Storage Information Retrieval
#
# This section defines a data source block that retrieves detailed information
# about a specific Hitachi VSP storage system using the "hitachi_vsp_one_storage" data source.
#
# The data source allows you to access configuration and property details of an existing
# storage device, enabling you to reference its attributes in Terraform configurations
# and automate resource dependencies.
#
# Adjust the parameters (for example, serial) to match your environment and retrieve
# information for the desired storage system.
#

data "hitachi_vsp_one_storage" "s12345" {
  serial = 12345
  with_estimated_configurable_capacities = true
}

output "s12345" {
  value = data.hitachi_vsp_one_storage.s12345
}




