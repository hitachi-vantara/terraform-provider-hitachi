#
# Hitachi VSP PAV LDEV (Assign/Unassign Alias)
#
# This example defines a resource block to assign one or more alias LDEVs to a base LDEV
# on a Hitachi Virtual Storage Platform (VSP) using HashiCorp Configuration Language (HCL).
#
# Resource behavior:
# - Create: Assigns the alias LDEVs listed in `alias_ldev_ids` to `base_ldev_id`.
# - Read: Verifies that the configured alias LDEVs remain assigned to the base LDEV.
# - Destroy (Unassign): Unassigns the alias LDEVs listed in `alias_ldev_ids`.
#
# Unassign (Destroy) steps:
# - To unassign alias devices, run:
#   terraform destroy -target=hitachi_vsp_pav_ldev.pav_ldev
# - Or remove the resource block and run `terraform apply`.
#


resource "hitachi_vsp_pav_ldev" "pav_ldev" {
  serial         = 12345
  base_ldev_id   = 100
  alias_ldev_ids = [101, 102]
}

output "pav_ldev_data" {
  value = resource.hitachi_vsp_pav_ldev.pav_ldev
}