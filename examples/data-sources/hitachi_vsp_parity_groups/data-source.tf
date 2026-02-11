#
# Hitachi VSP Parity Groups Data Retrieval
#
# This section defines a data source block to fetch information about specific
# parity groups from a Hitachi Virtual Storage Platform (VSP) using HashiCorp
# Configuration Language (HCL).
#
# The data source block "hitachi_vsp_parity_groups" retrieves details about
# parity groups associated with the provided parameters. This enables you to
# access configuration and property information for the specified parity groups.
#
# Customize the values of the parameters (serial, parity_group_ids, clpr_id, include_detail_info
# include_cache_info, drive_type_name) to match your environment, allowing you to retrieve 
# information about the desired parity groups.
#

data "hitachi_vsp_parity_groups" "myparitygroups" {
  serial           = 12345
  parity_group_ids = ["1-2", "1-3"]
}

output "myparitygroups" {
  value = data.hitachi_vsp_parity_groups.myparitygroups
}


# Example: Parity groups with optional detail flags
data "hitachi_vsp_parity_groups" "myparitygroups_with_details" {
  serial = 12345

  # - include_detail_info: (Optional, default: false)
  #   When true, the provider requests additional detailed information (FMC) for
  #   parity groups.
  include_detail_info = true

  # - include_cache_info: (Optional, default: false)
  #   When true, the provider requests cache/class related information (class) for
  #   parity groups.
  include_cache_info = true

  # Notes:
  # - You can enable either flag independently or both together.
}

output "myparitygroups_with_details" {
  value = data.hitachi_vsp_parity_groups.myparitygroups_with_details
}


# Example: Parity groups filtered by drive type + CLPR
data "hitachi_vsp_parity_groups" "myparitygroups_with_drive_type_and_clpr" {
  serial = 12345

  # Optional filters (server-side):
  # - drive_type_name (string)
  #   Filters parity groups by the drive type name reported by the storage.
  #   Valid values:
  #   - VSP One B20: SSD(QLC), SSD
  #   - VSP 5000 series: SAS, SSD(MLC), SSD(FMC), SSD, SCM
  #   - VSP E series: SAS, SSD(MLC), SSD
  #   - VSP G350/G370/G700/G900, VSP F350/F370/F700/F900: SAS, SSD(MLC), SSD(FMC), SSD(RI)
  drive_type_name = "SSD"

  # - clpr_id (int)
  #   Filters parity groups by CLPR number.
  clpr_id = 0

  # Optional detail flags (same pattern as other SAN data sources):
  # - include_detail_info=true adds detailInfoType=FMC
  # - include_cache_info=true  adds detailInfoType=class
  include_detail_info = true
  include_cache_info  = true
}

output "myparitygroups_with_drive_type_and_clpr" {
  value = data.hitachi_vsp_parity_groups.myparitygroups_with_drive_type_and_clpr
}
