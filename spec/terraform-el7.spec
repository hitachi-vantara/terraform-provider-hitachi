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


%pre
logfile="/var/log/hitachi_terraform_install.log"

# Remove old log if it exists
[ -f "$logfile" ] && rm -f "$logfile"

echo "[$(date)] Starting pre-install checks" | tee -a "$logfile"

# Check if same or different version is installed
installed_version=$(rpm -q --queryformat '%{VERSION}-%{RELEASE}' %{name} 2>/dev/null)
if [ $? -eq 0 ]; then
  echo "[$(date)] ERROR: Version $installed_version of %{name} is already installed." | tee -a "$logfile" >&2
  echo "[$(date)] Please remove or move it before reinstalling. Aborting." | tee -a "$logfile" >&2
  exit 1
fi

# Check if install directory exists
if [ -d "%{terraform}" ]; then
  echo "[$(date)] ERROR: Installation directory %{terraform} already exists." | tee -a "$logfile" >&2
  echo "[$(date)] Please remove or move it before reinstalling. Aborting." | tee -a "$logfile" >&2
  exit 1
fi

# Redirection of any errors generated from RPM scriptlet
# This will suppress the scriptlet error messages (like exit status 1, etc.)
exec 2>/dev/null

echo "[$(date)] Pre-install checks passed" | tee -a "$logfile"


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

# Log file
logfile="/var/log/hitachi_terraform_install.log"

echo "[$(date)] Starting installation of HV_Storage_Terraform" | tee -a "$logfile"

# Warning about overwrite
echo "[$(date)] WARN: Overwriting existing directories under %{terraform}" | tee -a "$logfile" >&2

echo "[$(date)] Installing terraform plugin for root"
mkdir -p /root/%{terraform_plugin}
ln -sf %{terraform}/bin/terraform-provider-hitachi  ${HOME}/%{terraform_plugin}/terraform-provider-hitachi 2>>"$logfile" | tee -a "$logfile"

tuser=$(logname)
if [[ $tuser != root ]]; then
  echo "[$(date)]  Installing terraform plugin for ${tuser}" | tee -a "$logfile"
  entry=$(grep "^${tuser}" /etc/passwd )
  entry="${entry%:*}"
  home="${entry##*:}"
  mkdir -p ${home}/%{terraform_plugin}
  ln -sf %{terraform}/bin/terraform-provider-hitachi  ${home}/%{terraform_plugin}/terraform-provider-hitachi 2>>"$logfile" | tee -a "$logfile"
fi

if [[ $1 -eq 1 ]]; then
    echo "[$(date)] Installation complete" | tee -a "$logfile"
    echo "[$(date)] Installation successful" | tee -a "$logfile"
fi


%postun
logfile="/var/log/hitachi_terraform_uninstall.log"

# Remove old log if it exists
[ -f "$logfile" ] && rm -f "$logfile"

echo "[$(date)] Starting uninstallation of HV_Storage_Terraform" | tee -a "$logfile"

if [[ $1 -eq 0 ]]; then
  echo "[$(date)] WARN: Deleting %{terraform} and contents" | tee -a "$logfile" >&2

  tuser=$(logname)
  if [[ $tuser != root ]]; then
    echo "[$(date)] Erasing : terraform plugin %{_DISPLAY_VERSION} for ${tuser}" | tee -a "$logfile"
    home=$(getent passwd "$tuser" | cut -d: -f6)
    rm -rf "${home}/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/%{_DISPLAY_VERSION}" >> "$logfile" 2>&1 | tee -a "$logfile" || true
    rmdir --ignore-fail-on-non-empty "${home}/.terraform.d/plugins/localhost/hitachi-vantara/hitachi" >> "$logfile" 2>&1 | tee -a "$logfile" || true
  fi

  echo "[$(date)] Erasing terraform plugin %{_DISPLAY_VERSION} for root" | tee -a "$logfile"
  rm -rf "/root/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/%{_DISPLAY_VERSION}" >> "$logfile" 2>&1 | tee -a "$logfile" || true
  rmdir --ignore-fail-on-non-empty "/root/.terraform.d/plugins/localhost/hitachi-vantara/hitachi" >> "$logfile" 2>&1 | tee -a "$logfile" || true

  echo "[$(date)] Removing install directory %{terraform}" | tee -a "$logfile"
  rm -rf %{terraform} >> "$logfile" 2>&1 | tee -a "$logfile" || true
  rmdir --ignore-fail-on-non-empty %{hitachi_base} >> "$logfile" 2>&1 | tee -a "$logfile" || true

  echo "[$(date)] Erase complete" | tee -a "$logfile"
  echo "[$(date)] Uninstallation complete" | tee -a "$logfile"
fi
