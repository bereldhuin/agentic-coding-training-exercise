# LeBonPoint Python Server (FastAPI)

Python FastAPI implementation of the LeBonPoint marketplace items API.

## API Specification

**This server implements the canonical OpenAPI specification located at `../api/openapi.yaml`.**

The OpenAPI specification is the single source of truth for the API contract. All endpoints, Pydantic models, and validation rules defined in this implementation conform to the canonical specification.

**Note**: FastAPI auto-generates an OpenAPI spec (available at `/docs`, `/redoc`, and `/openapi.json`). The auto-generated spec is useful for development but the canonical spec at `../api/openapi.yaml` is authoritative.

## Overview

This is a Python implementation of the marketplace items API using FastAPI. It provides the same REST API endpoints as the TypeScript implementation, sharing the same SQLite database.

## Architecture

The project follows **hexagonal architecture** (ports and adapters pattern):

```
src/
├── domain/              # Business logic (framework-agnostic)
│   ├── entities/        # Item entity, enums
│   ├── repositories/    # Repository port (abstract interface)
│   └── value_objects/   # Value objects
├── application/         # Use cases and services
│   └── use_cases/      # CreateItem, GetItem, ListItems, etc.
├── infrastructure/      # External adapters
│   ├── persistence/    # SQLite repository implementation
│   └── http/           # FastAPI routes and models
└── shared/             # Shared utilities
    ├── config.py       # Configuration management
    ├── errors.py       # Custom error classes
    └── cursor.py       # Cursor encoding/decoding
```

## Features

- **Full CRUD API**: Create, read, update, patch, and delete items
- **Full-text search**: SQLite FTS5 search with BM25 ranking
- **Filtering**: Filter by status, category, price range, location, etc.
- **Sorting**: Sort by any field with ascending/descending order
- **Pagination**: Cursor-based pagination for efficient data retrieval
- **Validation**: Pydantic models for request/response validation
- **OpenAPI**: Auto-generated API documentation (Swagger UI & ReDoc)
- **Hexagonal architecture**: Clean separation of concerns

## Requirements

- Python 3.11+
- SQLite database (shared with TypeScript implementation)

## Installation

1. **Install Poetry (if not already available):**

```bash
python -m pip install --user poetry
```

2. **Install project dependencies:**

```bash
cd server/python
poetry install --with dev
```

Poetry will create a reproducible lockfile (`poetry.lock`). You can activate the virtualenv with `poetry shell` or run commands via `poetry run <cmd>`.

## Configuration

The server can be configured via environment variables or a `.env` file:

```bash
# Server configuration
PORT=8000
HOST=0.0.0.0

# Database path (relative to server/python directory)
DATABASE_PATH=../database/db.sqlite

# Environment
ENVIRONMENT=development
```

## Running the Server

### Development mode (with auto-reload):

```bash
poetry run uvicorn src.main:app --reload --host 0.0.0.0 --port 8000
```

### Production mode:

```bash
poetry run uvicorn src.main:app --host 0.0.0.0 --port 8000 --workers 4
```

Or using the built-in main:

```bash
poetry run python -m src.main
```

## API Documentation

Once the server is running:

- **Swagger UI**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc
- **OpenAPI JSON**: http://localhost:8000/openapi.json

## API Endpoints

