# -*- rpm-spec -*-
# Spec file for HV_Storage_Terraform

Name:          HV_Storage_Terraform
Version:       %{_VERSION}
Release:       %{_BUILD_NUMBER}
Summary:       Hitachi Vantara storage plugin provider for Terraform
Vendor:        Hitachi Vantara
Group:         Adapters
License:       hiAdapterLicense
URL:           http://www.hitachivantara.com
Source0:       HV_Storage_Terraform-%{_VERSION}.tar.gz
ExclusiveArch: x86_64
BuildRoot:     %{_tmppath}/%{name}-%{_VERSION}.%{_release}-root-%(%{__id_u} -n)
AutoReqProv:   no

%description
Hitachi Terraform RPM Package

%define _missing_build_ids_terminate_build 0
%define debug_package %{nil}

%prep
%setup -q

%build
%define hitachi_base /opt/hitachi
%define terraform %{hitachi_base}/terraform
%define plugin_dir  .terraform.d/plugins
%define terraform_plugin %{plugin_dir}/localhost/hitachi-vantara/hitachi/%{_DISPLAY_VERSION}/linux_amd64

%define source_dir %{_topdir}/HV_Storage_Terraform-%{_VERSION}
%define examples_src %{source_dir}/examples
%define docs_src %{source_dir}/docs

%define examples_dst %{buildroot}/%{terraform}/examples
%define docs_dst %{buildroot}/%{terraform}/docs

%install
rm -rf $RPM_BUILD_ROOT

install -d %{buildroot}/%{terraform}/bin

# Create all the directory inside SRC examples path
for dir in %{examples_src}/*; do
  if [ -d "$dir" ]; then
    main_dir="examples"/${dir##*/}
    install -d %{buildroot}/%{terraform}/"$main_dir"
    # Recursively create subdirectories
    for sub_dir in "$dir"/*; do
      if [ -d "$sub_dir" ]; then
        sub_main_dir=${sub_dir##*/}
        new_dir="$main_dir/$sub_main_dir"
        install -d %{buildroot}/%{terraform}/"$new_dir"
      fi
    done
  fi
done

# Create all the directory inside SRC docs path
for dir in %{docs_src}/*; do
  if [ -d "$dir" ]; then
    main_dir="docs"/${dir##*/}
    install -d %{buildroot}/%{terraform}/"$main_dir"
    # Recursively create subdirectories
    for sub_dir in "$dir"/*; do
      if [ -d "$sub_dir" ]; then
        sub_main_dir=${sub_dir##*/}
        new_dir="$main_dir/$sub_main_dir"
        install -d %{buildroot}/%{terraform}/"$new_dir"
      fi
    done
  fi
done

# Install files
install -d %{buildroot}/%{terraform}/examples/provider
install %{examples_src}/provider/provider.tf   %{examples_dst}/provider/provider.tf 
install %{source_dir}/bin/terraform-provider-hitachi %{buildroot}/%{terraform}/bin/terraform-provider-hitachi
install %{docs_src}/index.md %{docs_dst}/index.md

# Data-source and resource files
for file in %{examples_src}/data-sources/*/*.tf; do
  dir_name=$(dirname "$file")
  last_path_name=${dir_name##*/}
  file_name=$(basename "$file")
  install "$file" %{examples_dst}/data-sources/"$last_path_name"/"$file_name"
done

for file in %{examples_src}/resources/*/*.tf; do
  dir_name=$(dirname "$file")
  last_path_name=${dir_name##*/}
  file_name=$(basename "$file")
  install "$file" %{examples_dst}/resources/"$last_path_name"/"$file_name"
done

# Docs files
for file in %{docs_src}/*/*.md; do
  dir_name=$(dirname "$file")
  last_path_name=${dir_name##*/}
  file_name=$(basename "$file")
  install "$file" %{docs_dst}/"$last_path_name"/"$file_name"
done

# Ensure .tf and .md files are not executable
find %{buildroot}/%{terraform}/examples -type f -name "*.tf" -exec chmod -x {} \;
find %{buildroot}/%{terraform}/docs -type f -name "*.md" -exec chmod -x {} \;

%files
%defattr(-,root,root,-)

# Top-level directory
%dir %{terraform}

# Binaries
%{terraform}/bin/terraform-provider-hitachi

# Examples and docs (include full tree recursively)
%{terraform}/examples
%{terraform}/docs

%post
chmod 755 -R %{hitachi_base}

echo "  Installing : terraform plugin for root"
mkdir -p /root/%{terraform_plugin}
ln -sf %{terraform}/bin/terraform-provider-hitachi  ${HOME}/%{terraform_plugin}/terraform-provider-hitachi 

tuser=$(logname)
if [[ $tuser != root ]]; then
  echo "  Installing : terraform plugin for ${tuser}"
  entry=$(grep "^${tuser}" /etc/passwd )
  entry="${entry%:*}"
  home="${entry##*:}"
  mkdir -p ${home}/%{terraform_plugin}
  ln -sf %{terraform}/bin/terraform-provider-hitachi  ${home}/%{terraform_plugin}/terraform-provider-hitachi
fi

if [[ $1 -eq 1 ]]; then
    echo "Installation complete"
fi


%preun

%postun
if [[ $1 -eq 0 ]]; then
  # Remove plugin for the user
  tuser=$(logname)
  if [[ $tuser != root ]]; then
    echo "  Erasing : terraform plugin %{_DISPLAY_VERSION} for ${tuser}"
    home=$(getent passwd "$tuser" | cut -d: -f6)
    rm -rf "${home}/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/%{_DISPLAY_VERSION}" || true
    rmdir --ignore-fail-on-non-empty "${home}/.terraform.d/plugins/localhost/hitachi-vantara/hitachi" || true
  fi

  # Remove plugin for root
  echo "  Erasing : terraform plugin %{_DISPLAY_VERSION} for root"
  rm -rf "/root/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/%{_DISPLAY_VERSION}" || true
  rmdir --ignore-fail-on-non-empty "/root/.terraform.d/plugins/localhost/hitachi-vantara/hitachi" || true

  # Remove install dir
  rm -rf %{terraform} || true
  rmdir --ignore-fail-on-non-empty %{hitachi_base} || true

  echo "Erase complete"
fi

%changelog
