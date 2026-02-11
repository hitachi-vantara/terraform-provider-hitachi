# Terraform Provider Hitachi

## Table of Contents
1. [Introduction](#introduction)
2. [Hardware requirements](#hardware-requirements)
3. [Software requirements](#software-requirements)
4. [Prerequisites & Dependencies](#prerequisites--dependencies)
5. [Build the Provider](#build-the-provider)
6. [Install RPM Package](#install-rpm-package)
7. [Check Terraform Providers](#check-terraform-providers)
8. [Uninstall RPM Package](#uninstall-rpm-package)
9. [Directory Layout](#directory-layout)
10. [Storage configuration operations](#storage-configuration-operations)
11. [Logging](#logging)
12. [Log Bundle Collection](#log-bundle-collection)
13. [Known Behavior](#known-behavior)
14. [Usage data collection](#usage-data-collection)
15. [Bundled TF Examples](#bundled-tf-examples)
16. [User Consent Script](#user-consent-script)
17. [Supported host modes](#supported-host-modes)
18. [Host mode options](#host-mode-options)
19. [The clean shell script](#the-clean-shell-script)
20. [Advanced workflows](#advanced-workflows)
21. [Terraform Compatibility](#terraform-compatibility)
22. [Linux Compatibility](#linux-compatibility)
23. [Planning for Mainframe Volume capacity](#planning-for-mainframe-volume-capacity)
24. [Privacy Notice](#privacy-notice)
25. [License](#license)

---

## Introduction

Hitachi Virtual Storage Platform One Block Storage Provider for HashiCorp Terraform 2.3.

Hitachi Virtual Storage Platform One Block Storage Provider for HashiCorp Terraform enables
IT and data center administrators to automate and manage the configuration of Hitachi block
storage systems (VSP One Block, VSP 5000 series, VSP E series, VSP F series, and VSP G
series) and VSP One SDS Block and Cloud systems.

VSP One Block Storage Provider for HashiCorp Terraform is installed on a virtual machine or
a bare metal server. It communicates storage operation requests from the HashiCorp
Terraform provider to the storage system. It accomplishes this communication through VSP
block storage systems and VSP One SDS Block and Cloud systems data source and
resource requests.

The configuration file for VSP One Block Storage Provider for HashiCorp Terraform allows
you to run VSP block storage systems and VSP One SDS Block and Cloud systems data
sources and resources.

**Example of VSP One SDS Block and Cloud Provider configuration file**

```hcl
terraform {
  required_providers {
    hitachi = {
      version = "2.3"
      source = "localhost/hitachi-vantara/hitachi"
    }
  }
}
provider "hitachi" {
  hitachi_vosb_provider {
    vosb_address = "192.0.2.100"
    username = var.hitachi_storage_user
    password = var.hitachi_storage_password
  }
}
```

You must be an admin user to run all the VSP One SDS Block and Cloud
Terraform data sources/resources modules.

**Example configuration for hitachi_vsp_one_provider**

```hcl
terraform {
  required_providers {
    hitachi = {
    version = "2.3.0"
    source  = "localhost/hitachi-vantara/hitachi"
    }
  }
}

provider "hitachi" {
  hitachi_vsp_one_provider {
             serial        = 12345
            management_ip = "10.10.11.12"                             
            username      = var.hitachi_storage_user
            password      = var.hitachi_storage_password
}
 }
```

**Example of VSP Block Provider configuration file**

```hcl
terraform {
  required_providers {
    hitachi = {
      version = "2.3"
      source = "localhost/hitachi-vantara/hitachi"
    }
  }
}
provider "hitachi" {
  san_storage_system {
    serial = 12345
    management_ip = "192.0.2.100"
    username = var.hitachi_storage_user
    password = var.hitachi_storage_password
  }
}    
```

**Schema**

**hitachi_vosb_provider**
  (block list) VSP One SDS Block and Cloud combines VSP One SDS Block, which 
  creates virtual storage systems from general-purpose servers, with VSP One SDS Cloud, 
  which enables deployment on AWS, Google Cloud Platform (GCP), and Microsoft Azure. 
- vosb_address (String) The host name or the IP address (IPv4) of the VSP One
SDS Block and Cloud.
- username (String) User name of the VSP One SDS Block and Cloud
- password (String) Password of the VSP One SDS Block and Cloud

**hitachi_vsp_one_provider**
 (Block List) VSP One Block Administrator is a configuration
management tool designed for VSP One Block 20 series, VSP One Block High End, and VSP
E series storage systems, simplifying and streamlining storage management.
- serial (Number) The serial number for VSP One Block Administrator
- management_ip (String) Management IP for the VSP One Block Administrator
- username (String) User name for the VSP One Block Administrator
- password (String) Password for the VSP One Block Administrator

**san_storage_system**
 (block list) VSP One Block 20 series, VSP One Block High End, VSP 5000 series, 
VSP E series, VSP F series, and VSP G series - all the VSP block storage systems present
in the Hardware requirements table are enterprise and/or mid-range storage solutions designed
to provide reliable and scalable block storage for a variety of environments. These systems focus
on simplifying data storage management while ensuring high availability and data
integrity.
- serial (Number) The serial number of the VSP block storage system
- management_ip (String) Management IP for the VSP block storage system
- username (String) User name for the VSP block storage system
- password (String) Password for the VSP block storage system 

---

## Hardware requirements

| VSP block storage systems | Microcode/Firmware |
|---------------------------|--------------------|
| VSP One Block 24 | A3-04-21-40/00 SVOS 10.4.1 |
| VSP One Block 26 | A3-04-21-40/00 SVOS 10.4.1 |
| VSP One Block 28 | A3-04-21-40/00 SVOS 10.4.1 |
| VSP One Block High End | A0-05-20-00/05 SVOS 10.5.1 |
| VSP 5100, 5500, 5100H, 5500H (SAS) | 90-09-26-00/00 |
| VSP 5200, 5600, 5200H, 5200H (SAS) | 90-09-26-00/00 |
| VSP E590, VSP E790 | 93-07-25-40/00 SVOS 9.8.7 |
| VSP E990 | 93-07-25-60/00 SVOS 9.8.7 |
| VSP E1090 | 93-07-25-80/00 SVOS 9.8.7 |
| VSP F350, VSP F370, VSP F700, FSP F900 | 88-08-14-x0/00 SVOS 9.6.0 |
| VSP G370, VSP G700, VSP G900 | 88-08-14-x0/00 SVOS 9.6.0 |
| VSP G350 | 88-08-15-20/01 |

The listed microcode versions are the minimum versions.

| VSP One SDS Block and Cloud Systems | Storage software version |
|-------------------------------------|--------------------------|
| Bare Metal | 01.18.02.40 |
| Cloud for AWS | 01.18.02.30 |
| Cloud for Microsoft Azure | 01.18.02.50 |
| Cloud for GCP | 01.18.02.60 |

Management of NVMe/TCP connections is not supported in the current release.

The following user roles are required when registering the storage system
user credentials:
- VSP block storage systems: Storage Administrator (Provisioning) roles
- VSP One SDS Block and Cloud systems: Storage roles

---

## Software requirements

| Category | Details |
|----------|---------|
| Software | Terraform - version 1.14.4 or higher |
| Host | VSP One Block Storage Provider for HashiCorp Terraform |
| CPU/vCPUs | 2 |
| Physical Memory | 8 GB |
| Hard Disk | 30 GB |
| Supported Operating Systems | Red Hat: 8.x, 9.1, 9.2, Oracle Enterprise Linux (OEL): 8.x, 9.1, 9.2 |

Hitachi recommends that you update the operating system with the latest software packages.

## Prerequisites & Dependencies

Ensure the following are installed:
- `jq`
- `terraform` (>= 1.14.0)
- `uuidgen`
- Standard runtime libraries (`glibc`, etc.)

---

## Build the Provider

### Build RPM package (for Developers)

Requires **Golang v1.24 or 1.25.x** and **superuser privileges**.

```bash
export GOPATH=/usr/local/go
cd <your hitachi terraform source code directory>
./build.sh <BUILD_NUMBER>
```

Example output path:
```
./rpmbuild/RPMS/x86_64/HV_Storage_Terraform-02.3-19.x86_64.rpm
```

### Build without RPM (for Developers)

```bash
make all
cd examples
terraform init && terraform apply
```

---

## Install RPM Package

Use the following procedure to install and configure VSP One Block Storage Provider for
HashiCorp Terraform.

**Before you begin**
Verify the Prequisites and dependencies

**Procedure**

1. From the Hitachi Vantara Support Portal (https://support.hitachivantara.com/en/anonymous-dashboard.html) 
Downloads page (login credentials required), search for terraform, click Hardware Download, and then 
download the VSP One Block Storage Provider for HashiCorp Terraform file. The Terraform modules are version
2.3.

2. Extract the following file from the distribution media kit installation TAR file:
HV_Storage_Terraform-02.3-XX.x86_64.tar.gz file.

3. Extract the installation rpm file on the Linux server using the following command:

```bash
tar -zxvf HV_Storage_Terraform-02.3-XX.x86_64.tar.gz
```

4. Upload the HV_Storage_Terraform-02.3-XX.x86_64.rpm file to the Linux host
where you want to install VSP One Block Storage Provider for HashiCorp Terraform.

5. Install VSP One Block Storage Provider for HashiCorp Terraform.

- Enter:

```bash
yum localinstall ./HV_Storage_Terraform-02.3-XX.x86_64.rpm
```

When prompted, input y to continue installation.

After the installation is complete, you can find VSP One Block Storage Provider for
HashiCorp Terraform in the following directory: /opt/hitachi/terraform/
examples

**Next steps**

Run the clean.sh script to clean up the Terraform directory. The cleanup is required to
remove, for example, Terraform generated state files.

See the clean.sh section for more information about clean.sh.

---

### Example Installation Output
```
Verifying...                          ################################# [100%]
Preparing...                          ################################# [100%]
[Tue Jun 17 11:13:25 EDT 2025] Starting pre-install checks
[Tue Jun 17 11:13:27 EDT 2025] Pre-install checks passed
Updating / installing...
   1:HV_Storage_Terraform-02.3.0-50     ################################# [100%]
[Tue Jun 17 11:13:28 EDT 2025] Starting installation of HV_Storage_Terraform
[Tue Jun 17 11:13:28 EDT 2025] WARN: Overwriting existing directories under /opt/hitachi/terraform
[Tue Jun 17 11:13:28 EDT 2025] Installing terraform plugin for user1
[Tue Jun 17 11:13:28 EDT 2025] Installation complete
[Tue Jun 17 11:13:28 EDT 2025] Installation successful

Hitachi Vantara LLC collects usage data such as storage model, storage serial number, operation name, status (success or failure),
and duration. This data is collected for product improvement purposes only. It remains confidential and it is not shared with any
third parties. To provide your consent, run bin/user_consent.sh from /opt/hitachi/terraform.
```

Log:
```
/var/log/hitachi/terraform/hitachi_terraform_install.log
```

Verify:
```bash
# Check the plugin terraform-provider-hitachi version (-v or --version)
cd ~/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/2.3.0/linux_amd64
./terraform-provider-hitachi -v
```

Example Output:
```text
Hitachi Terraform Provider version: 2.3
```

---

## Check Terraform Providers

Run inside your TF directory:

```bash
terraform init
terraform providers
```

Example output:
```text
provider[localhost/hitachi-vantara/hitachi] ~> 2.3
```

---

## Uninstall RPM Package

You can uninstall VSP One Block Storage Provider for HashiCorp Terraform, using the below
procedure.

Example files in the /opt/hitachi/terraform/examples/ path are
deleted during uninstall; back up files if required.

Procedure

1. Open the command line console.
2. Back up any Terraform configuration files.
3. Run the following command:

```bash
yum erase HV_Storage_Terraform.x86_64
```

When prompted, input y to continue uninstallation.

Uninstalling the Terraform provider will not delete logs from
the /var/log/hitachi/terraform directory.

Uninstalling the Terraform provider will not delete the
user_consent.json file from the /opt/hitachi/terraform/
directory.

---

### Example Uninstallation Output
```
[Tue Jun 17 11:13:14 EDT 2025] Starting uninstallation of HV_Storage_Terraform
[Tue Jun 17 11:13:14 EDT 2025] WARN: Deleting /opt/hitachi/terraform and contents
[Tue Jun 17 11:13:14 EDT 2025] Erasing terraform plugin 2.3.0 for user1
[Tue Jun 17 11:13:14 EDT 2025] Removing install directory /opt/hitachi/terraform
[Tue Jun 17 11:13:14 EDT 2025] Erase complete
[Tue Jun 17 11:13:14 EDT 2025] Uninstallation complete
```

Log:
```
/var/log/hitachi/terraform/hitachi_terraform_uninstall.log
```

---

## Directory Layout

Under `/opt/hitachi/terraform`:

| Path                          | Description                              |
|------------------------------|------------------------------------------|
| `bin/terraform-provider-hitachi` | Provider binary                    |
| `bin/user_consent.sh`        | Consent script                           |
| `bin/logbundle.sh`           | Consent script                           |
| `bin/.internal_config`       | For internal use (not for user)          |
| `docs/`                      | Documentation                            |
| `examples/`                  | Sample Terraform configs                 |
| `telemetry/*.json`           | Telemetry Usage data                     |
| `user_consent.json`          | User consent info                        |
| `logbundles/`                | Log bundle archives                        |

---

## Storage configuration operations

This section provides an example of how to perform the following tasks:

1. Create configuration changes on the storage system
2. Read a configuration from the storage system
3. Update a configuration on the storage system
4. Delete configuration from the storage system

The following examples to create, read, update, and delete are for a VSP One SDS Block
and Cloud system. Similar operations can be performed on san_storage_system
resources/data-sources.

After the installation of Terraform, go to the /opt/hitachi/terraform/examples/datasources
or /opt/hitachi/terraform/examples/resources directory.

Perform the following:
Verify that the provider file is configured. For example:

```hcl
terraform {
  required_providers {
    hitachi = {
      version = "2.3"
      source = "localhost/hitachi-vantara/hitachi"
    }
  }
}
provider "hitachi" {
  hitachi_vosb_provider {
    vosb_address = "192.0.2.100"
    username = var.hitachi_storage_user
    password = var.hitachi_storage_password
  }
}
```

Initialize the Terraform working directory:

```hcl
# terraform init
Initializing the backend...
Initializing provider plugins...
- Reusing previous version of localhost/hitachi-vantara/hitachi from the dependency
lock file
- Using previously-installed localhost/hitachi-vantara/hitachi v2.3

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
#
```

**1. Create configuration changes on the storage system**

Create a new volume on the storage system. Update the relevant details in the
hitachi_vosb_volume resource in the /opt/hitachi/terraform/examples/
resources directory. Select the resource file. In that resource file, this is the default
example. Update the default file:

```hcl
resource "hitachi_vosb_volume" "volumecreate" {
  vosb_address = "192.0.2.100"
  name = "test-volume-001"
  capacity_gb = 1
  storage_pool = "SP01"
  compute_nodes = []
  nick_name = "test-volume-001"
}
output "volumecreateData" {
  value = resource.hitachi_vosb_volume.volumecreate
}
```

Sample output:
```hcl
# terraform apply

Terraform used the selected providers to generate the following execution plan.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # hitachi_vosb_volume.volumecreate will be created
  + resource "hitachi_vosb_volume" "volumecreate" {
    + capacity_gb = 1
    + compute_nodes = []
    + id = (known after apply)
    + name = "test-volume-001"
    + nick_name = "test-volume-001"
    + storage_pool = "SP01"
    + volume = (known after apply)
    + vosb_address = "192.0.2.100"
  }
Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + volumecreateData = {
    + capacity_gb = 1
    + compute_nodes = []
    + id = (known after apply)
    + name = "test-volume-001"
    + nick_name = "test-volume-001"
    + storage_pool = "SP01"
    + volume = (known after apply)
    + vosb_address = "192.0.2.100"
  }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

hitachi_vosb_volume.volumecreate: Creating...
hitachi_vosb_volume.volumecreate: Creation complete after 9s [id=ace94e65-212b-4af8-
96fa-d8a4f5fec5cf]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

volumecreateData = {
  "capacity_gb" = 1
  "compute_nodes" = tolist([])
  "id" = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
  "name" = "test-volume-001"
  "nick_name" = "test-volume-001"
  "storage_pool" = "SP01"
  "volume" = tolist([
  {
    "compute_nodes" = tolist([])
    "data_reduction_progress_rate" = 0
    "data_reduction_status" = "Disabled"
    "full_allocated" = false
    "id" = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
    "name" = "test-volume-001"
    "nick_name" = "test-volume-001"
    "number_of_connecting_servers" = 0
    "number_of_snapshots" = 0
    "pool_id" = "7ab327f8-2a70-4a55-8b32-f14fd5c25a05"
    "pool_name" = "SP01"
    "protection_domain_id" = "67d99ad8-be21-4f77-916d-b54eb37fe1aa"
    "saving_effects" = tolist([
    {
      "post_capacity_data_reduction" = 0
      "pre_capacity_data_reduction_without_system_data" = 0
      "system_data_capacity" = 0
    },
    ])
    "saving_mode" = ""
    "saving_setting" = "Disabled"
    "snapshot_attribute" = "-"
    "snapshot_status" = ""
    "status" = "Normal"
    "status_summary" = "Normal"
    "storage_controller_id" = "b895bdae-ee54-4b5f-ab08-7134eb6b775f"
    "total_capacity" = 1024
    "used_capacity" = 0
    "volume_number" = 10043
    "volume_type" = "Normal"
  },
  ])
  "vosb_address" = "192.0.2.100"
}
```

**2. Read a configuration from the storage system**
Read the storage system configuration. Update the relevant details in the
hitachi_vosb_compute_nodes resource in the /opt/hitachi/terraform/
examples/data-sources directory. Select the data-source file. In that data-source file, this
is the default example. Update the default file:

```hcl
data "hitachi_vosb_compute_nodes" "computenodes" {
  vosb_address = "192.0.2.100"
  compute_node_name = "ComputeNode-15"
}
output "nodeoutput" {
  value = data.hitachi_vosb_compute_nodes.computenodes
}
```

Sample output:

```hcl
# terraform apply
data.hitachi_vosb_compute_nodes.computenodes: Reading...
data.hitachi_vosb_compute_nodes.computenodes: Read complete after 3s [id=b3ca158e-
9415-475f-8162-5cf86a2cc13d]

Changes to Outputs:
  + nodeoutput = {
    + compute_node_name = "ComputeNode-15"
    + compute_nodes = [
      + {
        + id = "b3ca158e-9415-475f-8162-5cf86a2cc13d"
        + nickname = "ComputeNode-15"
        + number_of_paths = 3
        + number_of_volumes = 1
        + os_type = "VMware"
        + paths = [
        + {
          + hba_name = "iqn.1993-08.org.debian.iscsi:01:109de7e4254k"
          + port_ids = [
            + "d92753e6-1c3e-4716-b560-a9872692cae9",
            + "6d0e8179-ce50-4752-b18d-792cf12bda80",
            + "66c334ac-d23b-4469-9a14-aed94d0a07f5",
          ]
        + protocol = "iSCSI"
      },
    ]
    + port_details = [
      + {
        + iscsi_initiator = "iqn.1994-04.jp.co.hitachi:rsd.sph.t.07072.000"
        + port_id = "d92753e6-1c3e-4716-b560-a9872692cae9"
        + port_name = "001-iSCSI-000"
      },
    + {
        + iscsi_initiator = "iqn.1994-04.jp.co.hitachi:rsd.sph.t.07072.001"
        + port_id = "6d0e8179-ce50-4752-b18d-792cf12bda80"
        + port_name = "000-iSCSI-001"
      },
    + {
        + iscsi_initiator = "iqn.1994-04.jp.co.hitachi:rsd.sph.t.07072.004"
        + port_id = "66c334ac-d23b-4469-9a14-aed94d0a07f5"
        + port_name = "002-iSCSI-004"
      },
    ]
    + total_capacity = 1945
    + used_capacity = 0
  },
  ]
  + id = "b3ca158e-9415-475f-8162-5cf86a2cc13d"
  + vosb_address = "192.0.2.100"
  }

You can apply this plan to save these new output values to the Terraform state,
without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

nodeoutput = {
  "compute_node_name" = "ComputeNode-15"
  "compute_nodes" = tolist([
  {
    "id" = "b3ca158e-9415-475f-8162-5cf86a2cc13d"
    "nickname" = "ComputeNode-15"
    "number_of_paths" = 3
    "number_of_volumes" = 1
    "os_type" = "VMware"
    "paths" = tolist([
    {
      "hba_name" = "iqn.1993-08.org.debian.iscsi:01:109de7e4254k"
      "port_ids" = tolist([
        "d92753e6-1c3e-4716-b560-a9872692cae9",
        "6d0e8179-ce50-4752-b18d-792cf12bda80",
        "66c334ac-d23b-4469-9a14-aed94d0a07f5",
    ])
      "protocol" = "iSCSI"
    },
    ])
    "port_details" = tolist([
    {
      "iscsi_initiator" = "iqn.1994-04.jp.co.hitachi:rsd.sph.t.07072.000"
      "port_id" = "d92753e6-1c3e-4716-b560-a9872692cae9"
      "port_name" = "001-iSCSI-000"
    },
    {
      "iscsi_initiator" = "iqn.1994-04.jp.co.hitachi:rsd.sph.t.07072.001"
      "port_id" = "6d0e8179-ce50-4752-b18d-792cf12bda80"
      "port_name" = "000-iSCSI-001"
    },
    {
      "iscsi_initiator" = "iqn.1994-04.jp.co.hitachi:rsd.sph.t.07072.004"
      "port_id" = "66c334ac-d23b-4469-9a14-aed94d0a07f5"
      "port_name" = "002-iSCSI-004"
    },
    ])
    "total_capacity" = 1945
    "used_capacity" = 0
    },
  ])
  "id" = "b3ca158e-9415-475f-8162-5cf86a2cc13d"
  "vosb_address" = "192.0.2.100"
}
#
```

**3. Update a configuration on the storage system**
Update the storage system configuration. For example, you can change the volume name.
Enter the following configuration:

```hcl
resource "hitachi_vosb_volume" "volumecreate" {
  vosb_address = "192.0.2.100"
  name = "NEW-test-volume-001"
  capacity_gb = 1
  storage_pool = "SP01"
  compute_nodes = []
  nick_name = "test-volume-001"
}
output "volumecreateData" {
  value = resource.hitachi_vosb_volume.volumecreate
}
```

Sample output:

```hcl
# terraform apply
hitachi_vosb_volume.volumecreate: Refreshing state... [id=ace94e65-212b-4af8-96fad8a4f5fec5cf]

Terraform used the selected providers to generate the following execution plan.
Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # hitachi_vosb_volume.volumecreate will be updated in-place
  ~ resource "hitachi_vosb_volume" "volumecreate" {
    id = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
    ~ name = "test-volume-001" -> "NEW-test-volume-001"
    ~ volume = [
    - {
      - compute_nodes = []
      - data_reduction_progress_rate = 0
      - data_reduction_status = "Disabled"
      - full_allocated = false
      - id = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
      - name = "test-volume-001"
      - nick_name = "test-volume-001"
      - number_of_connecting_servers = 0
      - number_of_snapshots = 0
      - pool_id = "7ab327f8-2a70-4a55-8b32-f14fd5c25a05"
      - pool_name = "SP01"
      - protection_domain_id = "67d99ad8-be21-4f77-916d-b54eb37fe1aa"
      - saving_effects = [
      - {
          - post_capacity_data_reduction = 0
          - pre_capacity_data_reduction_without_system_data = 0
          - system_data_capacity = 0
        },
        ]
        - saving_setting = "Disabled"
        - snapshot_attribute = "-"
        - status = "Normal"
        - status_summary = "Normal"
        - storage_controller_id = "b895bdae-ee54-4b5f-ab08-7134eb6b775f"
        - total_capacity = 1024
        - used_capacity = 0
        - volume_number = 10043
        - volume_type = "Normal"
        # (2 unchanged attributes hidden)
        },
      ] -> (known after apply)
      # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ volumecreateData = {
    id = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
  ~ name = "test-volume-001" -> "NEW-test-volume-001"
  ~ volume = [
    - {
        - compute_nodes = []
        - data_reduction_progress_rate = 0
        - data_reduction_status = "Disabled"
        - full_allocated = false
        - id = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
        - name = "test-volume-001"
        - nick_name = "test-volume-001"
        - number_of_connecting_servers = 0
        - number_of_snapshots = 0
        - pool_id = "7ab327f8-2a70-4a55-8b32-f14fd5c25a05"
        - pool_name = "SP01"
        - protection_domain_id = "67d99ad8-be21-4f77-916d-b54eb37fe1aa"
        - saving_effects = [
        - {
          - post_capacity_data_reduction = 0
          - pre_capacity_data_reduction_without_system_data = 0
          - system_data_capacity = 0
          },
        ]
        - saving_mode = ""
        - saving_setting = "Disabled"
        - snapshot_attribute = "-"
        - snapshot_status = ""
        - status = "Normal"
        - status_summary = "Normal"
        - storage_controller_id = "b895bdae-ee54-4b5f-ab08-7134eb6b775f"
        - total_capacity = 1024
        - used_capacity = 0
        - volume_number = 10043
        - volume_type = "Normal"
      },
    ] -> (known after apply)
    # (5 unchanged attributes hidden)
  }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

hitachi_vosb_volume.volumecreate: Modifying... [id=ace94e65-212b-4af8-96fad8a4f5fec5cf]
hitachi_vosb_volume.volumecreate: Modifications complete after 6s [id=ace94e65-212b-
4af8-96fa-d8a4f5fec5cf]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

volumecreateData = {
  "capacity_gb" = 1
  "compute_nodes" = tolist([])
  "id" = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
  "name" = "NEW-test-volume-001"
  "nick_name" = "test-volume-001"
  "storage_pool" = "SP01"
  "volume" = tolist([
    {
      "compute_nodes" = tolist([])
      "data_reduction_progress_rate" = 0
      "data_reduction_status" = "Disabled"
      "full_allocated" = false
      "id" = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
      "name" = "NEW-test-volume-001"
      "nick_name" = "test-volume-001"
      "number_of_connecting_servers" = 0
      "number_of_snapshots" = 0
      "pool_id" = "7ab327f8-2a70-4a55-8b32-f14fd5c25a05"
      "pool_name" = "SP01"
      "protection_domain_id" = "67d99ad8-be21-4f77-916d-b54eb37fe1aa"
      "saving_effects" = tolist([
      {
        "post_capacity_data_reduction" = 0
        "pre_capacity_data_reduction_without_system_data" = 0
        "system_data_capacity" = 0
      },
      ])
      "saving_mode" = ""
      "saving_setting" = "Disabled"
      "snapshot_attribute" = "-"
      "snapshot_status" = ""
      "status" = "Normal"
      "status_summary" = "Normal"
      "storage_controller_id" = "b895bdae-ee54-4b5f-ab08-7134eb6b775f"
      "total_capacity" = 1024
      "used_capacity" = 0
      "volume_number" = 10043
      "volume_type" = "Normal"
    },
  ])
  "vosb_address" = "192.0.2.100"
}
#
```

**4. Delete configuration from the storage system**
Sample output:

```hcl
# terraform destroy
hitachi_vosb_volume.volumecreate: Refreshing state... [id=ace94e65-212b-4af8-96fad8a4f5fec5cf]

Terraform used the selected providers to generate the following execution plan.
Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # hitachi_vosb_volume.volumecreate will be destroyed
  - resource "hitachi_vosb_volume" "volumecreate" {
    - capacity_gb = 1 -> null
    - compute_nodes = [] -> null
    - id = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf" -> null
    - name = "NEW-test-volume-001" -> null
    - nick_name = "test-volume-001" -> null
    - storage_pool = "SP01" -> null
    - volume = [
      - {
        - compute_nodes = []
        - data_reduction_progress_rate = 0
        - data_reduction_status = "Disabled"
        - full_allocated = false
        - id = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
        - name = "NEW-test-volume-001"
        - nick_name = "test-volume-001"
        - number_of_connecting_servers = 0
        - number_of_snapshots = 0
        - pool_id = "7ab327f8-2a70-4a55-8b32-f14fd5c25a05"
        - pool_name = "SP01"
        - protection_domain_id = "67d99ad8-be21-4f77-916d-b54eb37fe1aa"
        - saving_effects = [
          - {
              - post_capacity_data_reduction = 0
              - pre_capacity_data_reduction_without_system_data = 0
              - system_data_capacity = 0
            },
          ]
        - saving_setting = "Disabled"
        - snapshot_attribute = "-"
        - status = "Normal"
        - status_summary = "Normal"
        - storage_controller_id = "b895bdae-ee54-4b5f-ab08-7134eb6b775f"
        - total_capacity = 1024
        - used_capacity = 0
        - volume_number = 10043
        - volume_type = "Normal"
          # (2 unchanged attributes hidden)
        },
      ] -> null
    - vosb_address = "192.0.2.100" -> null
}

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - volumecreateData = {
    - capacity_gb = 1
    - compute_nodes = []
    - id = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
    - name = "NEW-test-volume-001"
    - nick_name = "test-volume-001"
    - storage_pool = "SP01"
    - volume = [
      - {
        - compute_nodes = []
        - data_reduction_progress_rate = 0
        - data_reduction_status = "Disabled"
        - full_allocated = false
        - id = "ace94e65-212b-4af8-96fa-d8a4f5fec5cf"
        - name = "NEW-test-volume-001"
        - nick_name = "test-volume-001"
        - number_of_connecting_servers = 0
        - number_of_snapshots = 0
        - pool_id = "7ab327f8-2a70-4a55-8b32-f14fd5c25a05"
        - pool_name = "SP01"
        - protection_domain_id = "67d99ad8-be21-4f77-916d-b54eb37fe1aa"
        - saving_effects = [
          - {
              - post_capacity_data_reduction = 0
              - pre_capacity_data_reduction_without_system_data = 0
              - system_data_capacity = 0
            },
          ]
        - saving_mode = ""
        - saving_setting = "Disabled"
        - snapshot_attribute = "-"
        - snapshot_status = ""
        - status = "Normal"
        - status_summary = "Normal"
        - storage_controller_id = "b895bdae-ee54-4b5f-ab08-7134eb6b775f"
        - total_capacity = 1024
        - used_capacity = 0
        - volume_number = 10043
        - volume_type = "Normal"
        },
      ]
    - vosb_address = "192.0.2.100"
  } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

hitachi_vosb_volume.volumecreate: Destroying... [id=ace94e65-212b-4af8-96fad8a4f5fec5cf]
hitachi_vosb_volume.volumecreate: Destruction complete after 9s

Destroy complete! Resources: 1 destroyed.
#
```

---

## Logging

| File Path                                                  | Purpose                                  |
|------------------------------------------------------------|------------------------------------------|
| `/var/log/hitachi/terraform/hitachi_terraform_install.log` | Installation log                         |
| `/var/log/hitachi/terraform/hitachi_terraform_uninstall.log` | Uninstallation log                      |
| `/var/log/hitachi/terraform/hitachi-terraform.log`         | Runtime log                              |

---

## Log Bundle Collection
The `logbundle.sh` script collects a support log bundle containing Terraform configuration files, state files, crash logs, plugin logs, and system/environment diagnostics. It is designed to assist with support, debugging, and issue investigation.

### Usage

```
/opt/hitachi/terraform/bin/logbundle.sh [tf_dir1 tf_dir2 ...]
```

+ Provide one or more directories to scan for .tf, .tfvars, and related Terraform files.

+ If no directories are provided, the script will prompt for input interactively.

  + Press Enter to use the default directories:

      + Default TF dirs: "." "/opt/hitachi/terraform/examples"

### Environment Variables
- `TF_MAX_LOGBUNDLES`.  
  Controls how many log bundles to retain. Older bundles beyond this limit will be deleted.  
  **Default: `3`

### Notes
+ If you're already inside a Terraform project directory and run the script without any arguments, only the default directories will be scanned.

+ If your project spans multiple folders, be sure to specify all relevant directories as arguments.

### Example

```
/opt/hitachi/terraform/bin/logbundle.sh .
```

### Sample Output
```
Max log bundles: 3 (can be set via environment variable TF_MAX_LOGBUNDLES)
📦 Collecting Terraform version info...
🔍 Searching Terraform directories under: .../examples/data-sources/hitachi_vosb_storage_pools
📦 Copying hitachi terraform plugin logs...
📦 Copying hitachi terraform plugin config and telemetry files...
📦 Collecting hitachi terraform plugin version...
📦 Collecting machine info...
📦 Copying logbundle script output log...
📦 Creating logbundle archive at /opt/hitachi/terraform/logbundles/hitachi_terraform_logbundle-20250616_162047.tar.gz...
✅ Log bundle created: /opt/hitachi/terraform/logbundles/hitachi_terraform_logbundle-20250616_162047.tar.gz
🧹 Cleaning up old logbundles, keeping only the last 3...
```

---

## Known Behavior

When unassigning fields in a list denoted by square brackets (ex: `[]`), **do not comment out** the parameter as this will cause Terraform to skip over the configuration parameter.  
Instead, pass it an empty string to remove all objects in the list.

As an example, we have a `port_names` field defined with two values:

```hcl
["002-iSCSI-001", "001-iSCSI-002"]
```

### Example: Assigning ports

```hcl
resource "hitachi_vss_block_compute_node" "mycompute" {
  vss_block_address  = "192.0.2.100"
  compute_node_name  = "ComputeNode"
  os_type            = "VMware"

  iscsi_connection {
    iscsi_initiator = "iqn.1998-01.com.vmware:node-06-0723aa94"
    port_names      = ["002-iSCSI-001", "001-iSCSI-002"]
  }
}
```

### Example: Unassigning ports

To unassign the two ports from the compute node, pass the following configuration:

```hcl
resource "hitachi_vss_block_compute_node" "mycompute" {
  vss_block_address  = "192.0.2.100"
  compute_node_name  = "ComputeNode"
  os_type            = "VMware"

  iscsi_connection {
    iscsi_initiator = "iqn.1998-01.com.vmware:node-06-0723aa94"
    port_names      = []
  }
}
```

---

## Usage data collection

Hitachi Vantara LLC collects usage data such as storage model, storage serial number,
operation name, status (success or failure), and duration. This data is collected for product
improvement purposes only. It remains confidential and it is not shared with any third parties.

user_consent.sh file path:

```hcl
/opt/hitachi/terraform/bin/user_consent.sh

```
After updating user consent, a record is saved at this location:
File path:

```hcl
/opt/hitachi/terraform/user_consent.json
```

The user_consent.sh script can be run again to disable/enable user consent.

Running Terraform data sources and resources results in successes or failures, which are
automatically logged to Hitachi Vantara module-specific counters in a usage.json file.

For example, after running the get lun facts task, a hv_hg_facts.get_luns success
counter is incremented.

Sample usage counters are shown below:

```hcl
{
"name": "san.datasource.hitachi_vsp_volume.GetLun","metrics": {
     "averageTimeInSec": 0.53,
     "success": 2,
     "failure": 0
}
```

The averageTimeInSec counter tracks the average call duration in seconds.


The usage.json file is available at: /opt/hitachi/terraform/telemetry
The failure counter is incremented if the API call fails while Terraform is running.
The usage.json and user_consent.json files are collected when the Terraform log bundle is generated.

### Requirements for the Terraform client to support Telemetry - Usage data collection

**Unrestricted Outgoing Traffic:**
Ensure that the client's firewall or security software allows outgoing HTTPS traffic on
port 443.

**Proxy Settings:**
- If the client is behind a proxy, verify that the proxy allows the CONNECT method on
port 443 for HTTPS connections.
- Configure proxy settings in the client application if needed.

**Trusted Certificates:**
Ensure the client's certificate store trusts the Certificate Authority (CA) that issued the
server's SSL/TLS certificate. This is crucial for establishing a secure connection.

**TLS/SSL Compatibility:**
Confirm that the client supports the required TLS versions (e.g., TLS 1.2 or 1.3) used
by the server.

**DNS Resolution:**
Make sure the client can resolve the API's domain name correctly to establish a
connection over port 443.

**Sample usage data collected**

```hcl
  "376b4e32-29f5-4e49-a2e5-58e0b7805b10": {
  "site": "376b4e32-29f5-4e49-a2e5-58e0b7805b10",
  "createDate": "2025-06-04T20:56:41.337660Z",
  "lastUpdate": "2025-06-04T20:56:41.337660Z",
  "sds_block": [
    {
     "model": "01.17.00.40",
     "serial": "",
     "vosb": [
       {
         "vosb.terraform.providerConfigure.GetStorageVersionInfo": {
           "success": 2,
           "failure": 0,
           "averageTimeInSec": 0.01
         }
       },
       {
         "vosb.resource.hitachi_vosb_add_drives_to_pool.GetDrivesInfo": {
           "success": 0,
           "failure": 1,
          "averageTimeInSec": 0.43
         }
       },
       {
         "vosb.datasource.hitachi_vosb_storage_drives.GetDrivesInfo": {
            "success": 2,
            "failure": 0,
            "averageTimeInSec": 0.46
         }
       }
     ]
   }
  ],
  "vsp": [
    {
       "model": "VSP 5600H",
       "serial": "40015",
  "san": [
    {
       "san.terraform.providerConfigure.GetStorageSystemInfo": {
          "success": 1,
          "failure": 0,
          "averageTimeInSec": 1.56
        }
    },
    {
       "san.datasource.hitachi_vsp_storage.GetStorageCapacity": {
              "success": 2,
              "failure": 0,
              "averageTimeInSec": 15.06
            }
          }
        ]
      }
    ]
  }
}
```
---

## Bundled TF Examples

```bash
cd /opt/hitachi/terraform/examples
cd data-sources/hitachi_vsp_storage

export TF_LOG=DEBUG
export TF_LOG_PATH="terraform.log"

# Clean up any existing state
./clean.sh   # or alternatively:
rm -rf .terraform .terraform.lock.hcl terraform.tfstate*

# Modify your .tf files as needed
terraform init
terraform apply
```

---

## User Consent Script

### Usage

```bash
cd /opt/hitachi/terraform
./bin/user_consent.sh
```

### Output

```
# bin/user_consent.sh

==================== USER CONSENT ====================

  Hitachi Vantara LLC collects usage data such as storage model, storage serial number, operation name, status (success or failure),
  and duration. This data is collected for product improvement purposes only. It remains confidential and it is not shared with any
  third parties.

======================================================

Do you consent to the collection of usage data? (Yes/No): Yes

✅ User consent has been recorded successfully.
```


### Example user consent:
Saves to:
```
/opt/hitachi/terraform/user_consent.json
```
```json
{
  "site_id": "...",
  "user_consent_accepted": true,
  "time": "...",
  "consent_history": [...]
}
```

---

## Supported host modes
The following table lists the supported host modes:

| Host mode | Type | Description | Required | Valid input values |
|-----------|------|-------------|----------|--------------------|
| host_mode | string | Host mode of the host group | No | Host mode string value. |
| | | | | The default value is "LINUX". Valid input values: |
| | | | | - LINUX |
| | | | | - VMWARE |
| | | | | - HP |
| | | | | - OPEN_VMS |
| | | | | - NETWARE |
| | | | | - WINDOWS |
| | | | | - HI_UX |
| | | | | - AIX |
| | | | | - VMWARE_EXTENSION |
| | | | | - WINDOWS_EXTENSION |
| | | | | - UVM |
| | | | | - HP_XP |
| | | | | - DYNIX |

---

## Host mode options

For host mode options, see:

https://docs.hitachivantara.com/r/en-us/svos/9.8.7/mk-97hm85026/managing-logical-volumes/configuring-hosts/host-modes-and-host-mode-options-for-host-facing-host-ports


---

## The clean shell script

The clean.sh script inside each example directory does the following:

**Clean Terraform Working Directory**
Removes all Terraform-generated state and configuration artifacts to reset the
environment. This includes the .terraform directory, lock files, and state files. Log
files are preserved.

**Removes**
.terraform/, .terraform.lock.hcl, terraform.tfstate,
terraform.tfstate.backup

**Conditions**
This script should only be used when you are:

- Fixing corrupted Terraform state or provider errors.
- Switching providers, backends, or environments.
- Starting fresh or troubleshooting persistent issues.
- After each re-installation of Terraform provider.

**To run, for example**

```bash
cd /opt/hitachi/terraform/examples/data-sources/hitachi_vsp_parity_group

./clean.sh

# perform your terraform commands.
```

**clean.sh script code**

```hcl
#!/bin/sh

echo "Start cleaning terraform files"
rm -rf .terraform
rm -rf .terraform.lock.hcl
rm -rf terraform.tfstate*
echo "Done"
```

---

## Advanced workflows

**Steps to perform the Add storage node task**

1. Create the configuration file.
Refer to the following section to export the configuration file:
hitachi_vosb_configuration_file (Resource)

2. Run the Add storage node module.
Refer to the following section, use the exported configuration file, and execute the
register storage node: hitachi_vosb_storage_node (Resource)

---

## Terraform Compatibility

Tested with:
- Terraform >= 1.14.4

---

## Linux Compatibility

Tested on:
- Oracle Enterprise Linux 8.10

---

## Planning for Mainframe Volume capacity

For mainframe only:

1. The formula for calculation of blocks

1 cylinder = 1740 blocks

raw_blocks = requested_cylinders * BLOCKS_PER_CYL

2. Calculate the number of boundary units (rounding up)

This ensures the result is divisible by the page unit (77,952)

boundary_units = math.ceil(raw_blocks / 77952)

3. Calculate the final capacity in blocks

final_blocks = boundary_units * 1740

---

## Privacy Notice

For details, see [Hitachi Vantara’s Privacy Policy](https://www.hitachivantara.com/en-us/company/legal/privacy.html).

---

## License

Proprietary — Hitachi Vantara LLC.
