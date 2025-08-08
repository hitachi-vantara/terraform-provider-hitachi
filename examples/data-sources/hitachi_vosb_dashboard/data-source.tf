#
# Hitachi VSP One SDS Block Dashboard Data Retrieval
#
# This section defines a data source block to fetch dashboard information from a
# Hitachi VSP One SDS Block using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vosb_dashboard" retrieves a dashboard of
# information associated with the provided parameters. This allows you to access
# an overview of various metrics and details about the storage system.
#
# Customize the value of the parameter (vosb_address) to match your environment,
# enabling you to retrieve dashboard information.
#

data "hitachi_vosb_dashboard" "dashboard" {
  vosb_address = "10.10.12.13"
}

output "dashboardoutput" {
  value = data.hitachi_vosb_dashboard.dashboard
}
