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
	go build -ldflags "-X main.vshVersion=$(VERSION)" -o build/${APP_NAME}_linux_amd64

get-bats:
	rm -rf bats-tests/bin/
	mkdir -p bats-tests/bin/core
	mkdir -p bats-tests/bin/plugins/bats-assert
	mkdir -p bats-tests/bin/plugins/bats-support
	curl -sL https://github.com/bats-core/bats-core/archive/v1.2.0.tar.gz | tar xvz --strip 1 -C bats-tests/bin/core
	curl -sL https://github.com/bats-core/bats-assert/archive/v2.0.0.tar.gz | tar xvz --strip 1 -C bats-tests/bin/plugins/bats-assert 
	curl -sL https://github.com/bats-core/bats-support/archive/v0.3.0.tar.gz | tar xvz --strip 1 -C bats-tests/bin/plugins/bats-support 

integration-tests:
	bats-tests/run.sh

clean:
	rm ./build/* || true
