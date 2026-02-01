## Context

This is the foundational change for the marketplace application. Before implementing the API and business logic, we need to establish the database layer with proper schema and seed data for development and testing.

## Goals / Non-Goals

- **Goals**:
  - Create SQLite database with items table matching the item entity specification
  - Create FTS5 virtual table for full-text search
  - Provide a seed script with 15 realistic random items
  - Provide a verification script to validate schema and data integrity
  - Make it possible to connect to and query the database

- **Non-Goals**:
  - API endpoints (covered in `01-add-marketplace-items-api`)
  - Repository/adapters (covered in `01-add-marketplace-items-api`)
  - Production migration tooling (this is for local development)

## Decisions

### Database: SQLite

**Rationale**: Lightweight, embedded, no separate server process, excellent for development and small-scale deployments.

### Library: better-sqlite3

**Rationale**: Synchronous API (simpler than callback-based sqlite3), fast, well-maintained, good TypeScript support.

### Schema Design

**Items table** (`items`):
```sql
CREATE TABLE items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL CHECK(length(title) >= 3 AND length(title) <= 200),
  description TEXT,
  price_cents INTEGER NOT NULL CHECK(price_cents >= 0),
  category TEXT,
  condition TEXT NOT NULL CHECK(condition IN ('new', 'like_new', 'good', 'fair', 'parts', 'unknown')),
  status TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'active', 'reserved', 'sold', 'archived')),
  is_featured INTEGER NOT NULL DEFAULT 0, -- Stored as 0/1, mapped to boolean in app
  city TEXT,
  postal_code TEXT,
  country TEXT NOT NULL DEFAULT 'FR',
  delivery_available INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at TEXT NOT NULL DEFAULT (datetime('now')),
  published_at TEXT,
  images TEXT NOT NULL DEFAULT '[]' -- JSON array
);
```

**FTS5 virtual table** (`items_fts`):
```sql
CREATE VIRTUAL TABLE items_fts USING fts5(
  title,
  description,
  content=items,
  content_rowid=rowid,
  tokenize='unicode61 remove_diacritics 1'
);
```

Using FTS5 external content table so `items` remains the source of truth.

### Seed Data Strategy

Generate 15 random items with:
- Realistic French marketplace data (locations like Strasbourg, Paris, Lyon, etc.)
- Mix of categories (Ordinateurs, Appareils Photo, Consoles, etc.)
- Various conditions and statuses
- Realistic prices in euros (stored as cents)
- Sample images (placeholder URLs)

### Scripts Structure

```
scripts/
├── db-init.ts      # Create database and schema
├── db-seed.ts      # Seed 15 random items
└── db-verify.ts    # Verify schema and data
```

### Verification Script

The verify script should:
1. Connect to the database
2. Check that all columns exist with correct types
3. Check that constraints (CHECK, DEFAULT) are properly set
4. Verify FTS5 table exists and is linked
5. Validate seed data (count=15, field values in expected ranges)

## Risks / Trade-offs

- **JSON storage for images**: Not queryable in SQL, but simple for MVP
  - *Acceptable*: Can move to separate table later if needed
- **SQLite limitations**: Single write connection
  - *Acceptable*: Sufficient for development and low-traffic scenarios
- **Seed data randomness**: May produce unrealistic combinations
  - *Mitigation*: Use curated lists of realistic values

## Migration Plan

N/A (new database)

## Open Questions

- **Database location**: Root of project or in `database/` subdirectory?
  - *Proposal*: Use `database/` subdirectory for organization
- **Database file name**: `marketplace.db` or `items.db`?
  - *Proposal*: `marketplace.db` for future expansion
