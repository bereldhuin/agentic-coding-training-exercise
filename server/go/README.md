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

## API Endpoints

### Health Check
- `GET /health` - Server health status

### Items
- `GET /v1/items` - List items with filtering, sorting, pagination
- `GET /v1/items/search` - Full-text search using FTS5 (requires `-tags fts5`)
- `POST /v1/items` - Create a new item
- `GET /v1/items/:id` - Get item by ID
- `PUT /v1/items/:id` - Update item (full replacement)
- `PATCH /v1/items/:id` - Update item (partial update)
- `DELETE /v1/items/:id` - Delete item

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

## Architecture

The project follows hexagonal architecture (ports and adapters):

```
server/go/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/                  # Business logic (framework-agnostic)
│   │   ├── entity/              # Item entity
│   │   ├── repository/          # Repository interfaces (ports)
│   │   └── valueobject/         # ItemImage, enums
│   ├── application/             # Use cases
│   │   ├── usecase/            # Business logic orchestration
│   │   └── dto/                # Data transfer objects
│   ├── infrastructure/          # External adapters
│   │   ├── persistence/        # SQLite repository
│   │   └── http/               # Gin routes, handlers, middleware
│   └── shared/                 # Shared utilities
│       ├── config/             # Configuration
│       ├── errors/             # Custom errors
│       ├── cursor/             # Pagination utilities
│       └── time/               # Time utilities
├── test/                       # Integration tests
└── pkg/                        # Public packages
```

### Layer Responsibilities

- **Domain**: Core business logic, entities, and repository interfaces
- **Application**: Use cases that orchestrate business operations
- **Infrastructure**: External concerns (database, HTTP, etc.)
- **Shared**: Common utilities and configuration

## OpenAPI Compliance

This server implements the canonical OpenAPI specification at `../api/openapi.yaml`.

All endpoints, schemas, validation rules, and error responses match the specification exactly.

## Database

The server uses SQLite with the following features:

- Shared database with other implementations (TypeScript, Python, Swift, Kotlin)
- **Full-text search using FTS5** (requires `-tags fts5` build flag)
- Cursor-based pagination
- JSON field for images array

### FTS5 Full-Text Search

This server uses SQLite's FTS5 extension for efficient full-text search on items. The search indexes the `title` and `description` fields.

**Implementation Details**:
- Virtual table: `items_fts` (auto-created on first run)
- Triggers automatically sync the FTS5 table on INSERT/UPDATE/DELETE
- Search uses `MATCH` operator with relevance ranking
- Query terms are automatically escaped for safe FTS5 queries

**Building with FTS5**:
```bash
# Required for full-text search to work
go build -tags fts5 -o bin/server cmd/server/main.go
```

The `mattn/go-sqlite3` driver requires FTS5 to be explicitly enabled at compile time via build tags.

## Error Handling

All errors follow the canonical format:
```json
{
  "error": {
    "code": "error_code",
    "message": "Human readable message",
    "details": {}
  }
}
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

## License

MIT
