#
# Hitachi VSP Dynamic Pools Data Retrieval
#
# This data source block fetches information about all dynamic pools
# from a Hitachi Virtual Storage Platform (VSP) using
# HashiCorp Configuration Language (HCL).
#
# The "hitachi_vsp_dynamic_pools" data source retrieves details such as
# pool status, usage rate, and threshold settings.
#
# Customize the values of the parameters (serial) to align with your
# environment, allowing you to retrieve information about all dynamic pools.
#

data "hitachi_vsp_dynamic_pools" "dynamicpools" {
  serial  = 12345
}

output "dynamicpools" {
  value = data.hitachi_vsp_dynamic_pools.dynamicpools
}
