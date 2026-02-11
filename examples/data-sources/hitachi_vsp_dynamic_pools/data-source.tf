#
# Hitachi VSP Dynamic Pools Data Retrieval
#
# This data source block fetches information about all dynamic pools
# from a Hitachi Virtual Storage Platform (VSP) using
# HashiCorp Configuration Language (HCL).
#
# The "hitachi_vsp_dynamic_pools" data source retrieves details such as
# pool status, usage rate, and threshold settings.
#
# Customize the values of the parameters (serial) to align with your
# environment, allowing you to retrieve information about all dynamic pools.
#

data "hitachi_vsp_dynamic_pools" "dynamicpools" {
  serial = 12345
}

output "dynamicpools" {
  value = data.hitachi_vsp_dynamic_pools.dynamicpools
}


# data "hitachi_vsp_dynamic_pools" "dynamicpools_with_filters" {
#   serial = 12345

#   # - pool_type: (Optional) Filters returned pools by type.
#   #   Supported values:
#   #     - "DP"  : Dynamic Provisioning pool
#   #     - "HTI" : Hitachi Thin Image pool
#   #   If omitted, the provider returns all pool types.
#   pool_type = "DP"

#   # Optional detail flags
#   #
#   # include_detail_info:
#   #   When true, the provider asks the storage API for additional detail blocks via
#   #   the `detailInfoType` query parameter. In this provider, it expands the request
#   #   with these detailInfoType values:
#   #     - FMC
#   #     - tierPhysicalCapacity
#   #     - efficiency
#   #     - formattedCapacity
#   #     - autoAddPoolVol
#   #     - tierDiskType
#   #   This can populate extra nested/advanced fields (when supported by the storage
#   #   firmware/API), and may increase API response size.
#   #
#   include_detail_info = true

#   # include_cache_info:
#   #   When true, the provider adds `class` to detailInfoType.
#   #   This can populate cache/class related fields (when supported by the storage).
#   include_cache_info = true
#   #
#   # Note:
#   #   You can set both flags together. If both are false/omitted, the provider
#   #   returns the baseline fields.
# }

# output "dynamicpools_with_filters" {
#   value = data.hitachi_vsp_dynamic_pools.dynamicpools_with_filters
# }


# For retrieving mainframe-specific dynamic pools, uncomment the following block:
# data "hitachi_vsp_dynamic_pools" "mainframe_pools" {
#   serial      = 12345
#   is_mainframe = true
# }

# output "mainframe_dynamic_pools" {
#   value = data.hitachi_vsp_dynamic_pools.mainframe_pools.dynamic_pools
# }

