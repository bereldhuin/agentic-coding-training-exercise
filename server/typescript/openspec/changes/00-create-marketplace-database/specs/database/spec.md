## ADDED Requirements

### Requirement: Database File Creation

The system SHALL create a SQLite database file at `database/marketplace.db` when the init script is executed.

#### Scenario: Initialize empty database
- **GIVEN** no database file exists
- **WHEN** running `npm run db:init`
- **THEN** create `database/marketplace.db` with proper schema

#### Scenario: Reinitialize existing database
- **GIVEN** a database file already exists
- **WHEN** running `npm run db:init` again
- **THEN** overwrite or drop and recreate the database schema

### Requirement: Items Table Schema

The system SHALL create an `items` table with all required columns, constraints, and defaults matching the item entity specification.

#### Scenario: Verify table structure
- **GIVEN** the database is initialized
- **WHEN** querying `PRAGMA table_info(items)`
- **THEN** return all columns: id, title, description, price_cents, category, condition, status, is_featured, city, postal_code, country, delivery_available, created_at, updated_at, published_at, images

#### Scenario: Verify NOT NULL constraints
- **GIVEN** the items table exists
- **WHEN** inserting a row without required fields (title, price_cents, condition, status, country)
- **THEN** the insert fails with a constraint violation

#### Scenario: Verify CHECK constraint on title length
- **GIVEN** the items table exists
- **WHEN** inserting a row with title length < 3
- **THEN** the insert fails with a CHECK constraint violation
- **AND** inserting with title length > 200 also fails

#### Scenario: Verify CHECK constraint on price_cents
- **GIVEN** the items table exists
- **WHEN** inserting a row with price_cents < 0
- **THEN** the insert fails with a CHECK constraint violation

#### Scenario: Verify CHECK constraint on condition enum
- **GIVEN** the items table exists
- **WHEN** inserting a row with condition not in ('new', 'like_new', 'good', 'fair', 'parts', 'unknown')
- **THEN** the insert fails with a CHECK constraint violation

#### Scenario: Verify CHECK constraint on status enum
- **GIVEN** the items table exists
- **WHEN** inserting a row with status not in ('draft', 'active', 'reserved', 'sold', 'archived')
- **THEN** the insert fails with a CHECK constraint violation

#### Scenario: Verify DEFAULT values
- **GIVEN** the items table exists
- **WHEN** inserting a row without specifying status, is_featured, delivery_available, country, created_at, updated_at, images
- **THEN** the row is created with status='draft', is_featured=0, delivery_available=0, country='FR', images='[]', and valid timestamps

### Requirement: FTS5 Full-Text Search Table

The system SHALL create an FTS5 virtual table `items_fts` for full-text search on title and description with accent-insensitive tokenization.

#### Scenario: Verify FTS5 table exists
- **GIVEN** the database is initialized
- **WHEN** querying `sqlite_master` for items_fts
- **THEN** the virtual table exists with type='table' and sql containing 'USING fts5'

#### Scenario: Verify FTS5 is external content table
- **GIVEN** the items_fts table exists
- **WHEN** examining the table definition
- **THEN** it references `content=items` and `content_rowid=rowid`

#### Scenario: Verify accent-insensitive tokenization
- **GIVEN** the items_fts table exists
- **WHEN** examining the tokenizer configuration
- **THEN** it uses 'unicode61 remove_diacritics 1' for accent-insensitive search

#### Scenario: Verify FTS5 columns
- **GIVEN** the items_fts table exists
- **WHEN** querying the schema
- **THEN** it indexes the `title` and `description` columns

### Requirement: Seed Data Generation

The system SHALL provide a script that generates 15 random items with realistic marketplace data.

#### Scenario: Seed 15 items
- **GIVEN** the database is initialized
- **WHEN** running `npm run db:seed`
- **THEN** insert exactly 15 items into the items table

#### Scenario: Verify seed data completeness
- **GIVEN** the seed script has run
- **WHEN** querying `SELECT COUNT(*) FROM items`
- **THEN** return 15

