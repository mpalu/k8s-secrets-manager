# Makefile
.PHONY: build test lint docker run

BINARY_NAME=k8s-secrets-manager
BUILD_DIR=build

help:
	@echo "Kubernetes Secret Manager"
	@echo ""
	@echo "Usage:"
	@echo "  make build      - Compile the project"
	@echo "  make test       - Execute the tests"
	@echo "  make lint       - Run golangci-lint"
	@echo "  make docker     - Build the Docker image"
	@echo "  make run        - Run the project"
	@echo "  make help       - Show this help message"

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/k8s-secrets-manager/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run

docker:
	@echo "Building Docker image..."
	@docker build -t k8s-secrets-manager .

run:
	@echo "Running project..."
	@go run ./cmd/k8s-secrets-manager/main.go

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