ci:
	make gen-all && make test

gen-all:
	make gen-wasm && make gen-js

gen-wasm:
	@echo "Generating wasm..."
	rm -rf ./build
	mkdir ./build
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./build
	GOOS=js GOARCH=wasm go build -o ./build/main.wasm ./wasm/main.go

gen-js:
	./scripts/inlinewasm.js ./build/main.wasm
	yarn && yarn build

test:
	yarn test
