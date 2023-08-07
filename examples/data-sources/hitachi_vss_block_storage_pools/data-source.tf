data "hitachi_vss_block_storage_pools" "pool" {
  vssb_address = ""
  // storage_pool_names = ["SP01"]
}

output "pool" {
  value = data.hitachi_vss_block_storage_pools.pool
}
