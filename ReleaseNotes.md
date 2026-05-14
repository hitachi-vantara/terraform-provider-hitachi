Release Notes for Hitachi Virtual Storage Platform One Block Storage Provider for HashiCorp Terraform 2.4.

Version 2.4 focuses on improving error handling and introduces limited support for VSP One Block 85.
This Version also adds support for NVMe, Thin Image Snapshots, Thin Image Advanced Snapshots, vClone

Version 2.4 also has the following resources and data-sources modules added.


New supported Modules [VSP One and VSP E-Series]: resources
	- hitachi_vsp_nvme
	- hitachi_vsp_snapshot
	- hitachi_vsp_snapshot_group


New supported Modules [VSP One and VSP E-Series]: data-sources
	- hitachi_vsp_nvme
	- hitachi_vsp_nvmes
	- hitachi_vsp_snapshot
	- hitachi_vsp_snapshot_family
	- hitachi_vsp_snapshot_group
	- hitachi_vsp_snapshot_groups
	- hitachi_vsp_snapshots
	- hitachi_vsp_vclone_parent_vols


**Known issues**

| Defect ID | Problem | Workaround |
|-----------|---------|------------|
| UCT-220 | While running Terraform modules, you might see the following error: Error 503 service unavailable because the service might be temporarily busy. | Wait a few minutes and then try to issue the request again | 
| UCT-222 | While running the hitachi_vosb_storage_drives module, sometimes, vendor_name is shown as "N/A". The reason is that if a valid vendor name cannot be obtained, "N/A" is shown. | Currently, there is no workaround.|
| UCT-431 | For VSP One Block 85, storage capacity is currently reported as 0. This is a known issue in the existing storage microcode. | This is a known issue in the existing storage microcode and is planned to be addressed in a future microcode release.|
| UCT-864 | For the hitachi_vsp_dynamic_pools and hitachi_vsp_parity_groups data-source modules, the include_cache_info  = true input parameter is not supported. | Currently, there is no workaround. |