# Anti-Claude Project Overview

## Purpose
A Salesforce MCP (Model Context Protocol) server implemented in Go.

## Tech Stack
- **Language:** Go 1.25.3
- **Module:** `salesforce-mcp-server`
- **Dependencies:** 
  - `github.com/google/uuid` - UUID generation
  - Node.js dependencies for MCP SDK

## Architecture
Clean Architecture with:
- **Domain Layer** (`internal/domain/`) - Entities, Value Objects, Repository Interfaces
- **UseCase Layer** (`internal/usecase/`) - Application Services
- Repository pattern with ISP-compliant interfaces (Reader/Writer)

## Code Style & Conventions
- Go standard formatting
- DDD (Domain-Driven Design) patterns
- Builder pattern for complex entity creation
- Structured errors with error codes
- Comprehensive nil-safety
- Defensive copies for immutability
- Table-driven tests

## Commands
```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Format code
go fmt ./...
```

## Project Structure
```
internal/
├── domain/
│   └── nippou/          # Nippou domain model
│       ├── nippou.go    # Entity, Value Objects, Repository interface
│       └── nippou_test.go
└── usecase/
    └── nippou/          # Nippou use cases
        ├── create.go    # Create UseCase
        └── create_test.go
```
