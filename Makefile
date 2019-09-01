.PHONY: test

APP_NAME := vsh
PLATFORMS := linux darwin
ARCHS := 386 amd64

cross-compile: clean
	mkdir -p ./build/
	for GOOS in $(PLATFORMS); do \
	  for GOARCH in $(ARCHS); do \
	  	export GOOS=$$GOOS; \
		export GOARCH=$$GOARCH; \
	    go build -o build/${APP_NAME}_$${GOOS}_$${GOARCH}; \
	  done \
	done
	ls build/

compile: clean
	go build -o build/${APP_NAME}_linux_amd64

integration-test:
	./test/kv2.sh
	./test/kv1.sh

clean:
	rm ./build/* || true
