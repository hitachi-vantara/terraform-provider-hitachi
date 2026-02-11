# Hitachi VSP Storage Maintenance Resource
#
# This example demonstrates the `hitachi_vsp_storage_maintenance` resource which
# can be used to invoke maintenance actions on a Hitachi VSP appliance. 
#
# Resource behavior (example):
# - Create: If `should_stop_all_volume_format` is true, the provider will
#   invoke the appliance action to stop all in-progress volume format operations.
# - Read: Verifies presence of the sentinel state; this resource acts as an
#   invocation sentinel rather than a persistent server-side object.
# - Delete: No backend action is taken on destroy; the resource simply clears
#   the local sentinel state.
#


resource "hitachi_vsp_storage_maintenance" "stop_formats" {
  serial = 12345
  // This only stops format for NORMAL format type
  should_stop_all_volume_format = true
}