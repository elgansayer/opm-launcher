BINARY_NAME=opm-launcher
BIN_DIR=bin

.PHONY: all build-all clean linux-amd64 linux-arm64 windows-amd64 darwin-amd64 darwin-arm64

all: build-all

build-all: linux-amd64 linux-arm64 windows-amd64 darwin-amd64 darwin-arm64

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 main.go

linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-linux-arm64 main.go

windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 main.go

darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 main.go

clean:
	rm -rf $(BIN_DIR)
