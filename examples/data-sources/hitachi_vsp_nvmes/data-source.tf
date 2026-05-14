#
# Hitachi VSP NVM Subsystems Data Retrieval
#
# This section defines a data source block to fetch information about multiple 
# NVM Subsystems from a Hitachi Virtual Storage Platform (VSP) using 
# HashiCorp Configuration Language (HCL).
#
# The data source block "hitachi_vsp_nvmes" retrieves a list of 
# subsystems based on the provided filters. This allows you to discover 
# subsystem IDs, names, and optionally enrich the output with detailed 
# port, host NQN, and namespace/path information.
#
# Provide the storage system serial number and optional filters. 
# You can use the "include_ports", "include_host_nqns", and "include_namespaces" 
# inputs to toggle the retrieval of extended configuration data.
#
# Search Patterns (Optional):
#   1. ID Filter:   [nvm_subsystem_id]   -> Matches a specific ID.
#   2. Name Filter: [nvm_subsystem_name] -> Matches by name (case-insensitive).
#   3. No Filter:   Returns all NVM subsystems for the storage system.
#
# Enrichment Inputs (Boolean):
#   - include_ports:      Fetches associated NVM ports for each subsystem.
#   - include_host_nqns:  Fetches authorized host NQNs and mapping paths.
#   - include_namespaces: Fetches LDEV namespace IDs and their Host NQN paths.
#



#####################
# Example: Fetch All Subsystems (Basic)
# Retrieves basic information for every NVM subsystem on the storage system.
data "hitachi_vsp_nvmes" "all" {
  serial = 12345
}

output "all_subsystems" {
  value = data.hitachi_vsp_nvmes.all.nvm_subsystems
}

output "total_count" {
  value = data.hitachi_vsp_nvmes.all.nvm_subsystem_count
}

#####################
# Example: Full Enrichment by Name
# Finds a specific subsystem and includes namespaces, ports, and host NQNs.
data "hitachi_vsp_nvmes" "enriched" {
  serial             = 12345
  nvm_subsystem_name = "PROD_ESXI_CLUSTER"
  include_namespaces = true
  include_ports      = true
  include_host_nqns  = true
}

output "enriched_subsystems" {
  value = data.hitachi_vsp_nvmes.enriched.nvm_subsystems
}

#####################
# Example: Filter by ID with Path Enrichment
# Retrieves a subsystem by ID and includes host NQN path mappings.
data "hitachi_vsp_nvmes" "filtered_id" {
  serial            = 12345
  nvm_subsystem_id  = 10
  include_host_nqns = true
}

output "subsystem_details" {
  value = data.hitachi_vsp_nvmes.filtered_id.nvm_subsystems
}
