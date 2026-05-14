################################################################################
# Hitachi VSP NVM Subsystem Management Resource
#
# This section defines resource blocks to manage the lifecycle of NVM Subsystems
# for NVMe-oF (FC-NVMe or NVMe over TCP) environments.
#
# The resource "hitachi_vsp_nvme" supports creating subsystems, 
# assigning ports, registering host NQNs, and managing namespaces (LDEVs)
# with specific host-to-namespace paths.
#
# -----------------------------------------------------------------------------------
# Attribute Detail & Constraints:
# -----------------------------------------------------------------------------------
# serial                    (Required) Storage system serial number.
# nvm_subsystem_id          (Optional) NVMe ID (0-2047). If omitted, the lowest 
#                           available ID is auto-assigned.
# virtual_nvm_subsystem_id  (Optional) Applicable for VSP 5000 series to support
#                           Virtual Storage Machines (VSM). Range: 0-2047.
# nvm_subsystem_name        (Optional) 1-32 characters. Cannot start/end with 
#                           spaces or start with a hyphen.
# host_mode                 (Optional) Defines the OS communication protocol.
#                           Options: LINUX/IRIX, VMWARE, VMWARE_EX, AIX.
# host_mode_options         (Optional) List of supported integer HMOs. 
#                           Currently, No HMOs are available on the NVM subsystem.
#                           Refer to Provisioning Guide for Open Systems.
# enable_namespace_security (Optional) Boolean. Toggles port-level security for 
#                           Namespaces (similar to LUN security). Default true.
# ports                     (Optional) Set of Port IDs (e.g., ["CL1-A"]).
#                           Elements added to the set are mapped to the subsystem.
#                           Elements removed from the set are unmapped.
# host_nqns                 (Optional) Block defining Host NQNs. 
#                           Sub-attributes: 'nqn', 'nickname'.
#                           New NQN blocks are registered to the subsystem.
#                           Removed NQN blocks are unmapped and deregistered.
# namespaces                (Optional) Block defining LDEV mapping.
#                           Sub-attributes: 'ldev_id/ldev_id_hex', 'nickname', 'path'.
#                           New 'ldev_id' blocks create a namespace mapping.
#                           Removed 'ldev_id' blocks delete the namespace mapping.
#                           'path' (Set of NQNs) manages specific host access to the namespace.
#                           For B series, ldev id should be a DRS.
# -----------------------------------------------------------------------------------
# Terraform Import Procedure:
#
# Step 1: Define a 'skeleton' resource block in your .tf file. 
#         At minimum, the 'serial' must be defined for config-assisted imports.
#
#         resource "hitachi_vsp_nvm_subsystem" "my_nvme" {
#           serial = 810045
#         }
#
# Step 2: Run the import command.
#         terraform import hitachi_vsp_nvm_subsystem.example <serial>/<nvm_subsystem_id>
#         Example:         
#         terraform import hitachi_vsp_nvm_subsystem.my_nvme 810045/15
#
# Step 3: Run 'terraform plan'. 
#         Terraform will output the current configuration as it exists on the 
#         storage hardware.
#
# Step 4: Synchronize your configuration. 
#         Copy the 'ports', 'host_nqns', 'namespaces', and other reported 
#         attributes from the plan output into your .tf resource block.
################################################################################
# Minimal skeleton block required for `terraform import`.
# Create this block in your .tf before running the import command.
resource "hitachi_vsp_nvm_subsystem" "imported" {
  serial = 810045
}

output "imported_nvme_id" {
  value = hitachi_vsp_nvm_subsystem.imported.id
}
################################################################################
# Example: Basic NVMe Auto-assigned ID
################################################################################
resource "hitachi_vsp_nvme" "basic_nvme" {
  serial = 810045
}

output "basic_nvme_out" {
  value = hitachi_vsp_nvme.basic_nvme.nvm_subsystems
}

################################################################################
# Example: Basic NVMe Explicit ID
################################################################################
resource "hitachi_vsp_nvme" "basic_nvme_explicit" {
  serial                    = 810045
  nvm_subsystem_id          = 15
  nvm_subsystem_name        = "APP_SERVER_NVME"
  host_mode                 = "VMWARE_EX"
  enable_namespace_security = false
}

output "basic_nvme_expl_out" {
  value = hitachi_vsp_nvme.basic_nvme_explicit.nvm_subsystems
}

