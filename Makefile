# Define variables
BINARY_NAME=bin/app
SRC_DIR=cmd/main
GO_FILES=$(SRC_DIR)/main.go
BUILD_DIR=bin
AIR_BINARY=bin/air

# Default target: build the project
.PHONY: all
all: build

# Build the Go binary
build:
	@echo "Building the project..."
	@go build -o $(BINARY_NAME) $(GO_FILES)
	@echo "Build completed. Binary: $(BINARY_NAME)"

# Run the Go project
run:
	@echo "Running the project..."
	@$(BINARY_NAME)

# Clean up the binary
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@echo "Clean completed."

# Install air for live reloading
install-air:
	@echo "Installing air..."
	@go install github.com/cosmtrek/air@latest
	@cp $(GOPATH)/bin/air $(AIR_BINARY)
	@echo "Air installed."

# Run the project using air for live reloading
air:
	@echo "Running with air..."
	@$(AIR_BINARY)

# Format Go code
fmt:
	@echo "Formatting Go code..."
	@go fmt ./...

# Run tests (if you have tests)
test:
	@echo "Running tests..."
	@go test ./...

# Help command
help:
	@echo "Available targets:"
	@echo "  build        - Build the Go project"
	@echo "  run          - Run the Go project"
	@echo "  clean        - Clean up the binary"
	@echo "  install-air  - Install air for live reloading"
	@echo "  air          - Run the project with air"
	@echo "  fmt          - Format Go code"
	@echo "  test         - Run Go tests"

