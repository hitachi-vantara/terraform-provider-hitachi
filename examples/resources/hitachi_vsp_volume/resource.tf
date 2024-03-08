//
// Hitachi VSP Volume Resource
//
// This section defines a Terraform resource block to create a Hitachi VSP volume.
// The resource "hitachi_vsp_volume" represents a volume on a Hitachi Virtual Storage
// Platform (VSP) and allows you to manage its configuration using Terraform.
//
// Customize the values of the parameters (serial, size_gb, pool_id) to match your
// desired volume configuration.
//

/* 
Refre to the below table for more information

| Property                      | Gateway Provider           | VSP Direct Connect Provider  |
|-------------------------------|----------------------------|------------------------------|
| serial                        | 12345 **                   | 12345 *                      |
| storage_id                    | ""   **                    | -                            |
| size_gb                       | 1  *                       | 1 *                          |
| pool_id                       | 2  *                       | 1 **                         |
| name                          | "SampleName"               | "SampleName"                 |
| system                        | "SampleSystemName"         | -                            |
| subscriber_id                 | ""                         | -                            |
| resource_group_id             | ""                         | -                            |
| ldev_id                       | 2                          | -                            |
| deduplication_compression_mode| ""                         | ""                           |
| pool_name                     | -                          | "PoolName"   **              |
| paritygroup_id                | -                          | "parity_group_id" **         |



- Not supported
** Either or one of the parameter required(Optinal Mandotary)
* Required parameters


User can comment and uncomments based on requirements
*/

resource "hitachi_vsp_volume" "mylun" {
  serial  = 12345
  size_gb = 1
  pool_id = 1

  //Optional parameters in both the provider
  name = "hitachi_vsp_volume"
  deduplication_compression_mode ="DISABLED"
  ldev_id = 0

}
