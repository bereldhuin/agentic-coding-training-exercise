# LeBonPoint Swift Server

A Swift Vapor implementation of the LeBonPoint Marketplace API. This server implements the canonical OpenAPI specification defined in `server/api/openapi.yaml` using the hexagonal architecture pattern.

## Overview

This Swift implementation:
- Implements the complete REST API for marketplace items management
- Shares the same SQLite database as the TypeScript and Python servers
- Follows hexagonal architecture (ports and adapters) pattern
- Uses Vapor 4 for modern async/await web framework
- Provides full-text search using SQLite FTS5
- Supports filtering, sorting, and cursor-based pagination

## Prerequisites

- Swift 5.9+ (for async/await support)
- macOS 13+ or Linux
- SQLite3 (for database)

## Installation

1. Navigate to the server directory:
```bash
cd server/swift
```

2. Resolve dependencies:
```bash
swift package resolve
```

3. Build the project:
```bash
swift build
```

## Configuration

Create a `.env` file in the `server/swift` directory (copy from `.env.example`):

```env
PORT=9000
DATABASE_PATH=../database/db.sqlite
ENVIRONMENT=development
LOG_LEVEL=info
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `9000` |
| `DATABASE_PATH` | Path to SQLite database | `../database/db.sqlite` |
| `ENVIRONMENT` | Environment (development, testing, production) | `development` |
| `LOG_LEVEL` | Logging level (debug, info, warning, error, critical) | `info` |

## Running the Server

### Development Mode

```bash
swift run
```

### Production Mode

1. Build the release binary:
```bash
swift build -c release
```

2. Run the binary:
```bash
.build/release/App
```

## Testing

Run all tests:
```bash
swift test
```

Run tests with verbose output:
```bash
swift test --verbose
```

## API Endpoints

The server implements the following endpoints (all prefixed with `/v1`):

### Health Check
- `GET /health` - Server health check

### Items
- `GET /v1/items` - List items with filtering, sorting, and pagination
- `POST /v1/items` - Create a new item
- `GET /v1/items/:id` - Get a specific item
- `PUT /v1/items/:id` - Replace an entire item
- `PATCH /v1/items/:id` - Partially update an item
- `DELETE /v1/items/:id` - Delete an item

### Query Parameters

#### Filtering (GET /v1/items)
- `status` - Filter by status (draft, active, reserved, sold, archived)
- `category` - Filter by category
- `min_price_cents` - Minimum price in cents
- `max_price_cents` - Maximum price in cents
- `city` - Filter by city
- `postal_code` - Filter by postal code
- `is_featured` - Filter by featured status (true/false)
- `delivery_available` - Filter by delivery availability (true/false)

#### Sorting
- `sort` - Sort format: `field:direction` (e.g., `created_at:desc`)
  - Fields: id, title, price_cents, created_at, updated_at, published_at
  - Directions: asc, desc

#### Pagination
- `limit` - Items per page (1-100, default: 20)
- `cursor` - Base64-encoded pagination cursor

#### Search
- `q` - Full-text search query (searches title and description)

## Architecture

This server follows the **hexagonal architecture** (ports and adapters) pattern:

```
Sources/App/
├── Domain/                  # Core business logic (framework-agnostic)
│   ├── Entities/            # Item entity
│   ├── Repositories/        # Repository protocol (port)
│   └── ValueObjects/        # Enums and value objects
├── Application/             # Use cases and services
│   ├── UseCases/           # Business logic orchestration
│   └── Services/           # Validation and domain services
├── Infrastructure/          # External adapters
│   ├── Persistence/        # SQLite repository implementation
│   └── HTTP/               # Vapor routes and controllers
└── Shared/                 # Shared utilities
    ├── Configuration.swift # Config management
    └── Errors.swift        # Custom error types
```

### Layer Responsibilities

1. **Domain Layer**: Core business entities and repository interfaces
2. **Application Layer**: Use cases that orchestrate business logic
3. **Infrastructure Layer**: External concerns (database, HTTP)
4. **Shared**: Cross-cutting concerns (errors, configuration)

## Database

The server uses SQLite with the shared database at `../database/db.sqlite`.

### Schema

The database schema is created by the TypeScript server. The Swift server:
- Uses the existing `items` table
- Uses the FTS5 full-text search table `items_fts`
- Does not perform migrations (assumes schema exists)

## API Compatibility

This Swift implementation is fully compatible with the TypeScript and Python implementations:

- Same OpenAPI specification
- Same database schema
- Same API behavior
- Same error response format
- Same pagination and sorting

## Docker Support

### Build Docker Image

```bash
docker build -t lebonpoint-swift-server .
```

### Run with Docker

```bash
docker run -p 9000:9000 \
  -v $(pwd)/../database:/app/database \
  -e PORT=9000 \
  -e DATABASE_PATH=/app/database/db.sqlite \
  lebonpoint-swift-server
```

### Docker Compose

See the root `docker-compose.yml` for running all servers together.

## Examples

### Create an Item

```bash
curl -X POST http://localhost:9000/v1/items \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Vintage Leather Sofa",
    "description": "Beautiful vintage leather sofa from the 1970s.",
    "price_cents": 35000,
    "category": "furniture",
    "condition": "good",
    "status": "draft",
    "city": "Strasbourg",
    "postal_code": "67000"
  }'
```

### List Items with Filters

```bash
curl "http://localhost:9000/v1/items?status=active&category=electronics&sort=created_at:desc&limit=10"
```

### Search Items

```bash
curl "http://localhost:9000/v1/items?q=iPhone"
```

### Get a Specific Item

```bash
curl http://localhost:9000/v1/items/1
```

### Update an Item (PATCH)

```bash
curl -X PATCH http://localhost:9000/v1/items/1 \
  -H "Content-Type: application/json" \
  -d '{"price_cents": 60000}'
```

### Delete an Item

```bash
curl -X DELETE http://localhost:9000/v1/items/1
```

## Error Responses

All errors follow the canonical format:

```json
{
  "error": {
    "code": "validation_error",
    "message": "Validation failed",
    "details": {
      "title": "Title must be at least 3 characters"
    }
  }
}
```

Error codes:
- `validation_error` - Request validation failed
- `not_found` - Resource not found
- `internal_error` - Server error

## Development

### Adding New Features

1. Define domain entities in `Domain/Entities/`
2. Define repository protocol in `Domain/Repositories/`
3. Implement use cases in `Application/UseCases/`
4. Implement repository in `Infrastructure/Persistence/`
5. Add HTTP models in `Infrastructure/HTTP/Models/`
6. Add controllers in `Infrastructure/HTTP/Controllers/`
7. Register routes in `routes.swift`

### Code Style

- Follow Swift API Design Guidelines
- Use 4 spaces for indentation
- Mark private methods with `private`
- Use `async`/`await` for asynchronous operations
- Prefer protocol-oriented programming

## License

MIT
