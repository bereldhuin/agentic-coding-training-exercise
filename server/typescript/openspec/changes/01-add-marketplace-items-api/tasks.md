## 1. Project Setup

- [ ] 1.1 Initialize npm project with `package.json` (TypeScript config)
- [ ] 1.2 Configure TypeScript (tsconfig.json) with strict mode
- [ ] 1.3 Install dependencies: express, zod, better-sqlite3
- [ ] 1.4 Install dev dependencies: vitest, @types/express, @types/better-sqlite3, tsx

## 2. Domain Layer

- [ ] 2.1 Create domain entities (Item, ItemImage value objects, Location value object)
- [ ] 2.2 Define enums for Condition and Status
- [ ] 2.3 Create repository port interface (IItemRepository)
- [ ] 2.4 Add unit tests for domain entities

## 3. Application Layer

- [ ] 3.1 Create validation service using Zod schemas
- [ ] 3.2 Implement CreateItem use case
- [ ] 3.3 Implement GetItem use case
- [ ] 3.4 Implement ListItems use case (with filters, pagination, FTS)
- [ ] 3.5 Implement UpdateItem use case (PUT)
- [ ] 3.6 Implement PatchItem use case (PATCH)
- [ ] 3.7 Implement DeleteItem use case
- [ ] 3.8 Add unit tests for all use cases with mocked repositories

## 4. Infrastructure Layer - Persistence

- [ ] 4.1 Create SQLite database connection module
- [ ] 4.2 Create database schema (items table, items_fts FTS5 table)
- [ ] 4.3 Implement SQLite repository adapter (IItemRepository implementation)
- [ ] 4.4 Implement FTS search logic with BM25 ranking
- [ ] 4.5 Implement cursor-based pagination
- [ ] 4.6 Add integration tests for repository with in-memory SQLite

## 5. Infrastructure Layer - HTTP

- [ ] 5.1 Create request/response DTOs
- [ ] 5.2 Implement Express routes for all endpoints
- [ ] 5.3 Create error handling middleware
- [ ] 5.4 Implement server startup module
- [ ] 5.5 Add API integration tests (using supertest)

## 6. Configuration and Scripts

- [ ] 6.1 Add npm scripts: dev, build, test, test:watch
- [ ] 6.2 Create environment configuration
- [ ] 6.3 Add database initialization script/migration

## 7. Server Testing with Skills

- [ ] 7.1 Start the server using the pm2-live-testing skill
- [ ] 7.2 Run through API verification checklist:
  - [ ] GET /v1/items (list with default sorting)
  - [ ] GET /v1/items?q=search (full-text search)
  - [ ] GET /v1/items?status=active (filtering)
  - [ ] GET /v1/items?sort=price_cents:asc (sorting)
  - [ ] GET /v1/items?limit=10 (pagination)
  - [ ] POST /v1/items (create)
  - [ ] GET /v1/items/{id} (get single)
  - [ ] PUT /v1/items/{id} (replace)
  - [ ] PATCH /v1/items/{id} (partial update)
  - [ ] DELETE /v1/items/{id} (delete)
- [ ] 7.3 Verify error responses (400, 404)
- [ ] 7.4 Clean up server process

## 8. Documentation

- [ ] 8.1 Add README with setup instructions
- [ ] 8.2 Document API endpoints (or reference OpenAPI spec)
- [ ] 8.3 Document hexagonal architecture structure
