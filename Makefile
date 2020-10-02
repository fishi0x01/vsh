APP_NAME := vsh
PLATFORMS := linux darwin
ARCHS := 386 amd64
VERSION := $(shell git describe --tags --always --dirty)

cross-compile: clean
	mkdir -p ./build/
	for GOOS in $(PLATFORMS); do \
		for GOARCH in $(ARCHS); do \
			GOOS=$$GOOS GOARCH=$$GOARCH \
				go build -ldflags "-X main.vshVersion=$(VERSION)" -o build/${APP_NAME}_$${GOOS}_$${GOARCH}; \
		done \
	done
	ls build/

compile: clean
	go build -ldflags "-X main.vshVersion=$(VERSION)" -o build/${APP_NAME}_$(shell uname | tr '[:upper:]' '[:lower:]')_amd64

get-bats:
	rm -rf test/bin/
	mkdir -p test/bin/core
	mkdir -p test/bin/plugins/bats-assert
	mkdir -p test/bin/plugins/bats-support
	mkdir -p test/bin/plugins/bats-file
	curl -sL https://github.com/bats-core/bats-core/archive/v1.2.0.tar.gz | tar xvz --strip 1 -C test/bin/core
	curl -sL https://github.com/bats-core/bats-assert/archive/v2.0.0.tar.gz | tar xvz --strip 1 -C test/bin/plugins/bats-assert
	curl -sL https://github.com/bats-core/bats-support/archive/v0.3.0.tar.gz | tar xvz --strip 1 -C test/bin/plugins/bats-support
	curl -sL https://github.com/bats-core/bats-file/archive/v0.2.0.tar.gz | tar xvz --strip 1 -C test/bin/plugins/bats-file

integration-tests:
	test/run.sh

local-vault-test-instance:
	bash -c ". test/util/util.bash && setup"

clean:
	rm ./build/* || true
