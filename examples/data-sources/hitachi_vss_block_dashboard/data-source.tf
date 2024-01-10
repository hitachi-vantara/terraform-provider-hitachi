#
# Hitachi VSS Block Dashboard Data Retrieval
#
# This section defines a data source block to fetch dashboard information from a
# Hitachi Virtual Storage System (VSS) using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vss_block_dashboard" retrieves a dashboard of
# information associated with the provided parameters. This allows you to access
# an overview of various metrics and details about the storage system.
#
# Customize the value of the parameter (vss_block_address) to match your environment,
# enabling you to retrieve dashboard information from the desired VSS.
#

data "hitachi_vss_block_dashboard" "dashboard" {
  vss_block_address = "10.10.12.13"
}

output "dashboardoutput" {
  value = data.hitachi_vss_block_dashboard.dashboard
}
