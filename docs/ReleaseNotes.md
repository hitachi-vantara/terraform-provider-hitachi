Release Notes for Hitachi Virtual Storage Platform One Block Storage Provider for HashiCorp Terraform 2.2.

Version 2.2 focuses on improving error handling and introduces limited support for VSP One Block 85 only over Fibre Channel (FC).
Version 2.2 also has following resources & data-sources modules supported added


New supported Modules [VSP One and VSP E-Series]: resources

- hitachi_vsp_one_iscsi_target
- hitachi_vsp_one_pool
- hitachi_vsp_one_port
- hitachi_vsp_one_server
- hitachi_vsp_one_server_hba
- hitachi_vsp_one_server_path
- hitachi_vsp_one_volume
- hitachi_vsp_one_volume_qos
- hitachi_vsp_one_volume_server_connection

New supported Modules [VSP One and VSP E-Series]: data-sources

- hitachi_vsp_one_iscsi_target
- hitachi_vsp_one_iscsi_targets
- hitachi_vsp_one_pool
- hitachi_vsp_one_pools
- hitachi_vsp_one_port
- hitachi_vsp_one_ports
- hitachi_vsp_one_server
- hitachi_vsp_one_server_hba
- hitachi_vsp_one_server_hbas
- hitachi_vsp_one_server_path
- hitachi_vsp_one_servers
- hitachi_vsp_one_storage
- hitachi_vsp_one_volume
- hitachi_vsp_one_volume_qos
- hitachi_vsp_one_volumes
- hitachi_vsp_one_volume_server_connection
- hitachi_vsp_one_volume_server_connections


**Known issues**

| Defect ID | Problem | Workaround |
|-----------|---------|------------|
| UCT-220 | While running Terraform modules, you might see the following error: Error 503 service unavailable because the service might be temporarily busy. | Wait a few minutes and then try to issue the request again | 
| UCT-222 | While running the hitachi_vosb_storage_drives module, sometimes, vendor_name is shown as "N/A". The reason is that if a valid vendor name cannot be obtained, "N/A" is shown. | Currently, there is no workaround.|
| UCT-431 | For VSP One Block 85, storage capacity is currently reported as 0. This is a known issue in the existing storage microcode. | This is a known issue in the existing storage microcode and is planned to be addressed in a future microcode release.|

