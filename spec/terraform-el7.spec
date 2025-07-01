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


%pre
logdir="/var/log/hitachi/terraform"
logfile="$logdir/hitachi_terraform_install.log"

# Ensure the log directory exists
mkdir -p "$logdir"
chmod 755 "$logdir"

# Remove old log if it exists
[ -f "$logfile" ] && rm -f "$logfile"

echo "[$(date)] Starting pre-install checks" | tee -a "$logfile"

installed_ver=$(rpm -qa %{name})
if [ "${installed_ver}" != "" ]; then
  echo "[$(date)] ERROR: $installed_ver is already installed." | tee -a "$logfile" >&2
  echo "[$(date)] Please uninstall it before reinstalling. Aborting." | tee -a "$logfile" >&2
  exit 1
fi

# Redirection of any errors generated from RPM scriptlet
# This will suppress the scriptlet error messages (like exit status 1, etc.)
exec 2>/dev/null

echo "[$(date)] Pre-install checks passed" | tee -a "$logfile"


%install
rm -rf $RPM_BUILD_ROOT

# Define terraform install root
install -d %{buildroot}/%{terraform}

# Copy everything from extracted source tree to terraform install location
cp -a * %{buildroot}/%{terraform}
find %{buildroot}/%{terraform}/docs -type f -name '*.md' -exec chmod 0644 {} \;
find %{buildroot}/%{terraform}/examples -type f -name '*.tf' -exec chmod 0644 {} \;
find %{buildroot}/%{terraform} -type f -name '*.json' -exec chmod 0644 {} \;


%files
%defattr(-,root,root)
%{terraform}


%post

# Log file
logdir="/var/log/hitachi/terraform"
logfile="$logdir/hitachi_terraform_install.log"

# Ensure the log directory exists
mkdir -p "$logdir"
chmod 755 "$logdir"

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

  %define user_consent_message %(sed 's/["`$\\]/\\&/g; s/^/echo "/; s/$/"/' BUILD/user_consent_message.txt | paste -sd';' -)
  echo
  %{user_consent_message}
  echo
fi


%postun
logdir="/var/log/hitachi/terraform"
logfile="$logdir/hitachi_terraform_uninstall.log"

# Ensure the log directory exists
mkdir -p "$logdir"
chmod 755 "$logdir"

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
