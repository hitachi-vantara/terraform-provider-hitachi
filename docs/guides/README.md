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
9. [Bundled Examples](#bundled-examples)
11. [User Consent Script](#user-consent-script)
12. [Config](#config)
13. [Terraform Compatibility](#terraform-compatibility)
14. [Linux Compatibility](#linux-compatibility)
15. [Privacy Notice](#privacy-notice)
16. [License](#license)

---

## Introduction

Terraform provider for Hitachi Vantara storage systems 2.1.

👉 Download latest RPM

---

## Prerequisites & Dependencies

Ensure the following are installed:
- `jq`
- `terraform` (>= 1.11.4, < 2.0.0)
- `uuidgen`
- Standard runtime libraries (`glibc`, etc.)

---

## Build the Provider

### 3.1 Build RPM package (for Developers)

Requires Golang v1.22 and superuser privileges.

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

# Install the new version: replace build number (40)
/usr/bin/rpm -Uvh HV_Storage_Terraform-02.1-40.x86_64.rpm
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

Logs:
```
/var/log/hitachi/terraform/hitachi_terraform_install.log
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

**Warning:** Will delete `/opt/hitachi/terraform/`. Back up important data before continuing.


```bash
/usr/bin/rpm -e HV_Storage_Terraform
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
| `docs/`                      | Documentation                            |
| `examples/`                  | Sample Terraform configs                 |
| `telemetry/*.json`           | Telemetry Usage data                     |
| `config.json`                | Configuration file                       |
| `user_consent.json`          | User consent info                        |

---

## Logging

| File Path                                                  | Purpose                                  |
|------------------------------------------------------------|------------------------------------------|
| `/var/log/hitachi/terraform/hitachi_terraform_install.log` | Installation log                         |
| `/var/log/hitachi/terraform/hitachi_terraform_uninstall.log` | Uninstallation log                      |
| `/var/log/hitachi/terraform/hitachi-terraform.log`         | Runtime log                              |

---

## Bundled Examples

```bash
cd /opt/hitachi/terraform/examples
cd data-sources/hitachi_vsp_storage
./clean.sh

# modify your .tf files
terraform init
terraform apply
```

---

## User Consent Script

```bash
cd /opt/hitachi/terraform
./bin/user_consent.sh
```

Saves to:
```
/opt/hitachi/terraform/user_consent.json
```

Example content:
```json
{
  "site_id": "...",
  "user_consent_accepted": true,
  "time": "...",
  "consent_history": [...]
}
```

---

## Config

Located at: `/opt/hitachi/terraform/config.json`

### Fields

| Field                  | Description |
|------------------------|-------------|
| `user_consent_message` | Message shown to the user. **Editable only by QA, developers, or other internal teams.** |
| `api_timeout`          | Timeout (seconds) for internal API operations. |
| `aws_timeout`          | Timeout (seconds) for AWS-related calls. |
| `aws_url`              | Optional URL for AWS integrations for Telemetry (can be blank). |

### Example

```json
{
  "user_consent_message": "...",
  "api_timeout": 300,
  "aws_timeout": 300,
  "aws_url": ""
}
```

---

## Terraform Compatibility

Tested with:
- \>= 1.11.4
- < 2.0.0

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
