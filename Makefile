# cogmit Makefile

BINARY_NAME=cogmit
VERSION?=0.1.0
BUILD_DIR=build
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.PHONY: build clean install test help

# Default target
all: build

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✅ Cross-platform build complete"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...
	@echo "✅ Tests complete"

# Install the binary to /usr/local/bin
install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installation complete"

# Uninstall the binary
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstall complete"

# Run the application
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Setup development environment
dev-setup:
	@echo "🛠️  Setting up development environment..."
	@go mod tidy
	@go mod download
	@echo "✅ Development setup complete"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@go vet ./...
	@echo "✅ Linting complete"

# Show help
help:
	@echo "cogmit Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build        Build the binary"
	@echo "  build-all    Build for multiple platforms"
	@echo "  clean        Clean build artifacts"
	@echo "  test         Run tests"
	@echo "  install      Install binary to /usr/local/bin"
	@echo "  uninstall    Remove binary from /usr/local/bin"
	@echo "  run          Build and run the application"
	@echo "  dev-setup    Setup development environment"
	@echo "  fmt          Format code"
	@echo "  lint         Lint code"
	@echo "  help         Show this help message"
