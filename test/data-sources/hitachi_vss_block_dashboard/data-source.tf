data "hitachi_vss_block_dashboard" "dashboard" {
  vss_block_address = "10.76.47.55"
}

output "dashboardoutput" {
  value = data.hitachi_vss_block_dashboard.dashboard
}
