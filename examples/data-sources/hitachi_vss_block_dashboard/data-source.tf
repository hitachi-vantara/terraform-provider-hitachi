data "hitachi_vss_block_dashboard" "dashboard" {
  vss_block_address = ""
}

output "dashboardoutput" {
  value = data.hitachi_vss_block_dashboard.dashboard
}
