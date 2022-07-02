OUT_DIR=output

all: saturn

prepare: test fmt
	mkdir -p output

saturn: prepare
	go mod tidy
	go build -o ${OUT_DIR}/saturn -trimpath ./cmd/saturn

unit-tests:
	go clean -testcache
	go test ./...

fmt:
	go fmt ./...
