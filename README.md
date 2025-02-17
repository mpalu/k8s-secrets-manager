# Kubernetes Secret Manager

A Kubernetes secrets manager that can be used as a CLI tool and as a HTTP server.

## Features

- Secret management via CLI or REST API
- Support for multiple namespaces
- Secret data validation
- Full CRUD operations (Create, Read, Update, Delete)
- Flexible configuration via YAML file
- Structured JSON logging
- Support for in-cluster and out-of-cluster execution

## Installation

### Prerequisites

- Go 1.21 or higher
- Access to a Kubernetes cluster
- Configured kubeconfig file (for out-of-cluster usage)

### Building from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/mpalu/k8s-secrets-manager.git
   cd k8s-secrets-manager
   ```

2. Build the binary:

   ```bash
   make build
   ```

3. Run tests:

   ```bash
   make test
   ```

4. Install the binary:
   ```bash
   sudo mv build/k8s-secret-manager /usr/local/bin/
   ```

## Usage

### CLI Mode

The CLI mode provides direct command-line access to manage Kubernetes secrets.

Available commands:

- `create`: Create a new secret
- `get`: Get a secret
- `update`: Update a secret
- `delete`: Delete a secret

### HTTP Server Mode

The HTTP server mode provides a RESTful API for managing Kubernetes secrets.

Available endpoints:

- `POST /api/v1/secrets`: Create a new secret
- `GET /api/v1/secrets`: Get all secrets
- `GET /api/v1/secrets/{name}`: Get a specific secret
- `PUT /api/v1/secrets/{name}`: Update a secret
- `DELETE /api/v1/secrets/{name}`: Delete a secret

### Configuration

The configuration is done via a YAML file.

```yaml
server:
  port: 8080
  host: "0.0.0.0"

kubernetes:
  inCluster: false
  kubeconfig: "" # Leave empty to use default location
```

### Running the Server

```bash
k8s-secret-manager server -c config.yaml
```

### Running the CLI

```bash
k8s-secret-manager create -n default -s mysecret -k key1=value1
```

## Contributing

Contributions are welcome! Please feel free to submit a PR.
