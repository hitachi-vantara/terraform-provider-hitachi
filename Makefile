# Makefile for Terraform Provider

TEST ?= $$(go list ./... | grep -v 'vendor')

HOSTNAME = localhost
NAMESPACE = hitachi-vantara
NAME = hitachi
BINARY = terraform-provider-${NAME}

# Use the TERRAFORM_VERSION and BUILD_NUMBER passed from build.sh, or use default
VERSION := $(or $(TERRAFORM_VERSION),2.1.1)
BUILD_NUMBER := $(or $(BUILD_NUMBER),1)

# Full version string: e.g., 2.0.x
SEMVER := ${VERSION}
# Full version string with build number: e.g., 2.0.x-1
FULL_VERSION := ${SEMVER}-${BUILD_NUMBER}

OS_ARCH = x86_64
LINUX_OS_ARCH = linux_amd64

GOVERSION = 1.24

.DEFAULT_GOAL := all

OPT_TERRAFORM=/opt/hitachi/terraform
INTERNAL_CONFIG_FILE=.internal_config

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
	mkdir -p $(OPT_TERRAFORM)
	mkdir -p $(OPT_TERRAFORM)/bin
	mkdir -p $(OPT_TERRAFORM)/telemetry
	@echo "📦 creating $(INTERNAL_CONFIG_FILE)"
	(cd hitachi/common/config/internal_config; go run create_internal_config.go $(INTERNAL_CONFIG_FILE))
	@echo "📦 copying $(INTERNAL_CONFIG_FILE) to $(OPT_TERRAFORM)/bin"
	cp hitachi/common/config/internal_config/$(INTERNAL_CONFIG_FILE) $(OPT_TERRAFORM)/bin
	cp hitachi/common/telemetry/user_consent.sh $(OPT_TERRAFORM)/bin
	cp scripts/logbundle.sh $(OPT_TERRAFORM)/bin
	cp ${BINARY} $(OPT_TERRAFORM)/bin
	cp -r docs $(OPT_TERRAFORM)
	chmod 0755 examples/data-sources/*/clean.sh
	chmod 0755 examples/resources/*/clean.sh
	cp -r --preserve=mode examples $(OPT_TERRAFORM)

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
