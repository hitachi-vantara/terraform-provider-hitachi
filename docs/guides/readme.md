# Terraform Provider Hitachi

## Build RPM package
You must be a superuser to build the package.

In build.sh, modify this variable to point to your hitachi terraform source code directory

    TERRAFORM_SRC_DIR

You must use Golang version 1.20, and set GOPATH to your Golang v1.20 directory.
```
# export GOPATH=/usr/local/go
# ./build
```

The RPM package is in /root/rpmbuild/RPMS/x86_64/HV_Storage_Terraform-02.0.0-1.el7.x86_64.rpm


## Install RPM package
You must be a superuser to install the rpm package.

Check what RPM version is installed

```
# /usr/bin/rpm -qa HV_Storage_Terraform
```

Uninstall the old version

```
# /usr/bin/rpm -e HV_Storage_Terraform
```

Install the new version

```
# /usr/bin/rpm -Uvh HV_Storage_Terraform-02.0.0-1.el7.x86_64.rpm
```

Check if the hitachi terraform plugin is installed. It must be linked to /opt/hitachi-vantara/storage-systems/terraform-provider/bin/terraform-provider-hitachi

```
# ls -l /root/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/2.0/linux_amd64/terraform-provider-hitachi 
lrwxrwxrwx 1 root root 61 Mar 22 19:58 /root/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/2.0/linux_amd64/terraform-provider-hitachi -> /opt/hitachi-vantara/storage-systems/terraform-provider/bin/terraform-provider-hitachi
```

## Use the tf samples
Navigate to /opt/hitachi-vantara/storage-systems/terraform-provider/samples

```
# cd /opt/hitachi-vantara/storage-systems/terraform-provider/samples
```

Go to any samples directory. Example:

```
# cd data/hitachi_storage_data
```

If not the first time using the sample directory, do cleanup

```
# rm .terraform .terraform.lock.hcl terraform.tfstate
# rm -rf san_settings
```

Modify provider.tf and storage.tf with your storage information, then do:

```
# terraform init
# terraform apply
```


## For developers only without the RPM package

First, build and install the provider.

```shell
$ make all
```

Then, navigate to the `examples` directory. 

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```