################################################################################
# Example: NVMe with ports, host, namespace
################################################################################
resource "hitachi_vsp_nvme" "nvme_info" {
  serial = 810045

  ports = ["CL1-A", "CL2-A"]

  host_nqns {
    nqn      = "nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000"
    nickname = "PROD_HOST_01"
  }

  namespaces {
    ldev_id  = 304 # or ldev_id_hex = "0x..."
    nickname = "GOLD_VOL_01"
    path     = ["nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000"]
  }
}

output "nvme_info_out" {
  value = hitachi_vsp_nvme.nvme_info.nvm_subsystems
}

################################################################################
# Example: High Availability Subsystem with Multiple Namespaces
################################################################################
resource "hitachi_vsp_nvme" "ha_subsystem" {
  serial           = 810045
  nvm_subsystem_id = 15
  host_mode        = "VMWARE_EX"

  ports = ["CL1-A", "CL1-B", "CL2-A", "CL2-B"]

  # Registering multiple hosts for a shared cluster
  host_nqns {
    nqn      = "nqn.2014-08.com.vmware:nvme:esx01"
    nickname = "ESX_NODE_A"
  }
  host_nqns {
    nqn      = "nqn.2014-08.com.vmware:nvme:esx02"
    nickname = "ESX_NODE_B"
  }

  # Mapping multiple LDEVs (Namespaces) to both hosts
  namespaces {
    ldev_id  = 1097 # or ldev_id_hex = "0x..."
    nickname = "VM_DATA_POOL_01"
    path     = ["nqn.2014-08.com.vmware:nvme:esx01", "nqn.2014-08.com.vmware:nvme:esx02"]
  }

  namespaces {
    ldev_id  = 1098 # or ldev_id_hex = "0x..."
    nickname = "VM_DATA_POOL_02"
    path     = ["nqn.2014-08.com.vmware:nvme:esx01", "nqn.2014-08.com.vmware:nvme:esx02"]
  }
}

output "ha_nvme" {
  value = hitachi_vsp_nvme.ha_subsystem.nvm_subsystems
}

################################################################################
# Example: Subsystem for VSP 5000 (Virtual Storage Machine)
# Pattern: Uses virtual_nvm_subsystem_id for multi-tenancy.
################################################################################
resource "hitachi_vsp_nvme" "vsm_nvme" {
  serial                    = 54321
  virtual_nvm_subsystem_id  = 100
  nvm_subsystem_name        = "VSM_TENANT_A"
  enable_namespace_security = true

  ports = ["CL3-A"]

  host_nqns {
    nqn = "nqn.2014-08.org.nvmexpress:uuid:tenant-a-host"
  }

  namespaces {
    ldev_id = 5000 # or ldev_id_hex = "0x..."
    path    = ["nqn.2014-08.org.nvmexpress:uuid:tenant-a-host"]
  }
}

output "vsm_nvme_out" {
  value = hitachi_vsp_nvme.vsm_nvme.nvm_subsystems
}

################################################################################
# PORT MANAGEMENT EXAMPLES
################################################################################

# Example A: Adding a Port
# Append the ID to the list. Terraform maps the new port to the subsystem.
resource "hitachi_vsp_nvme" "add_port" {
  serial = 810045
  ports  = ["CL1-A", "CL2-A", "CL3-A"] # Existing ports were ["CL1-A", "CL2-A"]
}

# Example B: Removing a Port
# Remove the ID from the set. Terraform unmaps only the missing port.
resource "hitachi_vsp_nvme" "remove_port" {
  serial = 810045
  ports  = ["CL1-A"]
}

# Example C: Deleting ALL Ports
# Set to an empty list. The provider unmaps every port currently assigned.
resource "hitachi_vsp_nvme" "clear_all_ports" {
  serial = 810045
  ports  = []
}

# Example D: Ignoring Port Changes
# Omitting the attribute entirely prevents Terraform from managing port state.
resource "hitachi_vsp_nvme" "ignore_ports" {
  serial = 810045
  # ports attribute is absent; do nothing.
}

################################################################################
# HOST NQN MANAGEMENT EXAMPLES
################################################################################

# Example A: Registering a new Host NQN
# Add a new host_nqns block. Terraform registers the NQN to the subsystem.
resource "hitachi_vsp_nvme" "add_host" {
  serial = 810045

  host_nqns {
    nqn      = "nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000"
    nickname = "PROD_HOST_01"
  }

  # Adding a second block registers an additional host
  host_nqns {
    nqn      = "nqn.2014-08.org.nvmexpress:uuid:661f9511-f30c-52e5-b827-557766551111"
    nickname = "PROD_HOST_02"
  }
}

