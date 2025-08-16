APP_NAME := vsh
SUPPORTED_PLATFORMS := linux darwin
SUPPORTED_ARCHS := amd64 arm64
VERSION := $(shell git describe --tags --always --dirty)

UNAME_M := $(shell uname -m)
ARCH := $(UNAME_M)
ifeq ($(UNAME_M),x86_64)
	ARCH=amd64
endif
ifneq ($(filter %86,$(UNAME_M)),)
	ARCH=386
endif
ifneq ($(filter arm%,$(UNAME_M)),)
  ARCH=arm
endif
ifneq ($(filter $(UNAME_M),arm64 aarch64 armv8b armv8l),)
	ARCH=arm64
endif

help: ## Prints help for targets with comments
	@grep -E '^[a-zA-Z0-9.\ _-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## Run golangci-lint for linting code
	golangci-lint run

format: ## format go code
	gofmt -s -w $(shell find . -name '*.go' -not -path "./vendor/*")
	gci write $(shell find . -name '*.go' -not -path "./vendor/*")
	golines -w $(shell find . -name '*.go' -not -path "./vendor/*")
	goimports -w $(shell find . -name '*.go' -not -path "./vendor/*")

compile-releases: clean ## Compile vsh binaries for multiple platforms and architectures strictly using vendor directory
	mkdir -p ./build/
	for GOOS in $(SUPPORTED_PLATFORMS); do \
		for GOARCH in $(SUPPORTED_ARCHS); do \
			GOOS=$$GOOS GOARCH=$$GOARCH \
				go build -mod vendor -trimpath -ldflags "-s -w -X main.vshVersion=$(VERSION)" -o build/${APP_NAME}_$${GOOS}_$${GOARCH}; \
				go build -mod vendor -tags 'notokenhelper' -trimpath -ldflags "-s -w -X main.vshVersion=$(VERSION)" -o build/${APP_NAME}_notokenhelper_$${GOOS}_$${GOARCH}; \
		done \
	done
	cd build/ && sha256sum * > SHA256SUM

compile: ## Compile vsh for platform based on uname
	go build -ldflags "-X main.vshVersion=$(VERSION)" -o build/${APP_NAME}_$(shell uname | tr '[:upper:]' '[:lower:]')_$(ARCH)

get-bats: ## Download bats dependencies to test directory
	rm -rf test/bin/
	mkdir -p test/bin/core
	mkdir -p test/bin/plugins/bats-assert
	mkdir -p test/bin/plugins/bats-support
	mkdir -p test/bin/plugins/bats-file
	curl -sL https://github.com/bats-core/bats-core/archive/v1.12.0.tar.gz | tar xvz --strip 1 -C test/bin/core
	curl -sL https://github.com/bats-core/bats-assert/archive/v2.0.0.tar.gz | tar xvz --strip 1 -C test/bin/plugins/bats-assert
	curl -sL https://github.com/bats-core/bats-support/archive/v0.3.0.tar.gz | tar xvz --strip 1 -C test/bin/plugins/bats-support
	curl -sL https://github.com/bats-core/bats-file/archive/v0.4.0.tar.gz | tar xvz --strip 1 -C test/bin/plugins/bats-file

integration-tests: ## Run integration test suites (requires bats - see get-bats)
	test/run-all-tests.sh

single-test: ## Run a single test suite, e.g., make single-test KV_BACKEND=KV2 VAULT_VERSION=1.6.1 TEST_SUITE=commands/cp
	KV_BACKEND=$(KV_BACKEND) VAULT_VERSION=$(VAULT_VERSION) TEST_SUITE=$(TEST_SUITE) test/run-single-test.sh

local-vault-standard-test-instance: ## Start a local vault container with standard setup provisioning
	bash -c ". test/util/common.bash && . test/util/standard-setup.bash && setup"

local-vault-concurrency-test-instance: ## Start a local vault container with concurrency setup provisioning
	bash -c ". test/util/common.bash && . test/util/concurrency-setup.bash && setup"

clean: ## Remove builds and vsh related docker containers
	docker rm -f vsh-integration-test-vault || true
	rm ./build/* || true

update-dependencies: ## Update all dependencies
	go get -u

.PHONY: vendor
vendor: ## synch dependencies in vendor/ directory
	go mod tidy
	go mod vendor

reinit-go-mod: ## remove and re-create mod files and vendor dir
	rm go.mod
	rm go.sum
	go mod init github.com/fishi0x01/vsh
	go mod tidy
	go mod vendor

