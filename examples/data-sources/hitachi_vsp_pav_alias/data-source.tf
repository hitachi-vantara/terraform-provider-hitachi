#
# Hitachi VSP PAV Alias Data Retrieval
#
# This section defines a data source block to fetch PAV alias information from a
# Hitachi Virtual Storage Platform (VSP) using HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vsp_pav_alias" retrieves details about PAV aliases
# associated with the provided parameters. This enables you to access the list of
# PAV aliases returned by the storage system.
#
# Customize the values of the parameters (serial, cu_number, base_ldev_id, alias_ldev_ids)
# to match your environment. You can retrieve all PAV aliases, filter by a specific
# Control Unit (CU), filter by a base LDEV, or return only specific alias LDEVs.
#
# Parameters:
# - serial: (Required) Serial number of the storage system.
# - cu_number: (Optional) Control Unit number (0-255). If omitted, returns all PAV aliases.
# - base_ldev_id: (Optional) Filters alias entries whose base volume matches this LDEV.
# - alias_ldev_ids: (Optional) List of alias LDEV IDs.
#

data "hitachi_vsp_pav_alias" "all_pav_aliases" {
  serial = 12345
}

output "all_pav_aliases" {
  value = data.hitachi_vsp_pav_alias.all_pav_aliases.pav_aliases
}

# Filter by CU number
# data "hitachi_vsp_pav_alias" "cu2_pav_aliases" {
#   serial    = 12345
#   cu_number = 2
# }
#
# output "cu2_pav_aliases" {
#   value = data.hitachi_vsp_pav_alias.cu2_pav_aliases.pav_aliases
# }
