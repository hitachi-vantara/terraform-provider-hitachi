#
# Hitachi VSP One Volumes Retrieval (Filtered List)
#
# This section defines a data source block that retrieves a list of existing
# volumes from the Hitachi Virtual Storage Platform One Block Administrator.
#
# The "hitachi_vsp_one_volumes" data source supports flexible filtering options,
# allowing you to narrow results by pool, server, nickname, or capacity range.
# This is useful for auditing, reporting, or referencing existing volumes in
# Terraform configurations.
#
# You can specify filtering parameters such as:
#   - pool_id or pool_name (partial match) to target specific pools
#   - server_id or server_nickname (partial match) to find volumes connected to servers
#   - nickname for partial name matching
#   - capacity ranges using min/max_total_capacity_mb or min/max_used_capacity_mb
#   - pagination controls via start_volume_id and requested_count
#
# Adjust the parameters to fit your environment and query requirements.
#

data "hitachi_vsp_one_volumes" "volumes" {
  serial    = var.serial_number
  pool_id   = 0
  pool_name = ""
  # server_id        = 100
  server_nickname = ""
  # nickname         = "test_volume"
  # max_total_capacity_mb = 102400
  # min_total_capacity_mb = 5120
  # max_used_capacity_mb  = 2048
  # min_used_capacity_mb  = 0
  start_volume_id = 6670
  requested_count = 20
}

output "volumes_info" {
  value = data.hitachi_vsp_one_volumes.volumes.volumes_info
}

output "volume_count" {
  value = data.hitachi_vsp_one_volumes.volumes.volume_count
}

output "total_count" {
  value = data.hitachi_vsp_one_volumes.volumes.total_count
}
