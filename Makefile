# Makefile
.PHONY: build test clean run-server run-cli help

BINARY_NAME=k8s-secret-manager
BUILD_DIR=build

help:
	@echo "Kubernetes Secret Manager"
	@echo ""
	@echo "Usage:"
	@echo "  make build      - Compile the project"
	@echo "  make test       - Execute the tests"
	@echo "  make clean      - Remove the generated files"
	@echo "  make run-server - Execute in server mode"
	@echo "  make run-cli    - Execute in CLI mode"
	@echo "  make help       - Show this help message"

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/k8s-secret-manager/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

run-server: build
	@echo "Starting server..."
	@./$(BUILD_DIR)/$(BINARY_NAME) server

run-cli: build
	@echo "Running CLI..."
	@./$(BUILD_DIR)/$(BINARY_NAME) $(filter-out $@,$(MAKECMDGOALS))

%:
	@: