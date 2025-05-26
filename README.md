# Resik - Golang Clean Architecture Implementation

[![Go Version](https://img.shields.io/github/go-mod/go-version/dedyf5/resik)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dedyf5/resik/ci.yml)](https://github.com/dedyf5/resik/actions)

Clean Architecture implementation in Golang for building REST and gRPC applications. This project provides various ready-to-use features for creating structured and maintainable applications.

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [How to run Resik?](#how-to-run-resik)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Credits](#credits)
- [License](#license)

## Features

- **Dual Protocol Support**
  - REST API
  - gRPC API

- **Core Functionality**
  - Multi-language Support (English, Bahasa Indonesia, 日本語)
  - Form Validation
  - JWT Authentication & Authorization
  - Dependency Injection
  - Structured Logging

- **Developer Experience**
  - API Documentation
  - Clean Architecture Implementation
  - Configuration Management

## Getting Started

### Prerequisites

- Ensure you have **[Go](https://go.dev/)** (version 1.23 or higher) installed on your system.
- You will need **[Make](https://www.gnu.org/software/make/)** to run build commands.
- If you plan to work with gRPC definitions, you'll need to install **[protoc](https://protobuf.dev/installation/)** (the Protocol Buffer compiler), **[`protoc-gen-go`](https://github.com/protocolbuffers/protobuf-go)**, and **[`protoc-gen-go-grpc`](https://github.com/grpc/grpc-go/tree/master/cmd/protoc-gen-go-grpc)**.
- If you intend to use API documentation generation, make sure you have **[`swag`](https://github.com/swaggo/swag)** installed.
- For dependency injection, ensure **[`wire`](https://github.com/google/wire)** is installed.

### Installation

1. Clone this repository:

    ```bash
    git clone https://github.com/dedyf5/resik.git
    cd resik
    go mod tidy
    ```

2. **(Optional)** Generate necessary files using Make if you intend to modify proto definitions, API documentation, or dependency injection setup:

    ```bash
    make generate
    ```

    This command will:
    - Generate gRPC files.
    - Generate API documentation.
    - Generate Wire dependency injection code.

### Configuration

Copy `.env.example` to `.env` and modify the values:

```bash
cp .env.example .env
```

### How to run Resik?

#### Run REST app

```bash
go run main.go rest
```

#### Run gRPC app

```bash
go run main.go grpc
```

#### Help

```bash
go run main.go --help
```

## API Documentation

- REST API: [Swagger UI](http://localhost:8081/swagger/index.html) (after running the app)
- gRPC: [View Proto Files](/app/grpc/proto/) (link to proto definitions)

## Testing

You can run the tests using the following command:

### Using `make test`

```bash
make test
```

### Using `go test`

```bash
go test ./...
```

## Credits

- [Dedy F. Setyawan](https://github.com/dedyf5) (Author)

## License

[MIT](/LICENSE)
