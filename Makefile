build-all: build-wasm build-js
.PHONY: build-all

build-wasm:
	@echo "Generating wasm..."
	@rm -rf ./build
	@mkdir ./build
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./build
	GOOS=js GOARCH=wasm go build -o ./build/main.wasm ./wasm/main.go
.PHONY: build-wasm

build-js:
	./scripts/inlinewasm.js ./build/main.wasm
	yarn && yarn build
.PHONY: build-js

test-all: test-go test-js
.PHONY: test-all

test-go:
	go test -v ./wasm/...
.PHONY: test-go

test-js:
	yarn test
.PHONY: test-js

install:
	go install github.com/agnivade/wasmbrowsertest@v0.7.0
	mv $(shell go env GOPATH)/bin/wasmbrowsertest $(shell go env GOPATH)/bin/go_js_wasm_exec
.PHONY: install
