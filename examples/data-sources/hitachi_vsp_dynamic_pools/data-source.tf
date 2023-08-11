data "hitachi_vsp_dynamic_pools" "dynamicpool" {
  serial  = 12345
  pool_id = 45
}

output "dynamicpool" {
  value = data.hitachi_vsp_dynamic_pools.dynamicpool
}
