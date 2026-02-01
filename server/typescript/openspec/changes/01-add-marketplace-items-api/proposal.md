# Change: Add Marketplace Items API

## Why

The application needs a backend API to manage marketplace listings (items) with full CRUD operations, filtering, pagination, and full-text search capabilities. This is the core feature for the marketplace platform.

## What Changes

- **NEW**: Create Express server with TypeScript following hexagonal architecture
- **NEW**: SQLite database with items table and FTS5 full-text search
- **NEW**: REST API endpoints for items management at `/v1/items`
- **NEW**: Support for filtering, sorting, pagination, and full-text search
- **NEW**: Comprehensive test suite (unit + integration)

## Impact

- Affected specs: New capability `items-api`
- Affected code: New project structure under `src/`
- New dependencies: express, sqlite3, better-sqlite3 (or similar), zod (validation)
- Architecture: Hexagonal (ports and adapters) pattern with clear separation of concerns
