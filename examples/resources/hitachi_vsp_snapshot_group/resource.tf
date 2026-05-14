################################################################################
# Hitachi VSP Snapshot Group Management Resource
#
# This section defines resource blocks to manage the lifecycle of Snapshot Groups.
# It acts as a container-level manager for both Thin Image (TI) Standard and 
# Thin Image Advanced (TIA) pairs grouped together for consistency.
#
# The resource "hitachi_vsp_snapshot_group" supports mass operations like split, 
# resync, restore, clone, vclone, vrestore on all pairs within the named group.
#
# -----------------------------------------------------------------------------------
# Logic Matrix: Snapshot Group Actions
# -----------------------------------------------------------------------------------
# Attribute                    TI Standard                    TI Advanced (TIA)
# -----------------------------------------------------------------------------------
# state = "read"               Refreshes current status       Refreshes current status
# state = "split"              Suspends all pairs in group    Suspends all pairs (S-VOLs usable)
# state = "resync"             Updates S-VOLs from P-VOLs     Updates S-VOLs from P-VOLs
# state = "restore"            Overwrites P-VOL from S-VOL    Overwrites P-VOL (Standard)
# state = "clone"              Triggers physical copies       Not Supported (Use vclone)
# state = "vclone"             N/A                            Triggers virtual copy
# state = "vrestore"           N/A                            Rapid Metadata Restore
# -----------------------------------------------------------------------------------
#
# Configuration Requirements:
#
#   - state: Defines the operation to apply to the group. 
#            Defaults to "read" (safe monitoring state).
#
#   - snapshot_group_name: Always required. Must exist on the storage system.
#
#   - serial: The serial number of the target VSP storage system.
#
#   - copy_speed: (Optional) Only applicable for "clone" state. 
#                 Influences I/O impact: slower, medium, or faster.
#
#   - auto_split: (Optional) Used with "resync" or "restore" to automatically 
#                 transition the group to a split state upon completion.
#
#   - retention_period_hours: (Optional) 
#                 1. Used with "split": Sets a lock on the snapshot data.
#                 2. Used with "resync": Requires auto_split=true to set lock 
#                    after the automatic split occurs.
#
# Notes:
#   - TIA Restore Behavior: When restoring a TIA group, the status transitions to "PSUS".
#   - TIA VRestore Behavior: User is responsible for pairing back vclones to 
#     snapshots in snapshot groups before running vrestore.
#   - terraform destroy command deletes the snapshot group.
#
################################################################################

################################################################################
# Example: Import (existing snapshot group)
# -----------------------------------------------------------------------------
# Use terraform import when the snapshot group already exists on storage and
# you want to bring it under Terraform management without re-creating it.
# The import reads the current state from storage and writes it to tfstate.
#
# Import ID format: <serial>/<snapshot_group_name>
#
# terraform import hitachi_vsp_snapshot_group.imported '12345/SG_PROD_DATABASE'

# Minimal skeleton block required for `terraform import`.
# Create this block in your .tf before running the import command.
resource "hitachi_vsp_snapshot_group" "imported" {
  serial              = 12345
  snapshot_group_name = "SG_PROD_DATABASE"
}

output "imported_snapshot_group_id" {
  value = hitachi_vsp_snapshot_group.imported.id
}

################################################################################
# Example 1: Standard Snapshot Group Operations
# Pattern: state="split" or "resync" or "restore"
# Result: Performs the requested action on all pairs associated with the group name.
# Use Case: Taking a point-in-time snapshot of a database for backup.
################################################################################
resource "hitachi_vsp_snapshot_group" "prod_db_group" {
  serial              = 54321
  snapshot_group_name = "SG_PROD_DATABASE"

  # Set to "split" to suspend the pairs and make S-VOLs accessible
  state = "split"
}

output "snapshot_group_prod" {
  value = hitachi_vsp_snapshot_group.prod_db_group
}