#### Scenario: Verify seed data variety - categories
- **GIVEN** the seed script has run
- **WHEN** querying `SELECT DISTINCT category FROM items`
- **THEN** return at least 3 different categories

#### Scenario: Verify seed data variety - conditions
- **GIVEN** the seed script has run
- **WHEN** querying `SELECT DISTINCT condition FROM items`
- **THEN** return at least 3 different condition values

#### Scenario: Verify seed data variety - statuses
- **GIVEN** the seed script has run
- **WHEN** querying `SELECT DISTINCT status FROM items`
- **THEN** return at least 2 different status values

#### Scenario: Verify seed data has valid prices
- **GIVEN** the seed script has run
- **WHEN** querying `SELECT price_cents FROM items`
- **THEN** all values are non-negative integers

#### Scenario: Verify seed data has valid locations
- **GIVEN** the seed script has run
- **WHEN** querying items with city and postal_code
- **THEN** at least 10 items have French city/postal_code combinations

#### Scenario: Verify seed data has images
- **GIVEN** the seed script has run
- **WHEN** querying `SELECT images FROM items`
- **THEN** at least 10 items have non-empty JSON arrays for images

### Requirement: Database Connection

The system SHALL allow connecting to the database using standard SQLite tools.

#### Scenario: Connect using sqlite3 CLI
- **GIVEN** the database file exists
- **WHEN** running `sqlite3 database/marketplace.db`
- **THEN** successfully connect and be able to run SQL queries

#### Scenario: Connect from Node.js
- **GIVEN** the database file exists and better-sqlite3 is installed
- **WHEN** creating a Database instance with the file path
- **THEN** successfully connect without errors

### Requirement: Database Verification Script

The system SHALL provide a verification script that validates schema correctness and data integrity.

#### Scenario: Verify all columns exist
- **GIVEN** the database exists
- **WHEN** running `npm run db:verify`
- **THEN** the script confirms all 16 columns exist in items table

#### Scenario: Verify all constraints exist
- **GIVEN** the database exists
- **WHEN** running `npm run db:verify`
- **THEN** the script confirms CHECK constraints for title, price_cents, condition, and status

#### Scenario: Verify FTS5 table is linked
- **GIVEN** the database exists
- **WHEN** running `npm run db:verify`
- **THEN** the script confirms items_fts exists and is properly linked to items

#### Scenario: Verify seed data count
- **GIVEN** the database has been seeded
- **WHEN** running `npm run db:verify`
- **THEN** the script confirms exactly 15 items exist

#### Scenario: Verify seed data quality
- **GIVEN** the database has been seeded
- **WHEN** running `npm run db:verify`
- **THEN** the script confirms all required fields have valid values (no NULL where not allowed, enums are valid, etc.)

#### Scenario: Detect missing database
- **GIVEN** no database file exists
- **WHEN** running `npm run db:verify`
- **THEN** the script exits with an error message indicating the database is missing

### Requirement: NPM Scripts

The system SHALL provide npm scripts for all database operations.

#### Scenario: Run db:init via npm
- **GIVEN** package.json contains the db:init script
- **WHEN** running `npm run db:init`
- **THEN** the database initialization script executes

#### Scenario: Run db:seed via npm
- **GIVEN** package.json contains the db:seed script
- **WHEN** running `npm run db:seed`
- **THEN** the seed script executes

#### Scenario: Run db:verify via npm
- **GIVEN** package.json contains the db:verify script
- **WHEN** running `npm run db:verify`
- **THEN** the verification script executes

### Requirement: TypeScript Configuration

The database scripts SHALL be written in TypeScript with proper types.

#### Scenario: Scripts compile successfully
- **GIVEN** the TypeScript scripts exist
- **WHEN** compiling the scripts
- **THEN** no TypeScript errors occur

#### Scenario: Scripts use tsx for execution
- **GIVEN** the scripts are written in TypeScript
- **WHEN** running via npm scripts
- **THEN** tsx is used to execute the scripts directly without precompilation
