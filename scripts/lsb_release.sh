#! /bin/bash
#
# The lsb_release.sh script provides similar functionality to the EOL'ed lsb_release utility
# found in various Linux Version 7 platforms.  
#
# Version 0.4 by J. Kahn
#   Restrict /etc/centos-release file usage to CentOS7.
# Version 0.5 by J. Kahn
#   Suppress tabs in VERSION variable (/etc/os-release file).
# Version 0.6 by J. Kahn
#   Use rhel not rhat. Redhat is quite attached to R.H.E.L toeken.
# Version 0.7 by J. Kahn
#   Added default display.
#


# options:
lsb_version="0.7"
lsb_id=0
lsb_release=0
lsb_short=0
lsb_default=1
lsb_file="/etc/os-release"
lsb_centos_file="/etc/centos-release"

centos_full="CentOS Linux"
centos_id="centos"
oel_full="Oracle Linux Server"
oel_id="oel"
rhat_full="Red Hat Enterprise Linux"
rhat_id="rhel"


lsb_help() {
    cat << EOF
Usage: lsb_release.sh [options]

Examples: lsb_release.sh -si
          lsb_release.sh -sr
          lsb_release.sh -i
          lsb_release.sh -r
          lsb_release.sh -rm
          lsb_release.sh -rn
          lsb_release.sh

Options:
  -h, --help            Usage information
  -v, --version         Utility version
  -s, --short           Suppress description of result

  -i, --id              Show distribution ID
  -r, --release         Show release number of this distribution
  -rm, --release-major  Show release major number of this distribution
  -rn, --release-minor  Show release minor number of this distribution
EOF
    exit 0
}


# Process lsb_release options
if [[ -e $lsb_file ]]; then
    # Define NAME and VERSION variables.
    . $lsb_file
else
    echo "ERROR: $lsb_file not found" >&2
    exit 1
fi

if [[ $NAME == "CentOS Linux" && $VERSION == "7 "* ]]; then
  if [[ -e $lsb_centos_file ]]; then
    # CentOS7 only provides full os version in /etc/centos-release file.
    raw=$(cat $lsb_centos_file)
    version=( $raw )
    VERSION=${version[3]}
  else
    echo "ERROR: $lsb_centos_file not found" >&2
    exit 1
  fi
fi

while  true; do
    case $1 in
        -h|--help)
	    lsb_help
	    ;;
	-v|--version)
            echo "$lsb_version"
	    exit 0
	    ;;
	-sr|-rs)
	    lsb_short=1
	    lsb_release=3
            lsb_default=0
	    ;;
	-si|-is)
	    lsb_short=1
	    lsb_id=1
            lsb_default=0
	    ;;
	-s|--short)
	    lsb_short=1
            lsb_default=0
	    ;;
	-i|--id)
	    lsb_id=1
            lsb_default=0
	    ;;
        -r|--release)
	    lsb_release=3
            lsb_default=0
	    ;;
        -rm|--release-major)
	    lsb_short=1
	    lsb_release=2
            lsb_default=0
	    ;;
        -rn|--release-minor)
	    lsb_short=1
	    lsb_release=1
            lsb_default=0
	    ;;
	*)
	    break
	    ;;
    esac
    shift
done

#
# Display full Distro Name by default.
# If short option specified, display short name id.
#
if [[ $lsb_default != 0 ]]; then
   lsb_id=1
   lsb_release=1
fi
if [[ $lsb_id != 0 ]]; then
    text=""
    if [[ $lsb_short == 0 ]]; then
        text="Distributor ID: "
    else
	case $NAME in	
	  $oel_full)
            NAME=$oel_id;;
	  $rhat_full)
            NAME=$rhat_id;;
	  $centos_full)
            NAME=$centos_id;;
	  *)
	    echo "Unknown Distributor: $NAME"
	esac
    fi
    text="$text$NAME"
    echo $text
fi

#
# Display full release version by default.
# If short option specified, display major.minor version.
# Suppress tabs in VERSION.
#
if [[ $lsb_release != 0 ]]; then
    text=""
    if [[ $lsb_short == 0 ]]; then
        text="Release: "
        text="$text$VERSION"
    else
        # Only display major.minor short release info.
        v=${VERSION//\t/ }
        v=( ${v//./ } )
	if [[ ${#v[@]} -lt 2 ]]; then
		echo "ERROR: Invalid OS version: ${VERSION}"
		exit 1
	fi
        case $lsb_release in
	    1) ver="${v[1]}" ;;
	    2) ver="${v[0]}" ;;
	    *) ver="${v[0]}.${v[1]}" ;;
	esac
        text="$text$ver"
    fi
    echo $text
fi
