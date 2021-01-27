APP_NAME := vsh
PLATFORMS := linux darwin
ARCHS := 386 amd64
VERSION := $(shell git describe --tags --always --dirty)

help: ## Prints help for targets with comments
	@grep -E '^[a-zA-Z0-9.\ _-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

cross-compile: clean ## Compile vsh binaries for multiple platforms and architectures
	mkdir -p ./build/
	for GOOS in $(PLATFORMS); do \
		for GOARCH in $(ARCHS); do \
			GOOS=$$GOOS GOARCH=$$GOARCH \
				go build -ldflags "-X main.vshVersion=$(VERSION)" -o build/${APP_NAME}_$${GOOS}_$${GOARCH}; \
		done \
	done
	ls build/

compile: clean ## Compile vsh for platform based on uname
	go build -ldflags "-X main.vshVersion=$(VERSION)" -o build/${APP_NAME}_$(shell uname | tr '[:upper:]' '[:lower:]')_amd64

get-bats: ## Download bats dependencies to test directory
	rm -rf test/bin/
	mkdir -p test/bin/core
	mkdir -p test/bin/plugins/bats-assert
	mkdir -p test/bin/plugins/bats-support
	mkdir -p test/bin/plugins/bats-file
	curl -sL https://github.com/bats-core/bats-core/archive/v1.2.0.tar.gz | tar xvz --strip 1 -C test/bin/core
	curl -sL https://github.com/bats-core/bats-assert/archive/v2.0.0.tar.gz | tar xvz --strip 1 -C test/bin/plugins/bats-assert
	curl -sL https://github.com/bats-core/bats-support/archive/v0.3.0.tar.gz | tar xvz --strip 1 -C test/bin/plugins/bats-support
	curl -sL https://github.com/bats-core/bats-file/archive/v0.2.0.tar.gz | tar xvz --strip 1 -C test/bin/plugins/bats-file

integration-tests: ## Run integration test suites (requires bats - see get-bats)
	test/run.sh

local-vault-test-instance: ## Start a local vault container with integration test provisioning
	bash -c ". test/util/util.bash && setup"

clean: ## Remove builds and vsh related docker containers
	docker rm -f vsh-integration-test-vault || true
	rm ./build/* || true
