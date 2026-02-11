# Hitachi VSP One SDS Block: Generate and Download Configuration File
#
# The `hitachi_vosb_configuration_file` resource allows you to create or download configuration files
# for a VSP One SDS Block system. It supports bare-metal and multiple cloud providers (AWS, Google Cloud, Azure),
# and can be used for maintenance operations such as adding or replacing storage nodes and drives.

### General Parameters
# - `vosb_address`: IP or hostname of VSP One SDS Block (**Required**).
# - `download_existconfig_only`: If `true`, skips creation and only downloads the latest file. Requires `download_path`.
# - `create_only`: If `true`, creates the file but skips downloading.
# - `download_path`: Path to save the downloaded file.
#   - Can be a directory or a full file path.
#   - `.tar.gz` is added if no extension is given for full file paths.
# - `create_configuration_file_param`: Block of `create` parameters for cloud providers. Not for bare-metal.

### Expected Cloud Provider Behavior
# - `expected_cloud_provider`: Validates input combinations only. Accepted values: `"google"`, `"aws"`, `"azure"`.
# - If the value does not match the actual environment, the request still proceeds.
# - VSP One SDS Block always applies its actual cloud provider behavior.
# - Additional parameters may be required based on the selected maintenance operation (`export_file_type`).
# - If not provided, all cloud-specific parameters are ignored.

### Maintenance Operation Types (`export_file_type`)
# - `"Normal"` *(default)*: No additional parameters required.
# - `"AddStorageNodes"`: See cloud-specific requirements below.
# - `"ReplaceStorageNode"`: See cloud-specific requirements below.
# - `"AddDrives"`: See cloud-specific requirements below.
# - `"ReplaceDrive"`: See cloud-specific requirements below.
# - Note: Exact parameters depend on the cloud provider (AWS, Azure, GCP).

### Configuration File Parameters
# For details, refer to the provider documentation, resource `.md` file in `docs/`.
# - `expected_cloud_provider`: `"google"`, `"azure"`, `"aws"`.
# - `export_file_type`: `"Normal"`, `"AddStorageNodes"`, `"ReplaceStorageNode"`, `"AddDrives"`, `"ReplaceDrive"`.
# - `machine_image_id`: VM image ID for node operations.
# - `number_of_drives`: Number of drives to add (6–24 for `AddDrives`).
# - `recover_single_drive`: Replaces a single removed drive if `true`.
# - `drive_id`: UUID of drive to replace. Not allowed if `recover_single_drive = true`.
# - `recover_single_node`: Recovers a node if `true`.
# - `node_id`: UUID of node to replace (required for `ReplaceStorageNode`).
# - `address_setting`: List (1–6 items) of storage node IPs for `AddStorageNodes`. Ignored for AWS.
# - `template_s3_url`: AWS only. URL of S3 bucket storing the configuration file.

### Outputs
# - `status`: Operation status returned by the VSP One SDS Block system.
# - `output_file_path`: Local path of the downloaded configuration file (if downloaded).

### Notes
# - Fields not applicable to the selected cloud provider or `export_file_type` are ignored.
# - Input combinations are validated at runtime.

# =====================================================================
# Bare-metal
# - Usage Modes: See `General Parameters` above.
# - No additional parameters needed. Cloud-specific parameters are ignored.

# =====================================================================
# AWS
# - Usage Modes: same as above.
# - For normal create, it needs create_configuration_file_param
# - Required:
#   - `create_configuration_file_param`
#   - `expected_cloud_provider = "aws"`
#   - `template_s3_url`
# - By `export_file_type`:
#   - `"Normal"`: None
#   - `"AddStorageNodes"`: `machine_image_id`
#   - `"ReplaceStorageNode"`: `machine_image_id`
#   - `"AddDrives"`: `number_of_drives` (6–24)
#   - `"ReplaceDrive"`: `machine_image_id`, and either `drive_id` or `recover_single_drive`

# =====================================================================
# Google Cloud (GCP)
# - Usage Modes: same as above.
# - For normal create, no need for create_configuration_file_param
# - Required (for other than normal create):
#   - `create_configuration_file_param`
#   - `expected_cloud_provider = "google"`
# - By `export_file_type`:
#   - `"Normal"`: None
#   - `"AddStorageNodes"`: `machine_image_id`, `address_setting`
#   - `"ReplaceStorageNode"`: `machine_image_id`, `node_id`, optional `recover_single_node`
#   - `"AddDrives"`: `number_of_drives` (6–24)
#   - `"ReplaceDrive"`: `machine_image_id`, and either `drive_id` or `recover_single_drive`

# =====================================================================
# Azure
# - Usage Modes: same as above.
# - For normal create, no need for create_configuration_file_param
# - Required (for other than normal create):
#   - `create_configuration_file_param`
#   - `expected_cloud_provider = "azure"`
# - By `export_file_type`:
#   - `"Normal"`: None
#   - `"AddStorageNodes"`: `machine_image_id`, `address_setting`. Optional: `compute_port_ipv6_address` in address_setting.
#   - `"ReplaceStorageNode"`: `machine_image_id`
#   - `"AddDrives"`: `number_of_drives` (6–24)
#   - `"ReplaceDrive"`: None

# =====================================================================
### Examples
# For more examples, see the cloud-specific example files.

# Example: Create and download configuration file. 
resource "hitachi_vosb_configuration_file" "download" {
  vosb_address              = var.vosb_address
  download_existconfig_only = false
  download_path             = "."
  create_only               = false
}

output "download_output" {
  value = resource.hitachi_vosb_configuration_file.download
}
