---
# generated by https://github.com/fbreckle/terraform-plugin-docs
page_title: "hitachi_vss_block_dashboard Data Source - terraform-provider-hitachi"
subcategory: "VSS Block Dashboard"
description: |-
  Obtains the information about Dashboard Information.
---

# hitachi_vss_block_dashboard (Data Source)

Obtains the information about Dashboard Information.

## Example Usage

```terraform
data "hitachi_vss_block_dashboard" "dashboard" {
  vss_block_address = "10.10.12.13"
}

output "dashboardoutput" {
  value = data.hitachi_vss_block_dashboard.dashboard
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `vss_block_address` (String) The host name or the IP address (IPv4) of the REST API server on Virtual Storage Software block.

### Read-Only

- `dashboard_info` (List of Object) This is output schema (see [below for nested schema](#nestedatt--dashboard_info))
- `id` (String) The ID of this resource.

<a id="nestedatt--dashboard_info"></a>
### Nested Schema for `dashboard_info`

Read-Only:

- `compute_node_count` (Number)
- `compute_port_count` (Number)
- `data_reduction` (Number)
- `drive_count` (Number)
- `fault_domain_count` (Number)
- `free_capacity_gb` (Number)
- `free_capacity_mb` (Number)
- `health_status` (List of Object) (see [below for nested schema](#nestedobjatt--dashboard_info--health_status))
- `storage_node_count` (Number)
- `storage_pool_count` (Number)
- `total_capacity_gb` (Number)
- `total_capacity_mb` (Number)
- `total_efficiency` (Number)
- `used_capacity_gb` (Number)
- `used_capacity_mb` (Number)
- `volume_count` (Number)

<a id="nestedobjatt--dashboard_info--health_status"></a>
### Nested Schema for `dashboard_info.health_status`

Read-Only:

- `status` (String)
- `type` (String)


