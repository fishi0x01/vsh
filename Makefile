.PHONY: test

APP_NAME := vsh
PLATFORMS := linux darwin
ARCHS := 386 amd64
BRANCH := $(shell git branch | grep \* | cut -d ' ' -f2)
TAG := $(shell git tag -l --points-at HEAD)
VERSION := $(shell [ -n "$(TAG)" ] && echo -n "$(TAG)" || echo -n "$(BRANCH)-SNAPSHOT")

cross-compile: clean
	mkdir -p ./build/
	for GOOS in $(PLATFORMS); do \
	  for GOARCH in $(ARCHS); do \
	  	export GOOS=$$GOOS; \
		export GOARCH=$$GOARCH; \
	    go build -ldflags "-X main.vshVersion=$(VERSION)" -o build/${APP_NAME}_$${GOOS}_$${GOARCH}; \
	  done \
	done
	ls build/

compile: clean
	go build -ldflags "-X main.vshVersion=$(VERSION)" -o build/${APP_NAME}_linux_amd64

integration-test:
	./test/kv2.sh
	./test/kv1.sh

clean:
	rm ./build/* || true
