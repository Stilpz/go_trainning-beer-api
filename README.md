# Beer API - Golang Training Project

A RESTful API for managing beer inventory built with Go, following Domain-Driven Design principles and clean architecture patterns. This training project demonstrates best practices in Go development, including dependency injection, database migrations, testing, and API documentation.

## 🎯 Project Overview

This API is structured around domain-driven design, with a focus on:
- **Business domain encapsulation**: Each business context (beer) is isolated in its own module
- **Clean architecture**: Clear separation between handlers, services, repositories, and models
- **Dependency injection**: Using Uber's Dig for managing dependencies
- **Comprehensive testing**: Unit tests with mocks using Mockery and go-sqlmock
- **API documentation**: Auto-generated Swagger/OpenAPI documentation

## 🏗️ Architecture & Project Structure

```
.
├── beer/                      # Beer domain module
│   ├── handler/              # HTTP handlers (controllers)
│   ├── service/              # Business logic layer
│   ├── repository/           # Data access layer
│   ├── model/                # Domain models/entities
│   ├── interfaces/           # Interface definitions
│   ├── external/             # External service integrations
│   └── mocks/                # Generated mocks for testing
├── cmd/                      # Application entry points
│   └── main.go              # Main application
├── configs/                  # Configuration packages
│   ├── generals/            # General configurations
│   │   ├── injector/        # Dependency injection setup
│   │   └── router/          # Route definitions
│   └── storage/             # Database configuration
│       ├── connection.go    # DB connection setup
│       ├── migration.go     # Migration management
│       └── migrations/      # SQL migration files
├── docs/                     # Auto-generated Swagger docs
├── pkg/                      # Shared packages/utilities
│   └── kit/                 # Common utilities (logger, etc.)
├── Makefile                  # Build and development commands
└── go.mod                    # Go module definition
```

## 🚀 Tech Stack

### Core Technologies
- **[Go 1.23+](https://golang.org/dl/)** - Programming language
- **[Echo Framework](https://github.com/labstack/echo)** - High-performance HTTP web framework
- **[PostgreSQL](https://www.postgresql.org/)** - Primary database

### Dependencies & Libraries
- **[Uber Dig](https://pkg.go.dev/go.uber.org/dig)** - Dependency injection container
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Database migration management
- **[zerolog](https://github.com/rs/zerolog)** - Fast and structured logging
- **[godotenv](https://github.com/joho/godotenv)** - Environment variable management
- **[pq](https://github.com/lib/pq)** - PostgreSQL driver
- **[Swag](https://github.com/swaggo/swag)** - Swagger documentation generator

### Testing & Development Tools
- **[Testify](https://github.com/stretchr/testify)** - Testing toolkit with assertions
- **[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)** - SQL driver mock for testing
- **[Mockery](https://github.com/vektra/mockery)** - Mock generator

## 📋 Prerequisites

Before running this project, ensure you have the following installed:

- **Go 1.23 or higher** - [Download Go](https://golang.org/dl/)
- **PostgreSQL** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Docker** (optional, recommended for containerization)
- **Docker Compose** (optional, for orchestration)
- **Make** (for using Makefile commands)

## ⚙️ Setup & Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Stilpz/go_trainning-beer-api.git
cd go_trainning-beer-api
```

### 2. Configure Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Server Configuration
API_PORT=8888

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=beer_api_db
DB_SSLMODE=disable

# Logger Configuration
LOGGER_DEBUG=true
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Install Development Tools

```bash
# Install Swagger/Swag
make install-swag

# Install Mockery for generating mocks
make mockery-install

# Install GolangCI-Lint (Windows with Chocolatey)
make install-lint-windows-chocolatey
```

### 5. Setup Database

Ensure PostgreSQL is running and create the database:

```sql
CREATE DATABASE beer_api_db;
```

The application will automatically run migrations on startup.

## 🏃 Running the Application

### Standard Execution

```bash
go run cmd/main.go
```

### With Debug Logging

```bash
go run cmd/main.go -debug=true
```

The API will be available at: `http://localhost:8888`

## 📚 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check endpoint |
| `GET` | `/docs/*` | Swagger API documentation |
| `GET` | `/` | List all beers |
| `GET` | `/:beerID` | Get beer details by ID |
| `GET` | `/:beerID/box-price` | Calculate box price for a beer |
| `POST` | `/` | Create a new beer |

### Swagger Documentation

Access the interactive API documentation at:
```
http://localhost:8888/docs/index.html
```

## 🧪 Testing

### Run All Tests with Coverage

```bash
make go-test
```

This command:
- Runs all tests in the `beer` module (excluding mocks)
- Generates coverage report (`coverage.out`)
- Displays coverage summary

### Generate HTML Coverage Report

```bash
make go-test-report
```

Opens the coverage report in your default browser.

### Generate Mocks

```bash
make mockery
```

Generates mocks based on `.mockery.yml` configuration.

## 🛠️ Development Workflow

### Generate/Update Swagger Documentation

```bash
make swag
```

This command:
- Generates Swagger docs from code annotations
- Creates OpenAPI JSON and YAML files
- Updates the `docs/` directory

### Available Make Commands

| Command | Description |
|---------|-------------|
| `make swag` | Generate Swagger documentation |
| `make install-swag` | Install Swag tool |
| `make mockery-install` | Install Mockery tool |
| `make mockery` | Generate mocks |
| `make go-test` | Run tests with coverage |
| `make go-test-report` | Generate HTML coverage report |

## 🔧 Configuration Files

- **`.mockery.yml`** - Mockery configuration for mock generation
- **`Makefile`** - Build automation and development tasks
- **`.env`** - Environment variables (create from example above)
- **`.gitignore`** - Git ignore rules

## 🧩 Key Features

### Dependency Injection
Uses Uber Dig for automatic dependency resolution and injection, making the codebase more maintainable and testable.

### Database Migrations
Automatic database migrations on startup using golang-migrate, with version control for schema changes.

### Structured Logging
zerolog provides fast, structured JSON logging with configurable debug levels.

### API Documentation
Auto-generated Swagger/OpenAPI documentation from code annotations, keeping docs in sync with implementation.

### Comprehensive Testing
Unit tests with mocks, integration testing support, and coverage reporting.

## 📝 Contributing

This is a training project. Feel free to fork and experiment!

## 📄 License

This project is for educational purposes.

## 👥 Author

**Stilpz** - [GitHub Profile](https://github.com/Stilpz)

---

**Training Organization**: Dropi

**API Version**: 1.0