### Items

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/items` | List items with filtering, sorting, pagination |
| POST | `/v1/items` | Create a new item |
| GET | `/v1/items/{id}` | Get a single item |
| PUT | `/v1/items/{id}` | Update (replace) an item |
| PATCH | `/v1/items/{id}` | Partially update an item |
| DELETE | `/v1/items/{id}` | Delete an item |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |

### Query Parameters (List Items)

| Parameter | Type | Description |
|-----------|------|-------------|
| `q` | string | Full-text search query |
| `status` | string | Filter by status |
| `category` | string | Filter by category |
| `min_price_cents` | integer | Minimum price in cents |
| `max_price_cents` | integer | Maximum price in cents |
| `city` | string | Filter by city |
| `postal_code` | string | Filter by postal code |
| `is_featured` | boolean | Filter by featured status |
| `delivery_available` | boolean | Filter by delivery availability |
| `sort_by` | string | Field to sort by |
| `sort_order` | string | Sort order (asc or desc) |
| `limit` | integer | Items per page (max 100) |
| `cursor` | string | Pagination cursor |

## Testing

### Run all tests:

```bash
pytest
```

### Run with coverage:

```bash
pytest --cov=src --cov-report=html
```

### Run specific test file:

```bash
pytest tests/unit/domain/test_enums.py
```

### Run with verbose output:

```bash
pytest -v
```

## API-First Development

This project follows an API-first development approach:

1. The canonical OpenAPI specification is at `../api/openapi.yaml`
2. All API changes should start with updating the specification
3. FastAPI's auto-generated spec should match the canonical spec
4. See `../api/README.md` for the complete API development workflow

## Code Quality

## Database

The server uses the shared SQLite database at `../database/db.sqlite` (created by the TypeScript server).

### Database Schema

The `items` table structure:

```sql
CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL CHECK(length(title) >= 3 AND length(title) <= 200),
    description TEXT,
    price_cents INTEGER NOT NULL CHECK(price_cents >= 0),
    category TEXT,
    condition TEXT NOT NULL CHECK(condition IN ('new', 'like_new', 'good', 'fair', 'parts', 'unknown')),
    status TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'active', 'reserved', 'sold', 'archived')),
    is_featured INTEGER NOT NULL DEFAULT 0,
    city TEXT,
    postal_code TEXT,
    country TEXT NOT NULL DEFAULT 'FR',
    delivery_available INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    published_at TEXT,
    images TEXT NOT NULL DEFAULT '[]'
);
```

### Full-Text Search

The server uses SQLite FTS5 for full-text search with BM25 ranking. Search is performed on the `title` and `description` fields.

## Item Model

### Item Fields

- `id` (integer): Unique identifier
- `title` (string): 3-200 characters
- `description` (string, optional): Item description
- `price_cents` (integer): Price in cents (>= 0)
- `category` (string, optional): Item category
- `condition` (enum): `new`, `like_new`, `good`, `fair`, `parts`, `unknown`
- `status` (enum): `draft`, `active`, `reserved`, `sold`, `archived`
- `is_featured` (boolean): Featured flag
- `city` (string, optional): City name
- `postal_code` (string, optional): Postal code
- `country` (string): Country code (default: "FR")
- `delivery_available` (boolean): Delivery availability
- `created_at` (string): ISO 8601 timestamp
- `updated_at` (string): ISO 8601 timestamp
- `published_at` (string, optional): ISO 8601 timestamp
- `images` (array): List of image objects

### Image Object

- `url` (string, required): Image URL
- `alt` (string, optional): Alt text
- `sort_order` (integer, optional): Sort order

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "code": "error_code",
    "message": "Human-readable message",
    "details": {}
  }
}
```

### Common Error Codes

- `validation_error`: Request validation failed (400)
- `not_found`: Resource not found (404)
- `internal_error`: Internal server error (500)

## Development

### Project Structure

```
server/python/
├── src/                    # Source code
│   ├── domain/            # Domain layer
│   ├── application/       # Application layer
│   ├── infrastructure/    # Infrastructure layer
│   ├── shared/           # Shared utilities
│   └── main.py           # Application entry point
├── tests/                # Test suite
│   ├── unit/            # Unit tests
│   └── integration/     # Integration tests
├── pyproject.toml       # Project configuration
├── poetry.lock          # Poetry lockfile (generated by `poetry install`)
├── .env-example        # Environment variables template
└── README.md           # This file
```

### Adding New Features

1. **Domain**: Define entities and value objects in `src/domain/`
2. **Repository port**: Define interface in `src/domain/repositories/`
3. **Repository implementation**: Implement in `src/infrastructure/persistence/`
4. **Use case**: Create use case in `src/application/use_cases/`
5. **HTTP layer**: Add Pydantic models and routes in `src/infrastructure/http/`
6. **Tests**: Write unit and integration tests

## License

MIT

## Contributing

This is part of the LeBonPoint project. Please follow the existing code style and patterns when contributing.
