# Hitachi VSP Volumes Data Retrieval from VSP Direct Connect
#
# The "hitachi_vsp_volumes" data source retrieves information about multiple
# volumes based on a specified range of logical device IDs (LDEVs). This allows
# you to access configuration and property details for all volumes within that range.
#
# Configure the parameters (serial, start_ldev_id/start_ldev_id_hex, end_ldev_id/end_ldev_id_hex, undefined_ldev)
# according to your environment. You may specify the LDEV range using either
# decimal IDs or hexadecimal strings.
#
# Setting "undefined_ldev" to true will include LDEVs that do not exist on the
# system; setting it to false will return only defined volumes.
#
# For advanced filtering and additional details, refer to the commented example
# blocks below. These demonstrate how to use options such as "filter_option",
# "include_detail_info", and "include_cache_info" to further customize the
# data retrieval for Hitachi VSP volumes.
#

data "hitachi_vsp_volumes" "volumes1" {
  serial        = 12345
  start_ldev_id = 280
  # or
  # start_ldev_id_hex = "0X118"
  end_ldev_id = 281
  # or
  # end_ldev_id_hex   = "0X119"
  undefined_ldev = false
}

output "volumes1" {
  value = data.hitachi_vsp_volumes.volumes1
}


# # Test VSP multiple volumes data source (legacy) with all advanced parameters
# data "hitachi_vsp_volumes" "multiple_volumes_test" {
#   serial         = 12345
  
#   # LDEV Option: Filters volumes based on specific conditions/types
#   # Available values:
#   #   "defined"        - Only defined/allocated volumes 
#   #   "undefined"      - Only undefined/unallocated volumes
#   #   "dpVolume"       - Dynamic pool volumes only
#   #   "luMapped"       - Only volumes mapped to logical units
#   #   "luUnmapped"     - Only volumes not mapped to logical units
#   #   "externalVolume" - External volumes (from external storage)
#   #   "mappedNamespace" - Volumes mapped to namespaces
#   #   "mainframe"      - Mainframe volumes (emulation types 3390-A and 3390-V)
#   filter_option    = "defined"
  
#   # Include Detail Info: Adds comprehensive volume information
#   # When true, includes: FMC info, external volume details, virtual serial numbers,
#   # data reduction/saving info, QoS settings, NGU ID, and other advanced properties
#   include_detail_info = true
  
#   # Include Cache Info: Adds cache-related information
#   # When true, includes: cache class details and cache-related properties
#   # Can only be specified for VSP 5000 series.
#   include_cache_info  = false
# }

# output "multiple_volumes_result" {
#   value = data.hitachi_vsp_volumes.multiple_volumes_test.volumes
# }
