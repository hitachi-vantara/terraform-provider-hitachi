#!/bin/bash
# Terraform Storage build script.

BUILD_MODE=${1:-Release}
echo "Build Mode: ${BUILD_MODE}"

TERRAFORM_NAME=HV_Storage_Terraform
TERRAFORM_VERSION=02.5.0
TERRAFORM_PKG="${TERRAFORM_NAME}-${TERRAFORM_VERSION}"
TERRAFORM_SOURCE_TAR="${TERRAFORM_PKG}.tar.gz"

TERRAFORM_DIR=$(pwd)
RPMBUILD_DIR=$(pwd)/rpmbuild


# Build mode: Debug or Release(default)
echo "Starting terraform provider build..."
make build


echo; echo "Preparing ${RPMBUILD_DIR}..."
rm -rf   ${RPMBUILD_DIR} || true
mkdir -p ${RPMBUILD_DIR}/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
mkdir -p ${RPMBUILD_DIR}/${TERRAFORM_PKG}/{bin,examples,docs}

# Populate rpmbuild with terraform files
echo; echo "Copying files to ${RPMBUILD_DIR}..."
cp     ${TERRAFORM_DIR}/spec/*.spec ${RPMBUILD_DIR}/SPECS
cp -rf ${TERRAFORM_DIR}/examples ${RPMBUILD_DIR}/${TERRAFORM_PKG}
cp -rf ${TERRAFORM_DIR}/docs ${RPMBUILD_DIR}/${TERRAFORM_PKG}
cp -f  ${TERRAFORM_DIR}/terraform-provider-hitachi ${RPMBUILD_DIR}/${TERRAFORM_PKG}/bin

# example: HV_Storage_Terraform-02.5.0.tar.gz
cd ${RPMBUILD_DIR}

tar -czf SOURCES/${TERRAFORM_SOURCE_TAR} ${TERRAFORM_PKG}


# RELEASE version
echo "%_topdir ${RPMBUILD_DIR}" > ~/.rpmmacros
RPMARGS="--target=x86_64  -bb"

echo; echo "Starting rpm build for ${BUILD_MODE} version..."
cd ${RPMBUILD_DIR}
rpmbuild ${RPMARGS} --define "_BUILD ${BUILD_MODE}" -v SPECS/terraform-el7.spec
rm ~/.rpmmacros || true

cd ${TERRAFORM_DIR}
cp ./rpmbuild/RPMS/x86_64/*.rpm ${TERRAFORM_DIR}
echo "Finished build rpm for ${BUILD_MODE} version..."

rm -rf   ${RPMBUILD_DIR} || true

rm -rf   terraform-provider-hitachi