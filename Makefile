# Makefile for Terraform Provider

TEST ?= $$(go list ./... | grep -v 'vendor')

HOSTNAME = localhost
NAMESPACE = hitachi-vantara
NAME = hitachi
BINARY = terraform-provider-${NAME}

# Use the TERRAFORM_VERSION and BUILD_NUMBER passed from build.sh, or use default
VERSION := $(or $(TERRAFORM_VERSION),2.1)
BUILD_NUMBER := $(or $(BUILD_NUMBER),1)

# Full version string: e.g., 2.0.x
SEMVER := ${VERSION}
# Full version string with build number: e.g., 2.0.x-1
FULL_VERSION := ${SEMVER}-${BUILD_NUMBER}

OS_ARCH = x86_64
LINUX_OS_ARCH = linux_amd64

GOVERSION = 1.22

.DEFAULT_GOAL := all

# called from dev environment
.PHONY: all
all: build install

# called from build.sh
.PHONY:  build
build:  mod telemetry
	@echo "🔧 Building provider version ${SEMVER}"
	go build -ldflags="-X main.version=${SEMVER}" -o ${BINARY}
	echo "${FULL_VERSION}" > version.txt

.PHONY:  telemetry
telemetry:
	go run hitachi/common/telemetry/createmap/createmap.go

# called from dev environment
.PHONY:  config
config:
	@echo "📦 creating config.json"
	(cd hitachi/common/config/main; go run create_config.go)
	@echo "📦 copying config.json to /opt/hitachi/terraform"
	mkdir -p /opt/hitachi/terraform
	mkdir -p /opt/hitachi/terraform/bin
	mkdir -p /opt/hitachi/terraform/telemetry
	cp hitachi/common/config/main/config.json /opt/hitachi/terraform
	cp hitachi/common/telemetry/user_consent.sh /opt/hitachi/terraform/bin
	cp ${BINARY} /opt/hitachi/terraform/bin
	cp -r docs /opt/hitachi/terraform
	cp -r examples /opt/hitachi/terraform

.PHONY: release
release:
	goreleaser release --rm-dist --snapshot --skip-publish --skip-sign

# called from dev environment
.PHONY: install
install: build config
	@echo "📦 Installing provider to local plugin path"
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${SEMVER}/${OS_ARCH}
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${SEMVER}/${LINUX_OS_ARCH}
	cp ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${SEMVER}/${OS_ARCH}/
	cp ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${SEMVER}/${LINUX_OS_ARCH}/

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
	@echo "Terraform Provider Makefile"
	@echo "  make                     - Build and install the provider"
	@echo "  make all                 - Build and install the provider, called from dev environment"
	@echo "  make build               - Compile the provider binary, called from build.sh"
	@echo "  make install             - Install to ~/.terraform.d/plugins"
	@echo "  make release             - Run goreleaser (snapshot)"
	@echo "  make test                - Run unit tests"
	@echo "  make testacc             - Run acceptance tests"
	@echo "  make clean               - Clean up build artifacts"
