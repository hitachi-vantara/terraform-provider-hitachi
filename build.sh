#!/bin/bash
# Terraform Storage build script.

display_usage() {
    echo "Format to run the script: ${0} <build_number>"
}

# Build versioning
DISPLAY_VERSION="2.2.0"         # For plugin and Makefile
RPM_VERSION="02.2.0"            # For RPM spec & filename
BUILD_NUMBER=$1
BUILD_MODE="Release"


if [[ -z $BUILD_NUMBER ]]; then
    echo "Missing parameter"
    display_usage
    exit 1
elif [[ $1 == --help || $1 == -h ]]; then
    display_usage
    exit 1
fi

echo "Build Mode: ${BUILD_MODE}"
echo "Build Number: ${BUILD_NUMBER}"
echo "Display Version: ${DISPLAY_VERSION}"
echo "RPM Version: ${RPM_VERSION}"

TERRAFORM_NAME="HV_Storage_Terraform"
TERRAFORM_PKG="${TERRAFORM_NAME}-${RPM_VERSION}"
TERRAFORM_SOURCE_TAR="${TERRAFORM_PKG}.tar.gz"

TERRAFORM_DIR=$(pwd)
RPMBUILD_DIR=$(pwd)/rpmbuild

# Pass TERRAFORM_VERSION and BUILD_NUMBER to Makefile
export TERRAFORM_VERSION=${DISPLAY_VERSION}
export BUILD_NUMBER=${BUILD_NUMBER}

# Build mode: Debug or Release(default)
echo "Starting terraform provider build..."
make build

echo; echo "Preparing ${RPMBUILD_DIR}..."
rm -rf ${RPMBUILD_DIR} || true
mkdir -p ${RPMBUILD_DIR}/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
mkdir -p ${RPMBUILD_DIR}/${TERRAFORM_PKG}/{bin,examples,docs,telemetry}

# Populate rpmbuild with terraform files
echo; echo "Copying files to ${RPMBUILD_DIR}..."
cp ${TERRAFORM_DIR}/spec/*.spec ${RPMBUILD_DIR}/SPECS
chmod 0755 ${TERRAFORM_DIR}/examples/data-sources/*/clean.sh
chmod 0755 ${TERRAFORM_DIR}/examples/resources/*/clean.sh
cp -rf ${TERRAFORM_DIR}/examples ${RPMBUILD_DIR}/${TERRAFORM_PKG}
cp -rf ${TERRAFORM_DIR}/examples ${RPMBUILD_DIR}/${TERRAFORM_PKG}
cp -rf ${TERRAFORM_DIR}/docs ${RPMBUILD_DIR}/${TERRAFORM_PKG}
cp -f ${TERRAFORM_DIR}/terraform-provider-hitachi ${RPMBUILD_DIR}/${TERRAFORM_PKG}/bin
cp -f ${TERRAFORM_DIR}/scripts/logbundle.sh ${RPMBUILD_DIR}/${TERRAFORM_PKG}/bin
cp -f ${TERRAFORM_DIR}/hitachi/common/telemetry/user_consent.sh ${RPMBUILD_DIR}/${TERRAFORM_PKG}/bin
# for ex in ${RPMBUILD_DIR}/${TERRAFORM_PKG}/examples/{data-sources,resources}/*; do
#     cp -f ${TERRAFORM_DIR}/hitachi/common/telemetry/user_consent_message.tf $ex
# done

CONSENT_MESSAGE_FOR_RPM_SPEC="${RPMBUILD_DIR}/BUILD/user_consent_message.txt"
echo; echo "Creating for rpm spec the file ${CONSENT_MESSAGE_FOR_RPM_SPEC}"
(cd ${TERRAFORM_DIR}/hitachi/common/config/spec_text; go run create_consent_spec_text.go ${CONSENT_MESSAGE_FOR_RPM_SPEC})

INTERNAL_CONFIG_FOR_RPM="${RPMBUILD_DIR}/${TERRAFORM_PKG}/bin/.internal_config"
echo; echo "Creating for rpm the file ${INTERNAL_CONFIG_FOR_RPM}"
(cd ${TERRAFORM_DIR}/hitachi/common/config/internal_config; go run create_internal_config.go ${INTERNAL_CONFIG_FOR_RPM})

echo; echo "Creating tarball..."
cd ${RPMBUILD_DIR}
tar -czf SOURCES/${TERRAFORM_SOURCE_TAR} ${TERRAFORM_PKG}

# Set the RPM build environment
echo "%_topdir ${RPMBUILD_DIR}" > ~/.rpmmacros
RPMARGS="--target=x86_64 -bb"

# Start RPM build for the specified version
echo; echo "Starting rpm build for ${BUILD_MODE} version..."
# echo rpmbuild ${RPMARGS} --define "_BUILD ${BUILD_MODE}"  --define "_VERSION ${RPM_VERSION}" --define "_BUILD_NUMBER ${BUILD_NUMBER}" -v SPECS/terraform-el7.spec

set -x
rpmbuild ${RPMARGS} \
  --define "_BUILD ${BUILD_MODE}" \
  --define "_VERSION ${RPM_VERSION}" \
  --define "_DISPLAY_VERSION ${DISPLAY_VERSION}" \
  --define "_BUILD_NUMBER ${BUILD_NUMBER}" \
  SPECS/terraform-el7.spec

# Check if the RPM build was successful
if [ $? -eq 0 ]; then
    echo "RPM build successful."
else
    echo "Error: RPM build failed."
    exit 1
fi
set +x

# Clean up the rpmmacros after the build
rm ~/.rpmmacros || true

# Copy the RPM to the original directory
echo; echo "Copying RPM to ${TERRAFORM_DIR}..."
cd ${RPMBUILD_DIR}/RPMS/x86_64/
RPM_FILE=$(ls *.rpm)

if [ -z "$RPM_FILE" ]; then
    echo "Error: No RPM file found in ${RPMBUILD_DIR}/RPMS/x86_64/"
    exit 1
else
    cp ${RPM_FILE} ${TERRAFORM_DIR}
    echo "RPM file copied to ${TERRAFORM_DIR}/${RPM_FILE}"
fi

echo "Finished build rpm for ${BUILD_MODE} version..."

# Optional clean-up
# rm -rf ${RPMBUILD_DIR} || true
# rm -rf terraform-provider-hitachi
