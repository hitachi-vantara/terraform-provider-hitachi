name:          HV_Storage_Terraform
Version:       02.0.0
Release:       1%{?dist}
Summary:       Hitachi Vantara storage plugin provider for terraform
Vendor:        Hitachi Vantara

Group:         Adapters 
License:       hiAdapterLicense
URL:           http://www.hitachivantara.com 
Source0:       HV_Storage_Terraform-02.0.0.tar.gz        
ExclusiveArch: x86_64
BuildRoot:     %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)
AutoReqProv:   no 

#BuildRequires:  
#Requires:       

%description
Hitachi Terraform RPM Package


%pre

#
# Check for OS distribution and minimum version
# RHEL family distributions only.
#
MIN_MAJOR_VER=7
MIN_MINOR_VER=0

os_file="/etc/os-release"
centos_file="/etc/centos-release"

# Collect OS platform and version info.
if [[ -e ${os_file} ]]; then
    . ${os_file}
else
    echo "error: ${os_file} not found" >&2
    exit 1
fi

# CentOS 7 only provides full os version in /etc/centos-release file.
if [[ ${ID} == "CentOS Linux" && ${VERSION} == "7 "* ]]; then
    if [[ -e $centos_file ]]; then
        raw=$(cat ${centos_file} )
        version=( $raw )
        VERSION=${version[3]}
    else
        echo "ERROR: ${centos_file} not found" >&2
        exit 1
    fi
fi
# Check if OS is compatible.
case ${ID} in
    'centos');;
    'ol');;
    'rhel');;
    *)
      echo "ERROR: Unsupported OS platform: ${NAME}" >&2
      echo "Supported platforms are: CentOS, OEL, and Red Hat." >&2
      exit 1;;
esac

# Suppress tabs in VERSION. Expected version format: Major.Minor.
v=${VERSION//\t/ }
v=( ${v//./ } )
if [[ ${#v[@]} -lt 2 ]]; then
    echo "ERROR: Invalid OS version: ${VERSION}" >&2
    exit 1
fi
MAJOR_VER=${v[0]}
MINOR_VER=${v[1]}

if [[ ( ${MAJOR_VER} -lt ${MIN_MAJOR_VER} ) ||  ( ${MINOR_VER} -lt ${MIN_MINOR_VER} ) ]]; then
    echo "ERROR: Terraform is only supported on ${MIN_MAJOR_VER}.${MIN_MINOR_VER} version or later"
    exit 1
fi

# Inhibit debug package collection for release builds.
%if "%{_BUILD}" == "Release"
  %define debug_package %{nil}
%endif


# XX: Disable Build ID checking in RedHat macro file.
%undefine _missing_build_ids_terminate_build


%prep
%setup -q


%build
%define hitachi_base /opt/hitachi
%define terraform %{hitachi_base}/terraform
%define plugin_dir  .terraform.d/plugins
%define terraform_plugin %{plugin_dir}/localhost/hitachi-vantara/hitachi/2.0/linux_amd64

%define source_dir %{_topdir}/HV_Storage_Terraform-02.0.0
%define examples_src %{source_dir}/examples
%define docs_src %{source_dir}/docs

%define examples_dst %{buildroot}/%{terraform}/examples
%define docs_dst %{buildroot}/%{terraform}/docs



%install
rm -rf $RPM_BUILD_ROOT

install -d %{buildroot}/%{terraform}
install -d %{buildroot}/%{terraform}/bin
install -d %{buildroot}/%{terraform}/log

#Create all the directory inside SRC examples path
for dir in %{examples_src}/*; do

  # Check if the directory is a directory
  if [ -d "$dir" ]; then

    # Print the directory name
    main_dir="examples"/${dir##*/}
    install -d %{buildroot}/%{terraform}/"$main_dir"
    # Recursively print the subdirectories of the current directory
    for sub_dir in "$dir"/*; do
      if [ -d "$sub_dir" ]; then
        sub_main_dir=${sub_dir##*/}
        new_dir="$main_dir/$sub_main_dir"
        install -d %{buildroot}/%{terraform}/"$new_dir"
      fi
    done

  fi
done

#Create all the directory inside SRC docs path
for dir in %{docs_src}/*; do

  # Check if the directory is a directory
  if [ -d "$dir" ]; then

    # Print the directory name
    main_dir="docs"/${dir##*/}
    install -d %{buildroot}/%{terraform}/"$main_dir"
    # Recursively print the subdirectories of the current directory
    for sub_dir in "$dir"/*; do
      if [ -d "$sub_dir" ]; then
        sub_main_dir=${sub_dir##*/}
        new_dir="$main_dir/$sub_main_dir"
        install -d %{buildroot}/%{terraform}/"$new_dir"
      fi
    done

  fi
done

install -d %{buildroot}/%{terraform}/examples/provider

install %{examples_src}/provider/provider.tf   %{examples_dst}/provider/provider.tf 
install %{source_dir}/bin/terraform-provider-hitachi %{buildroot}/%{terraform}/bin/terraform-provider-hitachi

install %{docs_src}/index.md %{docs_dst}/index.md

# Iterate over the list of data-source files and install them
for file in %{examples_src}/data-sources/*/*.tf; do

  # Get the directory name of the file
  dir_name=$(dirname "$file")
  last_path_name=${dir_name##*/}
  # Install the file
  file_name=$(basename "$file")
  # Install the file
  install "$file" %{examples_dst}/data-sources/"$last_path_name"/"$file_name"

done

# Iterate over the list of resource files and install them
for file in %{examples_src}/resources/*/*.tf; do

  # Get the directory name of the file
  dir_name=$(dirname "$file")
  last_path_name=${dir_name##*/}
  file_name=$(basename "$file")
  # Install the file
  install "$file" %{examples_dst}/resources/"$last_path_name"/"$file_name"