################################################################################
# Example 2: Thin Image Advanced (TIA) with Auto-Retention
# Pattern: state="resync" + auto_split=true + retention_period_hours
# Result: Synchronizes the group, then automatically splits it and sets 
# a 24-hour expiration lock on the TIA snapshots.
# Use Case: Automated daily recovery points that cannot be deleted prematurely.
################################################################################
resource "hitachi_vsp_snapshot_group" "tia_retention_mgmt" {
  serial              = 54321
  snapshot_group_name = "SG_DAILY_RETENTION"
  state               = "resync"

  auto_split = true
  # Enable TIA automation features
  retention_period_hours = 24 # 1 day lock
}

output "snapshot_group_tia" {
  value = hitachi_vsp_snapshot_group.tia_retention_mgmt
}

################################################################################
# Example 3: Snapshot Group Restore with Auto-Split
# Pattern: state="restore" + auto_split=true
# Result: Overwrites the data on the P-VOLs using the S-VOL data in the group. 
#         Once the restoration is initiated/completed, the pairs are 
#         automatically split to allow host I/O to the P-VOL.
# Use Case: Rolling back a volume group to a known good state after a 
#           failed software deployment or data corruption.
################################################################################
resource "hitachi_vsp_snapshot_group" "db_rollback" {
  serial              = 54321
  snapshot_group_name = "SG_PROD_DATABASE"
  state               = "restore"

  # Automatically split the pairs after the restore operation 
  # to resume normal host operations on the P-VOL.
  auto_split = true
}

output "snapshot_group_db" {
  value = hitachi_vsp_snapshot_group.db_rollback
}

################################################################################
# Example 4: Full Physical Clone Operation (TI Standard)
# Pattern: state="clone" + copy_speed
# Result: Triggers a background physical copy of all volumes in the group.
# Use Case: Creating a full performance-isolated copy for Dev/Test.
################################################################################
resource "hitachi_vsp_snapshot_group" "clone_operation" {
  serial              = 54321
  snapshot_group_name = "SG_DEV_RESTORE"
  state               = "clone"

  # Ensure the background copy doesn't impact Production I/O heavily
  copy_speed = "slower"
}

output "snapshot_group_clone" {
  value = hitachi_vsp_snapshot_group.clone_operation
}

################################################################################
# Example 5: Thin Image Advanced (TIA) Virtual Clone (vclone)
# Logic: Triggers an instantaneous, pointer-based virtual clone.
################################################################################
resource "hitachi_vsp_snapshot_group" "analytics_vclone" {
  serial              = 54321
  snapshot_group_name = "SG_PROD_DB_TIA"

  # Action: Create a TIA Virtual Clone (TIA only)
  state = "vclone"
}

output "snapshot_group_vclone" {
  value = hitachi_vsp_snapshot_group.analytics_vclone
}

################################################################################
# Example 6: Rapid TIA Recovery (vrestore)
# Logic: Reverts the P-VOL to a snapshot state using TIA metadata pointers.
# Note: User is responsible for pairing back vclones to snapshots before running vrestore
################################################################################
resource "hitachi_vsp_snapshot_group" "tia_recovery" {
  serial              = 54321
  snapshot_group_name = "SG_APP_DB"
  state               = "vrestore"
}

output "snapshot_group_vrestore" {
  value = hitachi_vsp_snapshot_group.tia_recovery
}

################################################################################
# Example 7: Thin Image Advanced (TIA) - Update with Retention
# Pattern: state="update_retention" + retention_period_hours
# Note: Updates retention period when snapshots in group are already in split
################################################################################
resource "hitachi_vsp_snapshot_group" "tia_update_with_lock" {
  serial                 = 54321
  snapshot_group_name    = "SG_APP_DB"
  state                  = "update_retention"
  retention_period_hours = 1
}

output "snapshot_group_retention" {
  value = hitachi_vsp_snapshot_group.tia_update_with_lock
}

