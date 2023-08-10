data "hitachi_vss_block_dashboard" "dashboard" {
  vss_block_address = "10.10.12.13"
}

output "dashboardoutput" {
  value = data.hitachi_vss_block_dashboard.dashboard
}
