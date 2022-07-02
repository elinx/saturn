OUT_DIR=output

all: saturn

prepare:
	mkdir -p output

saturn: prepare
	go mod tidy
	go build -o ${OUT_DIR}/saturn -trimpath ./cmd/saturn