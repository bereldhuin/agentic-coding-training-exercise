# LeBonPoint Marketplace Server (TypeScript)

TypeScript REST API server with SQLite database for the LeBonPoint marketplace application.

## API Specification

**This server implements the canonical OpenAPI specification located at `../api/openapi.yaml`.**

The OpenAPI specification is the single source of truth for the API contract. All endpoints, schemas, and validation rules defined in this implementation conform to the canonical specification.

## Features

- **RESTful API** for marketplace items management
- **Full-text search** with accent-insensitive matching (FTS5)
- **Filtering** by status, category, price range, location, and more
- **Sorting** by any field (created_at, price_cents, etc.)
- **Cursor-based pagination** for efficient data retrieval
- **Hexagonal architecture** (ports and adapters pattern)
- **Comprehensive validation** using Zod schemas
- **Type-safe** with TypeScript strict mode

## Architecture

The project follows hexagonal architecture (ports and adapters) with clear separation of concerns:

```
src/
├── domain/                 # Business logic (no external dependencies)
│   ├── entities/          # Domain entities and enums
│   ├── repositories/      # Repository port interfaces
│   └── value-objects/     # Value objects
├── application/           # Use cases and application services
│   ├── use-cases/        # Business use cases
│   └── services/         # Application services (validation)
├── infrastructure/        # External concerns
│   ├── persistence/      # SQLite repository adapter
│   └── http/             # Express routes and server
└── shared/               # Shared utilities and types
```

## Database

The database is located at `database/marketplace.db` (created by spec 00).

### Database Schema

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique identifier |
| title | TEXT | NOT NULL, CHECK(3-200 chars) | Item title |
| description | TEXT | | Item description |
| price_cents | INTEGER | NOT NULL, CHECK(>= 0) | Price in cents |
| category | TEXT | | Item category |
| condition | TEXT | NOT NULL, CHECK(enum) | Item condition (new, like_new, good, fair, parts, unknown) |
| status | TEXT | NOT NULL, CHECK(enum), DEFAULT 'draft' | Item status (draft, active, reserved, sold, archived) |
| is_featured | INTEGER | NOT NULL, DEFAULT 0 | Featured flag (0/1) |
| city | TEXT | | City name |
| postal_code | TEXT | | Postal code |
| country | TEXT | NOT NULL, DEFAULT 'FR' | Country code |
| delivery_available | INTEGER | NOT NULL, DEFAULT 0 | Delivery available flag (0/1) |
| created_at | TEXT | NOT NULL, DEFAULT datetime('now') | Creation timestamp |
| updated_at | TEXT | NOT NULL, DEFAULT datetime('now') | Update timestamp |
| published_at | TEXT | | Publication timestamp |
| images | TEXT | NOT NULL, DEFAULT '[]' | JSON array of image URLs |

### Full-Text Search

The database includes an FTS5 virtual table `items_fts` for full-text search on title and description with accent-insensitive tokenization.

## API Endpoints

### Base URL
```
http://localhost:3000
```

### Endpoints

#### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2026-01-30T17:36:03.653Z"
}
```

#### GET /v1/items
List items with optional filtering, sorting, pagination, and full-text search.

**Query Parameters:**
- `status` - Filter by status (draft, active, reserved, sold, archived)
- `category` - Filter by category
- `min_price_cents` - Minimum price in cents
- `max_price_cents` - Maximum price in cents
- `city` - Filter by city
- `postal_code` - Filter by postal code
- `is_featured` - Filter by featured flag (true/false)
- `delivery_available` - Filter by delivery availability (true/false)
- `sort` - Sort by field:direction (e.g., `price_cents:asc`, `created_at:desc`)
- `limit` - Number of items per page (1-100, default: 20)
- `cursor` - Cursor for pagination (base64-encoded JSON)
- `q` - Full-text search query

**Response:**
```json
{
  "items": [
    {
      "id": 1,
      "title": "Item Title",
      "description": "Item description",
      "price_cents": 10000,
      "category": "Category",
      "condition": "good",
      "status": "active",
      "is_featured": false,
      "city": "Strasbourg",
      "postal_code": "67000",
      "country": "FR",
      "delivery_available": true,
      "created_at": "2026-01-30T17:36:03.653Z",
      "updated_at": "2026-01-30T17:36:03.653Z",
      "published_at": "2026-01-30T17:36:03.653Z",
      "images": [
        {
          "url": "https://example.com/image.jpg",
          "alt": "Image description",
          "sort_order": 1
        }
      ]
    }
  ],
  "next_cursor": "eyJpZCI6MjB9"
}
```

#### POST /v1/items
Create a new item.

**Request Body:**
```json
{
  "title": "Item Title",
  "description": "Item description",
  "price_cents": 10000,
  "category": "Category",
  "condition": "good",
  "status": "draft",
  "is_featured": false,
  "city": "Strasbourg",
  "postal_code": "67000",
  "country": "FR",
  "delivery_available": true,
  "images": [
    {
      "url": "https://example.com/image.jpg",
      "alt": "Image description"
    }
  ]
}
```

**Response:** `201 Created` with the created item.

#### GET /v1/items/:id
Get a single item by ID.

**Response:** `200 OK` with the item or `404 Not Found`.

#### PUT /v1/items/:id
Replace an entire item (all fields except id are required).

**Request Body:** Same as POST (all required fields).

**Response:** `200 OK` with the updated item or `404 Not Found`.

#### PATCH /v1/items/:id
Partially update an item (only provide fields to update).

**Request Body:**
```json
{
  "title": "New Title"
}
```

**Response:** `200 OK` with the updated item or `404 Not Found`.

#### DELETE /v1/items/:id
Delete an item.

**Response:** `204 No Content` or `404 Not Found`.

### Error Responses

All errors follow this format:

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
- `validation_error` - Invalid input data (400)
- `not_found` - Resource not found (404)
- `internal_error` - Server error (500)

## Development

### Prerequisites

- Node.js 18+
- npm

### Installation

```bash
npm install
```

### Build

```bash
npm run build
```

### Run Development Server

```bash
npm run dev
```

The server will start on port 3000 (or the port specified in PORT environment variable).

### Run Tests

```bash
npm test
```

Run tests in watch mode:

```bash
npm run test:watch
```

## Database Scripts

### Initialize Database

```bash
npm run db:init
```

### Seed Database

```bash
npm run db:seed
```

### Verify Database

```bash
npm run db:verify
```

## Technology Stack

- **Runtime**: Node.js with TypeScript
- **Web Framework**: Express
- **Database**: SQLite with better-sqlite3
- **Validation**: Zod (schemas match OpenAPI specification)
- **Testing**: Vitest + Supertest
- **Executor**: tsx (TypeScript execution)

## API-First Development

This project follows an API-first development approach:

1. The canonical OpenAPI specification is at `../api/openapi.yaml`
2. All API changes should start with updating the specification
3. Server implementations are validated against the canonical spec
4. See `../api/README.md` for the complete API development workflow

## Environment Variables

- `PORT` - Server port (default: 3000)
- `NODE_ENV` - Environment (set to `test` for testing)

## Connecting to Database

Using SQLite CLI:

```bash
sqlite3 database/marketplace.db
```

Using Node.js:

```typescript
import Database from 'better-sqlite3';
const db = new Database('database/marketplace.db');
const items = db.prepare('SELECT * FROM items').all();
```
