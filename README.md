### Terraform Provider Hitachi


## Build RPM package:
You must be a superuser to build the package.

You must use Golang version 1.22 and set GOPATH to your Golang v1.22 directory.
```script
    export GOPATH=/usr/local/go
    cd <your hitachi terraform source code directory>
    ./build.sh $BUILD_NUMBER
```

The RPM package is in 
```
    ./rpmbuild/RPMS/x86_64/HV_Storage_Terraform-02.1-1.x86_64.rpm
```


## Install RPM package
You must be a superuser to install the RPM package.

**Warning:**
This installation will **overwrite** the contents of `/opt/hitachi/terraform/` (including the directories: `bin/`, `docs/`, `examples/`). 
Ensure you have backups of any important data before proceeding with the installation.

If you proceed with the installation, these directories will be replaced by the new version of the provider. 

**Installation steps:**
1. **Check if the package is already installed:**
    ```script
    /usr/bin/rpm -qa HV_Storage_Terraform
    ```


2. **Uninstall the old version (if necessary):**
    ```script
    See below section 'Uninstall RPM package'
    ```


3. **Install the new version:**
    ```script
    /usr/bin/rpm -Uvh HV_Storage_Terraform-02.1-1.x86_64.rpm
    ```


4. **Check if the Hitachi Terraform plugin is installed properly:**
    ```
    After installation, the Hitachi provider plugin should be linked to the following path:
    ```
    ```
    ls -l /root/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/2.1/linux_amd64/terraform-provider-hitachi 

    lrwxrwxrwx. 1 root root 80 Aug  9 21:55 /root/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/2.1/linux_amd64/terraform-provider-hitachi -> /opt/hitachi/terraform/bin/terraform-provider-hitachi
    ```

A detailed installation log is written to:

```
/var/log/hitachi_terraform_install.log
```

## Uninstall RPM package

**Warning:**
If you need to uninstall the package, be aware that it will **delete the entire directory** of `/opt/hitachi/terraform/` and its subdirectories (`bin/`, `docs/`, `examples/`). 
Ensure you have backups of any important data before proceeding with the uninstallation.

**Uninstallation steps:**
1. **Uninstall the package:**
    ```script
    /usr/bin/rpm -e HV_Storage_Terraform
    ```

2. **Verify that the plugin has been removed:**
    ```script
    ls /root/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/2.1/linux_amd64/terraform-provider-hitachi
    ```

A detailed uninstallation log is written to:

```
/var/log/hitachi_terraform_uninstall.log
```

## Use the tf samples
1. **Navigate to /opt/hitachi/terraform/examples:**
    ```
    cd /opt/hitachi/terraform/examples
    ```

2. **Go to any examples directory. Example:**
    ```
    cd data-sources/hitachi_vsp_storage
    ```

3. **If this is **not the first time** using the examples directory, you may want to clean up previous configurations:**
    ```
    rm .terraform .terraform.lock.hcl terraform.tfstate
    ```

4. **Modify `provider.tf` with your VSP or VOSB storage information, then run:**
    ```
    terraform init
    terraform apply
    ```


## For developers only without the RPM package

If you're not using the RPM package and want to build and install the provider manually, follow these steps:

1. **Build and install the provider:**
   ```shell
   make all
    ```
2. **Then, navigate to the `examples` directory.**
    ```shell
    cd examples
    ```

3. **Run the following command to initialize the workspace and apply the sample configuration.**
    ```shell
    terraform init && terraform apply
    ```
