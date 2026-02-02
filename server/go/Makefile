.PHONY: build clean run

BINARY_NAME=go-server
OUTPUT_DIR=../../bin
BUILD_FLAGS=-tags "fts5"

build:
	go build $(BUILD_FLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/server

clean:
	rm -f $(OUTPUT_DIR)/$(BINARY_NAME)

run:
	go run $(BUILD_FLAGS) ./cmd/server
