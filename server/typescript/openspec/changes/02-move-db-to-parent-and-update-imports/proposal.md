# Change: Move Database to Parent Directory and Update Imports

## Why

The database should be shared at the server level (not isolated within the typescript service) to allow potential future services (e.g., a Python service, Go service) to access the same database. Moving the database to `../database/db.sqlite` establishes a better project structure for multi-service architectures.

## What Changes

- **MODIFIED**: Database location moves from `database/marketplace.db` to `../database/db.sqlite`
- **MODIFIED**: All database path references in scripts must be updated
- **MODIFIED**: All imports and path references in TypeScript server code must be updated
- **MODIFIED**: npm scripts must use the updated path
- **BREAKING**: Any hardcoded paths to the old location will break

## Impact

- Affected specs: `database` (MODIFIED requirements)
- Affected code:
  - `scripts/db-init.ts` - database path
  - `scripts/db-seed.ts` - database path
  - `scripts/db-verify.ts` - database path
  - Any repository/connection code in typescript server
  - package.json npm scripts
- Dependencies: Requires `00-create-marketplace-database` to be completed first
- Must be completed before: `01-add-marketplace-items-api` uses the database
