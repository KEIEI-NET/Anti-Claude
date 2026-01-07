# Suggested Commands

## Testing
```bash
# Run all tests
go test ./... -v

# Run specific package tests
go test ./internal/domain/nippou/... -v
go test ./internal/usecase/nippou/... -v

# Run with coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

## Code Quality
```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run static analysis (if staticcheck installed)
staticcheck ./...
```

## Build
```bash
# Build the project
go build ./...

# Clean build cache
go clean -cache
```

## Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```
