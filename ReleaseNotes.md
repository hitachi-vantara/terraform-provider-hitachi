Release Notes for Hitachi Virtual Storage Platform One Block Storage Provider for HashiCorp Terraform 2.3.

Version 2.3 focuses on improving error handling and introduces limited support for VSP One Block High End only over Fibre Channel (FC).
This version also adds support for Mainframe z16 and SVOS 10.5.1.
Version 2.3 also has the following resources and data-sources modules added.


New supported Modules [VSP One and VSP E-Series]: resources

- hitachi_vsp_pav_ldev
- hitachi_vsp_storage_maintenance

New supported Modules [VSP One and VSP E-Series]: data-sources

- hitachi_vsp_pav_alias
- hitachi_vsp_supported_host_modes


**Known issues**

| Defect ID | Problem | Workaround |
|-----------|---------|------------|
| UCT-220 | While running Terraform modules, you might see the following error: Error 503 service unavailable because the service might be temporarily busy. | Wait a few minutes and then try to issue the request again | 
| UCT-222 | While running the hitachi_vosb_storage_drives module, sometimes, vendor_name is shown as "N/A". The reason is that if a valid vendor name cannot be obtained, "N/A" is shown. | Currently, there is no workaround.|
| UCT-431 | For VSP One Block High End, storage capacity is currently reported as 0. This is a known issue in the existing storage microcode. | This is a known issue in the existing storage microcode and is planned to be addressed in a future microcode release.|
| UCT-864 | For the hitachi_vsp_dynamic_pools and hitachi_vsp_parity_groups data-source modules, the include_cache_info  = true input parameter is not supported. | Currently, there is no workaround. |

