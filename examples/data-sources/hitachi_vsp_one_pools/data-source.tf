#
# Hitachi VSP One Pools Information Retrieval
#
# This section defines a data source block that retrieves information about
# multiple storage pools managed by Hitachi storage systems through the
# Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_pools" data source allows you to query and filter
# storage pools based on various criteria such as name patterns, status,
# or configuration state. This enables bulk operations, monitoring, and
# automated resource discovery across multiple pools.
#
# This data source is useful for:
# - Discovering all available pools in a storage system
# - Filtering pools based on naming conventions or status
# - Monitoring pool capacity and utilization across the environment
# - Automating pool selection for volume provisioning
#


# Get all pools on the storage system
data "hitachi_vsp_one_pools" "all" {
  serial = 12345
}

# Output all pool information
output "all_pools" {
  value = data.hitachi_vsp_one_pools.all.data
}

# Get pools with a specific name pattern
data "hitachi_vsp_one_pools" "production_pools" {
  serial      = 12345
  name_filter = "production"
}

output "production_pool_info" {
  value = data.hitachi_vsp_one_pools.production_pools.data
}

# Get pools with normal status
data "hitachi_vsp_one_pools" "healthy_pools" {
  serial        = 12345
  status_filter = "Normal"
}

# Output filtered results
output "healthy_pools_info" {
  value = data.hitachi_vsp_one_pools.healthy_pools.data
}

