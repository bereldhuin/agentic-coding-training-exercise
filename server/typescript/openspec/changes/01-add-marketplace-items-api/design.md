## Context

This is a greenfield project implementing a marketplace items API. The system needs to support CRUD operations, filtering, pagination, and full-text search for marketplace listings.

## Goals / Non-Goals

- **Goals**:
  - Working Express server with TypeScript
  - Hexagonal architecture (ports/adapters) for clean separation of concerns
  - SQLite storage with FTS5 full-text search
  - RESTful API at `/v1/items`
  - Comprehensive test coverage
  - Run and verify the server using skills

- **Non-Goals**:
  - Authentication/authorization (out of scope for this change)
  - Image upload/handling (URLs only for now)
  - User management
  - Real-time features

## Decisions

### Architecture: Hexagonal (Ports and Adapters)

**Rationale**: Separates business logic from infrastructure, making the codebase testable and maintainable.

**Layers**:
- `domain/`: Business entities (Item, ItemImage, value objects)
- `application/`: Use cases/orchestrators (CreateItem, ListItems, etc.)
- `infrastructure/`: External concerns (Express routes, SQLite repository)
- `interfaces/`: Ports (repository interface) and adapters (Express adapter, SQLite adapter)

**Structure**:
```
src/
├── domain/
│   ├── entities/
│   │   └── item.ts
│   └── repositories/
│       └── item.repository.ts (port interface)
├── application/
│   ├── use-cases/
│   │   ├── create-item.use-case.ts
│   │   ├── list-items.use-case.ts
│   │   ├── get-item.use-case.ts
│   │   ├── update-item.use-case.ts
│   │   └── delete-item.use-case.ts
│   └── services/
│       └── validation.service.ts
├── infrastructure/
│   ├── persistence/
│   │   ├── sqlite.ts
│   │   └── item.repository.impl.ts (adapter)
│   └── http/
│       ├── routes.ts
│       └── server.ts
└── shared/
    └── types.ts
```

### Database: SQLite with FTS5

**Rationale**: Lightweight, embedded, no separate server needed, excellent FTS5 support.

**Schema**:
- `items`: Main table with all fields
- `items_fts`: FTS5 virtual table (external content) for title/description search
- Images stored as JSON in `images` column

### Pagination: Cursor-based

**Rationale**: More efficient than offset for large datasets, consistent results.

**Implementation**: Base64-encoded JSON cursor with `id` + `created_at` for tie-breaking.

### Validation: Zod schemas

**Rationale**: Runtime type validation, TypeScript integration, clear error messages.

### Testing: Vitest

**Rationale**: Fast, native ESM support, similar API to Jest.

## Risks / Trade-offs

- **Hexagonal complexity**: May feel over-engineered for a simple API
  - *Mitigation*: Keep it minimal, avoid premature abstractions
- **SQLite limitations**: Single write connection, not suitable for high concurrency
  - *Acceptable*: For MVP/single-server deployment
- **Cursor pagination complexity**: More complex than offset
  - *Mitigation*: Provide clear examples in docs
- **JSON storage for images**: Not queryable, denormalized
  - *Acceptable*: Simple for MVP, can move to separate table later

## Migration Plan

N/A (greenfield project)

## Open Questions

None identified. All requirements are clear from the spec.
