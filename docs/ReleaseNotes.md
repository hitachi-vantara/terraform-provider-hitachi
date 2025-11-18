Release Notes for Hitachi Virtual Storage Platform One Block Storage Provider for HashiCorp Terraform 2.1.2.

Version 2.1.2 focuses on improving error handling and introduces limited support for VSP One Block 85 only over Fibre Channel (FC).

No new features or breaking changes are included.

**Known issues**

| Defect ID | Problem | Workaround |
|-----------|---------|------------|
| UCT-220 | While running Terraform modules, you might see the following error: Error 503 service unavailable because the service might be temporarily busy. | Wait a few minutes and then try to issue the request again | 
| UCT-222 | While running the hitachi_vosb_storage_drives module, sometimes, vendor_name is shown as "N/A". The reason is that if a valid vendor name cannot be obtained, "N/A" is shown. | Currently, there is no workaround.|
| UCT-431 | For VSP One Block 85, storage capacity is currently reported as 0. This is a known issue in the existing storage microcode. | This is a known issue in the existing storage microcode and is planned to be addressed in a future microcode release.|
| UCT-581 | User receives an Intermittent “401 Unauthorized” error with data source hitachi_vsp_volumes with VSP One Block storage systems | The workaround is to wait a few seconds (around 15 seconds) before trying again. |
| UCT-582 | Add path fails on VSP One Block when capacity saving is disabled.| For VSP One Block storage systems, before executing the Add LUN to host group task, enable the “capacity saving“ property on the LDEV. [“hitachi_vsp_hostgroup“ module]. The “Enable Capacity Saving” option is not supported for volume creation in this release (v2.1.2). This behavior is consistent with previous versions and is currently by design. Support for this capability is being evaluated for a future release. |



