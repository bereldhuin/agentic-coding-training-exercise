# Change: Create Marketplace Database

## Why

The marketplace application needs a SQLite database with the proper schema to store items. Before building the API layer, we need to establish the data layer with a working database, schema, and seed data for development/testing.

## What Changes

- **NEW**: SQLite database initialization with items schema
- **NEW**: Database migration/creation script
- **NEW**: Seed script that creates 15 random items with various data
- **NEW**: Verification script that validates all fields are correctly implemented
- **NEW**: npm scripts for database operations (init, seed, verify)

## Impact

- Affected specs: New capability `database`
- Creates: `database/` directory with SQLite file and scripts
- Dependencies: better-sqlite3 or sqlite3
- Prerequisite for: `01-add-marketplace-items-api` (the API layer depends on this database)
