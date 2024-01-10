data "hitachi_vss_block_storage_pools" "pool" {
  vssb_address = "10.76.47.55"
  // storage_pool_names = ["SP01"]
}

output "pool" {
  value = data.hitachi_vss_block_storage_pools.pool
}
