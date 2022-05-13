OUT_DIR=output

all: saturn

prepare:
	mkdir -p output

saturn: prepare
	go build -o ${OUT_DIR}/saturn ./cmd/saturn