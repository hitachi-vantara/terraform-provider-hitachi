# Hitachi VSP Volumes Data Retrieval from VSP Direct Connect
#
# The data source block "hitachi_vsp_volumes" retrieves details about volumes within a
# specified range of logical device IDs (LDEVs). This allows you to access configuration
# and property information for the specified volumes.
#
# Customize the values of the parameters (serial, start_ldev_id, end_ldev_id, undefined_ldev)
# to match your environment. By doing so, you can retrieve information about the desired range
# of volumes while indicating whether undefined LDEVs should be included or not.
#

data "hitachi_vsp_volumes" "volumes1" {
  serial         = 12345
  start_ldev_id  = 280
  end_ldev_id    = 285
  undefined_ldev = false
}

output "volumes1" {
  value = data.hitachi_vsp_volumes.volumes1
}
