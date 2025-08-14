#
# Hitachi VSP Dynamic Pool Data Retrieval
#
# This section defines a data source block to fetch information about a specific
# dynamic pool from a Hitachi Virtual Storage Platform (VSP) using HashiCorp
# Configuration Language (HCL).
#
# The data source block "hitachi_vsp_dynamic_pool" retrieves details about a
# particular dynamic pool based on the provided parameters.
#
# Customize the values of the parameters (serial, either pool_id or pool_name) to align with your
# environment, allowing you to retrieve information about the desired dynamic pool.
#

data "hitachi_vsp_dynamic_pool" "dynamicpool" {
  serial  = 12345
  pool_id = 45
  # or
  # pool_name = "testPool"
}

output "dynamicpool" {
  value = data.hitachi_vsp_dynamic_pool.dynamicpool
}
