# LeBonPoint Go Server

Go implementation of the LeBonPoint marketplace API using Gin framework.

## Prerequisites

- Go 1.21 or higher
- CGO enabled (required for go-sqlite3)
- GCC/Clang compiler (for CGO)

### Installing Go

Download from [golang.org](https://golang.org/dl/) or use your package manager.

### Verifying CGO

```bash
go env CGO_ENABLED
# Should output "1"
```

If CGO is disabled, enable it:
```bash
export CGO_ENABLED=1
```

On macOS, install Xcode Command Line Tools:
```bash
xcode-select --install
```

On Ubuntu/Debian:
```bash
sudo apt-get install build-essential
```

## Installation

1. Navigate to the server directory:
```bash
cd server/go
```

2. Install dependencies:
```bash
go mod download
```

## Running the Server

### Development Mode

Run directly (requires FTS5 build tag):
```bash
go run -tags fts5 cmd/server/main.go
```

Or use the Makefile:
```bash
make run
```

### Production Build

Build the binary (requires FTS5 build tag for full-text search):
```bash
go build -tags fts5 -o bin/server cmd/server/main.go
```

Or use the Makefile:
```bash
make build
```

Run the binary:
```bash
./bin/server
```

**Note**: The `-tags fts5` flag is required to enable SQLite FTS5 full-text search support. Without it, the search functionality will not work.

## Makefile

The project includes a Makefile for common tasks:

```bash
make build   # Build the server with FTS5 support
make clean   # Remove built binaries
make run     # Run the server with FTS5 support
```

## Configuration

The server uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8081 |
| `DATABASE_PATH` | Path to SQLite database | ../database/db.sqlite |
| `ENV` | Environment (development/production) | development |

Create a `.env` file to override defaults (see `.env.example`).

## Testing

Run all tests (requires FTS5 build tag):
```bash
go test -tags fts5 ./... -v
```

Run tests with coverage:
```bash
go test -tags fts5 -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Development

### Code Formatting
```bash
go fmt ./...
```

### Static Analysis
```bash
go vet ./...
```

### Linting (if golangci-lint is installed)
```bash
golangci-lint run
```

## Docker Support

Build Docker image:
```bash
docker build -t lebonpoint-go .
```

Run container:
```bash
docker run -p 8081:8081 -v $(pwd)/../database:/app/database lebonpoint-go
```
