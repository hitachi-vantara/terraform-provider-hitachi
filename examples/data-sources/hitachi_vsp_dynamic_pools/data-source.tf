#
# Hitachi VSP Dynamic Pool Data Retrieval
#
# This section defines a data source block to fetch information about a specific
# dynamic pool from a Hitachi VSP One SDS Block Storage Platform (VSP) using HashiCorp
# Configuration Language (HCL).
#
# The data source block "hitachi_vsp_dynamic_pools" retrieves details about a
# particular dynamic pool based on the provided parameters.
#
# Customize the values of the parameters (serial, pool_id) to align with your
# environment, allowing you to retrieve information about the desired dynamic pool.
#
data "hitachi_vsp_dynamic_pools" "dynamicpool" {
  serial  = 12345
  pool_id = 45
}

output "dynamicpool" {
  value = data.hitachi_vsp_dynamic_pools.dynamicpool
}