done

# Iterate over the list of docs files and install them
for file in %{docs_src}/*/*.md; do

  # Get the directory name of the file
  dir_name=$(dirname "$file")
  last_path_name=${dir_name##*/}
  file_name=$(basename "$file")
  # Install the file
  install "$file" %{docs_dst}/"$last_path_name"/"$file_name"

done

%define terraform %{hitachi_base}/terraform

%define mytffiles %{_builddir}/mytffiles.txt

#list all files and store in a temporary file 
for file in %{terraform}/examples/*/*/*.tf; do
    echo "$file" >> %{mytffiles}
done

%define mydocfiles %{_builddir}/mydocfiles.txt

#list all files and store in a temporary file 
for file in %{terraform}/docs/*/*.md; do
    echo "$file" >> %{mydocfiles}
done

%files -f %{mytffiles}
%files -f %{mydocfiles}
%dir %{terraform}
%defattr(-,root,root,-)

%{terraform}/bin/terraform-provider-hitachi
%{terraform}/examples/provider/provider.tf
%{terraform}/docs/index.md

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
# Reserved for future expansion.
#if [[ $1 -eq 0 ]]; then
#    # Do some tasks prior to uninstall
#    systemctl stop puma
#elif [[ $1 -eq 1 ]]; then
#    # Do some tasks prior to upgrade
#    echo "Preparing uninstall for upgrading..."
#fi

%postun
if [[ $1 -eq 0 ]]; then
  # Erase terraform plugin for user.
  tuser=$(logname)
  if [[ $tuser != root ]]; then
    echo "  Erasing : terraform plugin for ${tuser}"
    entry=$(grep "^${tuser}" /etc/passwd )
    entry=$(grep "^${tuser}" /etc/passwd )
    entry="${entry%:*}"
    home="${entry##*:}"
    rm -rf ${home}/%{plugin_dir}/hitachi-vantara   || true
    rmdir --ignore-fail-on-non-empty ${home}/.terraform.d/plugins || true
    rmdir --ignore-fail-on-non-empty ${home}/.terraform.d         || true
  fi
  # Erase terraform plugin for root.
  echo "  Erasing : terraform plugin for root"
  rm -rf ${HOME}/%{plugin_dir}/hitachi-vantara   || true
  rmdir --ignore-fail-on-non-empty ${HOME}/.terraform.d/plugins || true
  rmdir --ignore-fail-on-non-empty ${HOME}/.terraform.d         || true

  # Delete /opt/hitachi/terraform-provider-hitachi
  #echo "  Erasing %{terraform}"
  rm -rf %{terraform} || true

  # Delete /opt/hitachi if empty.
  #echo "  Erasing %{hitachi_base}"
  rmdir --ignore-fail-on-non-empty %{hitachi_base} || true
  echo "Erase complete"

# Reserved for future expansion.
#elif [[ $1 -eq 1 ]]; then
#    # Do some tasks after upgrade
#    echo "Perform some uninstalled post-tasks for upgrading..."
#fi
fi

%changelog

