# Terraform Provider Hitachi

## Table of Contents
1. [Introduction](#introduction)
2. [Prerequisites & Dependencies](#prerequisites--dependencies)
3. [Build the Provider](#build-the-provider)
   - [3.1 Build RPM package](#build-rpm-package)
   - [3.2 Build without RPM](#for-developers-without-rpm)
4. [Install RPM Package](#install-rpm-package)
5. [Check Terraform Providers](#check-terraform-providers)
6. [Uninstall RPM Package](#uninstall-rpm-package)
7. [Directory Layout](#directory-layout)
8. [Logging](#logging)
9. [Log Bundle Collection](#log-bundle-collection)
10. [Bundled TF Examples](#bundled-tf-examples)
11. [User Consent Script](#user-consent-script)
12. [Terraform Compatibility](#terraform-compatibility)
13. [Linux Compatibility](#linux-compatibility)
14. [Privacy Notice](#privacy-notice)
15. [License](#license)

---

## Introduction

Terraform provider for Hitachi Vantara storage systems 2.1.


---

## Prerequisites & Dependencies

Ensure the following are installed:
- `jq`
- `terraform` (>= 1.11.4)
- `uuidgen`
- Standard runtime libraries (`glibc`, etc.)

---

## Build the Provider

### 3.1 Build RPM package (for Developers)

Requires **Golang v1.22** and **superuser privileges**.

```bash
export GOPATH=/usr/local/go
cd <your hitachi terraform source code directory>
./build.sh <BUILD_NUMBER>
```

Example output path:
```
./rpmbuild/RPMS/x86_64/HV_Storage_Terraform-02.1-40.x86_64.rpm
```

### 3.2 Build without RPM (for Developers)

```bash
make all
cd examples
terraform init && terraform apply
```

---

## Install RPM Package

**Warning:** Installation will overwrite `/opt/hitachi/terraform/`. Back up important data before continuing.

### Steps:
```bash
# Check if the package is already installed:
/usr/bin/rpm -qa HV_Storage_Terraform

# Uninstall the old version:
/usr/bin/rpm -e HV_Storage_Terraform

# Install the new version: replace build number (XX)
/usr/bin/rpm -Uvh HV_Storage_Terraform-02.1-XX.x86_64.rpm
```

### Example Installation Output
```
Verifying...                          ################################# [100%]
Preparing...                          ################################# [100%]
[Tue Jun 17 11:13:25 EDT 2025] Starting pre-install checks
[Tue Jun 17 11:13:27 EDT 2025] Pre-install checks passed
Updating / installing...
   1:HV_Storage_Terraform-02.1-50     ################################# [100%]
[Tue Jun 17 11:13:28 EDT 2025] Starting installation of HV_Storage_Terraform
[Tue Jun 17 11:13:28 EDT 2025] WARN: Overwriting existing directories under /opt/hitachi/terraform
[Tue Jun 17 11:13:28 EDT 2025] Installing terraform plugin for root
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
cd ~/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/2.1/linux_amd64
./terraform-provider-hitachi -v
```

Example Output:
```text
Hitachi Terraform Provider version: 2.1
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
provider[localhost/hitachi-vantara/hitachi] ~> 2.1
```

---

## Uninstall RPM Package

**Warning:** Will delete dirs and files (except user_consent.json) under `/opt/hitachi/terraform/`. Back up important data before continuing.


```bash
/usr/bin/rpm -e HV_Storage_Terraform
```

### Example Uninstallation Output
```
[Tue Jun 17 11:13:14 EDT 2025] Starting uninstallation of HV_Storage_Terraform
[Tue Jun 17 11:13:14 EDT 2025] WARN: Deleting /opt/hitachi/terraform and contents
[Tue Jun 17 11:13:14 EDT 2025] Erasing terraform plugin 2.1 for root
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

## Terraform Compatibility

Tested with:
- Terraform >= 1.11.4 and <= 1.12.2

---

## Linux Compatibility

Tested on:
- Oracle Enterprise Linux 8.10

---

## Privacy Notice

For details, see [Hitachi Vantara’s Privacy Policy](https://www.hitachivantara.com/en-us/company/legal/privacy.html).

---

## License

Proprietary — Hitachi Vantara LLC.
