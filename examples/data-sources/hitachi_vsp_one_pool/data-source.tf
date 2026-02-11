#
# Hitachi VSP One Pool Information Retrieval
#
# This section defines a data source block that retrieves detailed information
# about a specific storage pool managed by Hitachi storage systems through the
# Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_pool" data source allows you to access configuration
# and property details of an existing storage pool, enabling you to reference pool
# information in Terraform configurations and automate resource dependencies.
#
# This data source is useful for retrieving pool details such as capacity,
# encryption status, drive configuration, and performance metrics for further
# automation, monitoring, or validation purposes.
#


# Get information about a specific pool
data "hitachi_vsp_one_pool" "target_pool" {
  serial  = 12345
  pool_id = 0
}

# Output detailed pool information
output "pool_details" {
  value = data.hitachi_vsp_one_pool.target_pool.data
}