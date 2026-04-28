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
  - [Database Migration](#database-migration)
  - [How to run Resik?](#how-to-run-resik)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Environment Variables](#environment-variables)
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
  - Rate Limiting
  - Dependency Injection
  - Database Migration Support
  - Structured Logging

- **Developer Experience**
  - API Documentation
  - Clean Architecture Implementation
  - Configuration Management

## Getting Started

### Prerequisites

- Ensure you have **[Go](https://go.dev/)** (version 1.26 or higher) installed on your system.
- You will need **[Make](https://www.gnu.org/software/make/)** to run build commands.
- If you plan to work with gRPC definitions, you'll need to install **[protoc](https://protobuf.dev/installation/)** (the Protocol Buffer compiler), **[protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)**, and **[protoc-gen-go-grpc](https://github.com/grpc/grpc-go/tree/master/cmd/protoc-gen-go-grpc)**.
- If you intend to use API documentation generation, make sure you have **[swag](https://github.com/swaggo/swag)** installed.
- For dependency injection, ensure **[wire](https://github.com/google/wire)** is installed.

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

For a detailed list of all available configuration options, see the [Environment Variables](#environment-variables) section.

### Database Migration

Before running the application for the first time, or after any database schema changes, you need to run the database migrations. Use the following command:

```bash
go run main.go migrate up
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

### REST API

- **Server Address**:
  - Default: `http://localhost:8081`
  - Configurable in `.env` using: `REST_MODULE_HOST`, `REST_MODULE_PORT`.
- **Swagger UI**:
  - Accessible at the server address (e.g., http://localhost:8081/docs/swagger/index.html) after running the app.
  - Configurable in `.env` using: `REST_MODULE_PUBLIC_SCHEMA`, `REST_MODULE_PUBLIC_HOST`, `REST_MODULE_PUBLIC_PORT`.

### gRPC API

- **Server Address**:
  - Default: `grpc://localhost:7071`
  - Configurable in `.env` using: `GRPC_MODULE_HOST`, `GRPC_MODULE_PORT`.
- **API Definitions (Proto Files)**:
  - Located under the [app/grpc/](/app/grpc/) and [core/](/core/) directories.

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

## Environment Variables

The application can be configured using environment variables. These variables can be set in a `.env` file in the root directory.

### App Configuration
| Variable | Description | Default | Example |
| :--- | :--- | :--- | :--- |
| `APP_NAME` | The name of the application | `Resik` | `Resik` |
| `APP_NAME_KEY` | A unique identifier key for the application | `resik` | `resik` |
| `APP_VERSION` | The version of the application | `0.1` | `0.1` |

### REST Module Configuration
| Variable | Description | Default | Example |
| :--- | :--- | :--- | :--- |
| `REST_MODULE_NAME` | Name of the REST module | - | `REST` |
| `REST_MODULE_NAME_KEY` | Unique identifier key for the REST module | - | `rest` |
| `REST_MODULE_ENV` | Running environment (`development`, `staging`, `production`) | `development` | `development` |
| `REST_MODULE_HOST` | Host address to bind the REST server | - | `127.0.0.1` |
| `REST_MODULE_PORT` | Port number to bind the REST server | - | `8081` |
| `REST_MODULE_PUBLIC_HOST` | Public host address for documentation and URLs | - | `127.0.0.1` |
| `REST_MODULE_PUBLIC_PORT` | Public port number for documentation and URLs | - | `8081` |
| `REST_MODULE_PUBLIC_SCHEMA` | Public protocol schema (`http` or `https`) | - | `http` |
| `REST_MODULE_PUBLIC_BASE_PATH` | Public base path prefix | - | `/` |
| `REST_MODULE_LANG_DEFAULT` | Default language for response messages | `en` | `en` |
| `REST_HTTP_READ_HEADER_TIMEOUT` | Max duration for reading request headers | - | `3s` |
| `REST_HTTP_READ_TIMEOUT` | Max duration for reading the entire request | - | `15s` |
| `REST_HTTP_WRITE_TIMEOUT` | Max duration for writing the response | - | `30s` |
| `REST_HTTP_IDLE_TIMEOUT` | Max duration to wait for the next request | - | `90s` |
| `REST_DATABASE_MAX_OPEN_CONS` | Max number of open database connections | `30` | `30` |
| `REST_DATABASE_MAX_IDLE_CONS` | Max number of idle database connections | `5` | `5` |
| `REST_DATABASE_CONN_MAX_LIFETIME` | Max time a database connection can be reused | `5m` | `5m` |
| `REST_DATABASE_CONN_MAX_IDLETIME` | Max time a database connection can be idle | `5m` | `5m` |
| `REST_DATABASE_IS_DEBUG` | Enable SQL query logging for debugging | `false` | `true` |
| `REST_AUTH_EXPIRES` | Expiration duration for authentication tokens | - | `24h` |
| `REST_LOG_FILE` | Path to the log file for this module | - | `logs/rest.log` |

### gRPC Module Configuration
| Variable | Description | Default | Example |
| :--- | :--- | :--- | :--- |
| `GRPC_MODULE_NAME` | Name of the gRPC module | - | `gRPC` |
| `GRPC_MODULE_NAME_KEY` | Unique identifier key for the gRPC module | - | `grpc` |
| `GRPC_MODULE_ENV` | Running environment (`development`, `staging`, `production`) | `development` | `development` |
| `GRPC_MODULE_HOST` | Host address to bind the gRPC server | - | `127.0.0.1` |
| `GRPC_MODULE_PORT` | Port number to bind the gRPC server | - | `7071` |
| `GRPC_MODULE_LANG_DEFAULT` | Default language for response messages | `en` | `en` |
| `GRPC_DATABASE_MAX_OPEN_CONS` | Max number of open database connections | `30` | `30` |
| `GRPC_DATABASE_MAX_IDLE_CONS` | Max number of idle database connections | `5` | `5` |
| `GRPC_DATABASE_CONN_MAX_LIFETIME` | Max time a database connection can be reused | `5m` | `5m` |
| `GRPC_DATABASE_CONN_MAX_IDLETIME` | Max time a database connection can be idle | `5m` | `5m` |
| `GRPC_DATABASE_IS_DEBUG` | Enable SQL query logging for debugging | `false` | `true` |
| `GRPC_AUTH_EXPIRES` | Expiration duration for authentication tokens | - | `24h` |
| `GRPC_LOG_FILE` | Path to the log file for this module | - | `logs/grpc.log` |

### Shared Database Configuration
| Variable | Description | Default | Example |
| :--- | :--- | :--- | :--- |
| `DATABASE_ENGINE` | Database engine type (`mysql`, `postgres`) | `mysql` | `mysql` |
| `DATABASE_HOST` | Database server host | - | `127.0.0.1` |
| `DATABASE_PORT` | Database server port | - | `3306` |
| `DATABASE_USERNAME` | Database username | - | `user` |
| `DATABASE_PASSWORD` | Database password | - | `password` |
| `DATABASE_SCHEMA` | Database name or schema | - | `dbname` |
| `DATABASE_HEALTHCHECK_TIMEOUT` | Timeout for database health checks | - | `2s` |

### Shared Redis Configuration
| Variable | Description | Default | Example |
| :--- | :--- | :--- | :--- |
| `REDIS_HOST` | Redis server host | - | `127.0.0.1` |
| `REDIS_PORT` | Redis server port | - | `6379` |
| `REDIS_USERNAME` | Redis username (if using ACL) | - | `user` |
| `REDIS_PASSWORD` | Redis password | - | `password` |
| `REDIS_DATABASE` | Redis database index | `0` | `0` |
| `REDIS_POOL_SIZE` | Redis connection pool size | - | `5` |
| `REDIS_HEALTHCHECK_TIMEOUT` | Timeout for Redis health checks | - | `500ms` |

### Rate Limit Configuration
| Variable | Description | Default | Example |
| :--- | :--- | :--- | :--- |
| `RATE_LIMIT_DRIVER` | Storage driver for rate limiting (`redis`, `memory`) | `memory` | `redis` |
| `RATE_LIMIT_PERIOD` | Time window for rate limiting | - | `1m` |
| `RATE_LIMIT_LIMIT` | Maximum requests allowed per period | - | `100` |
| `RATE_LIMIT_PREFIX` | Cache key prefix for rate limiting | - | `rate_limit` |

### Shared Auth Configuration
| Variable | Description | Default | Example |
| :--- | :--- | :--- | :--- |
| `AUTH_SIGNATURE_KEY` | Secret key used for signing JWT tokens | - | `secret` |
| `AUTH_HASH_MEMORY` | Memory limit for Argon2 hashing in KB | `8192` | `65536` |
| `AUTH_HASH_ITERATIONS` | Number of iterations for Argon2 hashing | `3` | `3` |

## Credits

- [Dedy F. Setyawan](https://github.com/dedyf5) (Author)

## License

[MIT](/LICENSE)
