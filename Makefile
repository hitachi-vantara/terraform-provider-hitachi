# Makefile for Terraform
#
# If multple CPU cores are available, use them to accelerate build process.
#MAKEFLAGS += -j$(shell nproc)
#
# Use sodo for go as it can be in /usr/local/bin.
# Chck go version as it is a moving target.
#

TEST?=$$(go list ./... | grep -v 'vendor')

HOSTNAME=localhost
NAMESPACE=hitachi-vantara
NAME=hitachi
BINARY=terraform-provider-${NAME}
VERSION?=2.0
PATCH_VERSION?=0
BUILD_NUMBER?=1
OS_ARCH=x86_64
LINUX_OS_ARCH=linux_amd64

GOVERSION=1.20

.DEFAULT_GOAL := all

.PHONY: all
all: build install

.PHONY:  build
build:  mod
	go build -o ${BINARY}
	echo "${VERSION}.${PATCH_VERSION}-${BUILD_NUMBER}" > version.txt

.PHONY: release
release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${LINUX_OS_ARCH}
	cp ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${LINUX_OS_ARCH}

.PHONY: test
test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

.PHONY: testacc
testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   

.PHONY: mod
mod:
	go mod tidy --compat=${GOVERSION}


.PHONY: clean
clean:
	rm -rf rpmlib rpmbuild
	rm -f *.rpm version.txt

.PHONY: help
help:
	@echo ""
	@echo "Makefile Options:"
	@echo "  make                           - Build and install all components"
	@echo "  make all                       - Build and install all components"
	@echo "  make build                     - Build only"
	@echo "  make release                   - Build release version"
	@echo " "

