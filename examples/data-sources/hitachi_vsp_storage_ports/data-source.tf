#
# Hitachi VSP Storage Ports Data Retrieval
#
# This section defines a data source block to fetch information about a specific
# storage port from a Hitachi Virtual Storage Platform (VSP) using HashiCorp
# Configuration Language (HCL).
#
# The data source block "hitachi_vsp_storage_ports" retrieves details about a
# storage port associated with the provided parameters. This allows you to access
# configuration and property information for the specified storage port.
#
# Customize the values of the parameters (serial, port_id) to match your
# environment, enabling you to retrieve information about the desired storage port.
#

data "hitachi_vsp_storage_ports" "storageports" {
  serial  = 12345
  port_id = "CL4-C"
}

output "storageports" {
  value = data.hitachi_vsp_storage_ports.storageports
}

# # Storage ports with additional filters
# data "hitachi_vsp_storage_ports" "storageports_with_filters" {
#   serial = 12345

#   # Include detailed information for ports
#   # (Optional, default: false)
#   #   When true, the provider requests additional “detailInfoType” data for ports.
#   #   This can populate extra fields that are not returned in the baseline ports
#   #   response 
#   include_detail_info = true

#   # Optional filters (do NOT set port_id when using these)
#   #   Server-side filter to return only ports of the specified type.
#   #   Allowed values:
#   #     - FIBRE      : Fibre Channel ports
#   #     - SCSI       : SCSI ports
#   #     - ISCSI      : iSCSI ports
#   #     - NVME_TCP   : NVMe/TCP ports
#   #     - ENAS       : NAS-related ports
#   #     - ESCON      : ESCON ports
#   #     - FICON      : FICON ports
#   #   If omitted, ports of all types are returned.
#   port_type = "FIBRE"

#   #   Server-side filter to return only ports matching an attribute.
#   #   Allowed values:
#   #     - TAR  : Target port
#   #     - MCU  : Initiator port
#   #     - RCU  : RCU target port
#   #     - ELUN : External port
#   #   If omitted, ports of all attributes are returned.
#   #   Cannot be used with port_id.
#   port_attributes = "TAR"
# }

# output "storageports_with_filters" {
#   value = data.hitachi_vsp_storage_ports.storageports_with_filters
# }
