# Hitachi VSP Volumes Data Retrieval from VSP Direct Connect
#
# The "hitachi_vsp_volumes" data source retrieves information about multiple
# volumes based on a specified range of logical device IDs (LDEVs). This allows
# you to access configuration and property details for all volumes within that range.
#
# Configure the parameters (serial, start_ldev_id, end_ldev_id, undefined_ldev)
# according to your environment. You may specify the LDEV range using either
# decimal IDs or hexadecimal strings.
#
# Setting "undefined_ldev" to true will include LDEVs that do not exist on the
# system; setting it to false will return only defined volumes.
#

data "hitachi_vsp_volumes" "volumes1" {
  serial         = 12345
  start_ldev_id  = 280
  end_ldev_id    = 281
  undefined_ldev = false
}

output "volumes1" {
  value = data.hitachi_vsp_volumes.volumes1
}
