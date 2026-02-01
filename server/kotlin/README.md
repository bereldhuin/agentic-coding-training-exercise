# LeBonPoint Kotlin Server

Kotlin implementation of the LeBonPoint Marketplace API using Ktor framework.

## Features

- RESTful API implementing the canonical OpenAPI specification
- Hexagonal architecture (ports and adapters)
- SQLite database with connection pooling
- Full-text search using FTS5
- Cursor-based pagination
- Dependency injection with Koin
- Comprehensive test coverage

## Prerequisites

- Kotlin 1.9+
- JDK 17+
- Gradle 8+

## Installation

```bash
# Install dependencies
./gradlew build
```

## Running the Server

```bash
# Run with default configuration (port 8080)
./gradlew run

# Run with custom port
PORT=9000 ./gradlew run

# Run with custom database path
DATABASE_PATH=/path/to/db.sqlite ./gradlew run
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `DATABASE_PATH` - Path to SQLite database (default: ../database/db.sqlite)
- `ENV` - Environment: development, testing, production (default: development)

## Running Tests

```bash
# Run all tests
./gradlew test

# Run with coverage
./gradlew test jacocoTestReport
```

## API Endpoints

- `GET /health` - Health check
- `GET /v1/items` - List items with filtering, sorting, and pagination
- `POST /v1/items` - Create a new item
- `GET /v1/items/{id}` - Get item by ID
- `PUT /v1/items/{id}` - Replace an item
- `PATCH /v1/items/{id}` - Partially update an item
- `DELETE /v1/items/{id}` - Delete an item

## Architecture

The project follows hexagonal architecture with clear separation of concerns:

```
src/main/kotlin/com/lebonpoint/
├── domain/              # Business logic (framework-agnostic)
│   ├── entities/        # Item entity, value objects
│   ├── repositories/    # Repository interfaces (ports)
│   └── valueobjects/    # ItemImage, enums
├── application/         # Use cases and services
│   ├── usecases/       # CreateItem, GetItem, ListItems, etc.
│   └── services/       # Validation service
├── infrastructure/      # External adapters
│   ├── persistence/    # SQLite repository implementation
│   ├── http/           # Ktor routes and models
│   └── di/             # Koin dependency injection
└── shared/             # Shared utilities
    ├── Configuration.kt
    └── Errors.kt
```

## Database

The server uses SQLite database shared with other implementations. The database schema includes:

- `items` table with indexes
- `items_fts` FTS5 virtual table for full-text search

## OpenAPI Compliance

This server implements the canonical OpenAPI specification at `../api/openapi.yaml`.

All API responses are validated against the specification to ensure compatibility with TypeScript, Python, and Swift implementations.