# Example B: Removing a Host NQN
# Delete the specific host_nqns block. Terraform deregisters that NQN.
resource "hitachi_vsp_nvme" "remove_host" {
  serial = 810045

  # PROD_HOST_02 block has been removed; it will be deleted from the array.
  host_nqns {
    nqn      = "nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000"
    nickname = "PROD_HOST_01"
  }
}

# Example C: Deleting ALL Host NQNs
# Set to an empty list. The provider will deregister every Host NQN currently in the subsystem.
resource "hitachi_vsp_nvme" "clear_all_hosts" {
  serial    = 810045
  host_nqns = []
}

# Example D: Ignoring Host NQN Changes
# Omitting the attribute entirely prevents Terraform from managing Host NQN state.
resource "hitachi_vsp_nvme" "ignore_hosts" {
  serial = 810045
  # host_nqns block is absent; do nothing.
}

################################################################################
# NAMESPACE MANAGEMENT EXAMPLES
################################################################################

# Example A: Adding a Namespace (Mapping an LDEV)
# Use 'ldev_id' or 'ldev_id_hex' to map a volume to the subsystem.
resource "hitachi_vsp_nvme" "add_namespace" {
  serial = 810045

  # Mapping using Decimal ID
  namespaces {
    ldev_id  = 300
    nickname = "DB_VOL_01"
    path     = ["nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000"]
  }

  # Mapping using Hexadecimal ID
  namespaces {
    ldev_id_hex = "01:2D" # Decimal 301
    nickname    = "DB_VOL_02"
    path        = ["nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000"]
  }
}

# Example B: Updating Namespace Access (Path Management)
# Adding or removing NQNs from the 'path' set will map/unmap access paths.
resource "hitachi_vsp_nvme" "update_paths" {
  serial = 810045

  namespaces {
    ldev_id  = 300
    nickname = "DB_VOL_01"
    # Added a second Host NQN to the path list
    path = [
      "nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000",
      "nqn.2014-08.org.nvmexpress:uuid:661f9511-f30c-52e5-b827-557766551111"
    ]
  }
}

# Example C: Removing a Namespace
# Removing the block will unmap the paths first, then delete the namespace mapping.
resource "hitachi_vsp_nvme" "remove_namespace" {
  serial = 810045

  # DB_VOL_02 block was removed. 
  # Terraform will automatically clear its paths before unmapping the LDEV.
  namespaces {
    ldev_id  = 300
    nickname = "DB_VOL_01"
    path     = ["nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000"]
  }
}

# Example D: Deleting ALL Namespaces
# Setting to an empty list removes all LDEV mappings from the subsystem.
resource "hitachi_vsp_nvme" "clear_all_namespaces" {
  serial     = 810045
  namespaces = []
}

# Example E: Ignoring Namespace Changes
# Omitting the block prevents Terraform from managing LDEV mappings or paths.
resource "hitachi_vsp_nvme" "ignore_namespaces" {
  serial = 810045
  # namespaces attribute is absent; existing array mappings are preserved.
}

################################################################################
# NAMESPACE PATH (MAPPING) MANAGEMENT EXAMPLES
################################################################################

# Example A: Mapping to Multiple Hosts
# Defining multiple NQNs in the path attribute creates a mapping for each host.
resource "hitachi_vsp_nvme" "multi_path" {
  serial = 810045

  namespaces {
    ldev_id  = 300
    nickname = "SHARED_VOL"
    path = [
      "nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000",
      "nqn.2014-08.org.nvmexpress:uuid:661f9511-f30c-52e5-b827-557766551111"
    ]
  }
}

# Example B: Removing a Single Path
# Removing one NQN from the list triggers a deregistration of that specific path.
resource "hitachi_vsp_nvme" "remove_single_path" {
  serial = 810045

  namespaces {
    ldev_id  = 300
    nickname = "SHARED_VOL"
    # Host ...1111 removed; only ...0000 remains mapped.
    path = ["nqn.2014-08.org.nvmexpress:uuid:550e8400-e29b-41d4-a716-446655440000"]
  }
}

# Example C: Deleting ALL Paths (Unmapping from all Hosts)
# IMPORTANT: Setting path to [] or omitting it entirely inside a namespace block 
# will cause Terraform to delete all existing mappings for that LDEV.
# There is no "retain" behavior for paths within a managed namespace.
resource "hitachi_vsp_nvme" "clear_all_paths" {
  serial = 810045

  namespaces {
    ldev_id  = 300
    nickname = "ISOLATED_VOL"
    path     = [] # This will unmap the LDEV from every host NQN
    # or
    # path attribute is absent does delete all path
  }
}
