# Makefile for mahin CLI

BINARY_NAME=mahin
GO=go
GOFLAGS=-ldflags="-s -w"
WINRES=go-winres
ICON=assets/icon.ico

# Default: build for current OS
.PHONY: build
build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

# Embed Windows icon resource
.PHONY: winres
winres:
	$(GO) install github.com/tc-hib/go-winres@v0.3.3
	$(WINRES) init --icon $(ICON)

# Build for Windows (with embedded icon)
.PHONY: build-windows
build-windows: winres
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME).exe .
	rm -f *.syso
	rm -rf winres

# Build for Windows arm64 (with embedded icon)
.PHONY: build-windows-arm
build-windows-arm: winres
	GOOS=windows GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME).exe .
	rm -f *.syso
	rm -rf winres

# Build for macOS (Apple Silicon)
.PHONY: build-mac-arm
build-mac-arm:
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

# Build for macOS (Intel)
.PHONY: build-mac-intel
build-mac-intel:
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

# Build for Linux (amd64)
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

# Build for Linux (arm64)
.PHONY: build-linux-arm
build-linux-arm:
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

# Cross-compile all 6 targets into dist/
.PHONY: build-all
build-all: winres
	mkdir -p dist
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-windows-arm64.exe .
	rm -f *.syso
	rm -rf winres
	GOOS=linux   GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	@echo "Built all targets in dist/"

# Run
.PHONY: run
run: build
	./$(BINARY_NAME)

# Clean
.PHONY: clean
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe *.syso
	rm -rf winres dist

# Install locally (copies to /usr/local/bin on Mac/Linux)
.PHONY: install
install: build
	cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Tidy modules
.PHONY: tidy
tidy:
	$(GO) mod tidy
