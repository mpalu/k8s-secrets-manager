# Makefile
.PHONY: build test clean run-server run-cli help

BINARY_NAME=k8s-secret-manager
BUILD_DIR=build

help:
	@echo "Kubernetes Secret Manager"
	@echo ""
	@echo "Usage:"
	@echo "  make build      - Compila o projeto"
	@echo "  make test       - Executa os testes"
	@echo "  make clean      - Remove arquivos gerados"
	@echo "  make run-server - Executa em modo servidor"
	@echo "  make run-cli    - Executa em modo CLI"

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