data "hitachi_vosb_dashboard" "dashboard" {
  vosb_address = "10.76.47.55"
}

output "dashboardoutput" {
  value = data.hitachi_vosb_dashboard.dashboard
}
