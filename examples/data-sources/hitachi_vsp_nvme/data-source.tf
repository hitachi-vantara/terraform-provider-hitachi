#
# Hitachi VSP NVM Subsystem Data Retrieval
#
# This section defines a data source block to fetch information about a specific
# NVM Subsystem from a Hitachi Virtual Storage Platform (VSP) using
# HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vsp_nvme" retrieves details about the NVM 
# subsystem associated with the provided parameters. This allows you to access 
# configuration, associated namespaces (LDEVs), host NQNs, and path mappings.
#
# Provide the storage system serial number and the NVM subsystem ID.
#

#####################
# Example: NVM Subsystem Lookup
data "hitachi_vsp_nvme" "nvme" {
  serial           = var.serial_number
  nvm_subsystem_id = 1
}

output "nvme_subsystems" {
  description = "List of NVM subsystem info"
  value       = data.hitachi_vsp_nvme.nvme.nvm_subsystems
}

output "nvme_count" {
  description = "Total number of NVM subsystems returned"
  value       = data.hitachi_vsp_nvme.nvme.nvm_subsystem_count
}
