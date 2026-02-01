## Context

After implementing `00-create-marketplace-database`, the database resides at `database/marketplace.db` within the typescript service directory. To support a multi-service architecture where the database might be accessed by multiple services, we need to relocate it to a shared parent directory.

## Goals / Non-Goals

- **Goals**:
  - Move database file from `database/marketplace.db` to `../database/db.sqlite`
  - Update all TypeScript code to reference the new path
  - Update all scripts (init, seed, verify) to use the new path
  - Ensure all tests still pass after the migration
  - Create parent `database/` directory if it doesn't exist

- **Non-Goals**:
  - Changing the database schema
  - Modifying the seed data
  - Adding new features

## Decisions

### New Database Location: `../database/db.sqlite`

**Rationale**:
- Parent directory allows database to be shared across services
- `db.sqlite` is more generic than `marketplace.db` in case the schema expands
- Relative path `../database/` from typescript service resolves to `/server/database/`

**Final path structure**:
```
server/
├── database/
│   └── db.sqlite          # New location
└── typescript/
    ├── scripts/
    │   ├── db-init.ts     # Updated to use ../database/db.sqlite
    │   ├── db-seed.ts     # Updated to use ../database/db.sqlite
    │   └── db-verify.ts   # Updated to use ../database/db.sqlite
    └── src/
        └── infrastructure/
            └── persistence/
                └── sqlite.ts  # Updated to use ../database/db.sqlite
```

### Path Resolution Strategy

**Option 1**: Relative path `../database/db.sqlite`
- Pros: Works from any location in typescript directory
- Cons: Breaks if script is run from a different working directory

**Option 2**: Absolute path from environment variable
- Pros: Configurable, works from anywhere
- Cons: Requires environment setup

**Decision**: Use relative path with `__dirname` resolution for scripts, and a shared config constant for runtime code.

### Migration Strategy

1. Create `../database/` directory if it doesn't exist
2. Copy existing database to new location (if exists)
3. Update all path references in code
4. Verify connection works
5. Optionally remove old database file

## Risks / Trade-offs

- **Breaking change**: Any tools or scripts expecting the old path will fail
  - *Mitigation*: Document the change clearly, update all references at once
- **Relative path fragility**: Scripts run from different directories may fail
  - *Mitigation*: Use `path.resolve(__dirname, '../database/db.sqlite')` in all scripts
- **Test failures**: Tests may have hardcoded paths
  - *Mitigation*: Update test fixtures and mocks

## Migration Plan

1. Update `00-create-marketplace-database` specs to reflect new location
2. Modify all database scripts to use new path
3. Update TypeScript server connection code
4. Run `npm run db:verify` to confirm everything works
5. Update documentation

## Open Questions

- **Should we create the database directory in this change or expect it to exist?**
  - *Proposal*: Create it in the db-init script if it doesn't exist

- **What if the old database exists with data?**
  - *Proposal*: Copy the data to the new location, then prompt user to delete the old file

- **Should tests use a separate in-memory database or the file-based one?**
  - *Proposal*: Tests should continue using in-memory database for isolation
