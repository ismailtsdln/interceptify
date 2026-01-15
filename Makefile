.PHONY: build clean build-all

BINARY_NAME=interceptify
VERSION=1.0.0

build:
	go build -o $(BINARY_NAME) main.go

build-all:
	# Linux
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 main.go
	# macOS (Intel)
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 main.go
	# macOS (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64 main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe main.go

clean:
	rm -f $(BINARY_NAME)
	rm -rf bin/
