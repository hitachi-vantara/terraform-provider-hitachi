data "hitachi_vosb_block_dashboard" "dashboard" {
  vosb_block_address = "10.76.47.55"
}

output "dashboardoutput" {
  value = data.hitachi_vosb_block_dashboard.dashboard
}
