## MODIFIED Requirements

### Requirement: Database File Creation

The system SHALL create a SQLite database file at `../database/db.sqlite` (relative to the typescript service directory, resolving to `/server/database/db.sqlite`) when the init script is executed.

#### Scenario: Initialize database at parent directory location
- **GIVEN** no database file exists at `../database/db.sqlite`
- **WHEN** running `npm run db:init`
- **THEN** create `../database/db.sqlite` with proper schema
- **AND** create the `../database/` directory if it doesn't exist

#### Scenario: Migrate existing database
- **GIVEN** a database file exists at the old location `database/marketplace.db`
- **WHEN** running `npm run db:init` after this change
- **THEN** copy the existing database to `../database/db.sqlite`
- **AND** log a message informing about the migration
- **AND** prompt the user to manually delete the old file if desired

#### Scenario: Reinitialize existing database at new location
- **GIVEN** a database file already exists at `../database/db.sqlite`
- **WHEN** running `npm run db:init` again
- **THEN** overwrite or drop and recreate the database schema at the new location

### Requirement: Database Connection Path Resolution

The system SHALL resolve the database path correctly from scripts and runtime code.

#### Scenario: Script path resolution
- **GIVEN** a database script is run from within the typescript directory
- **WHEN** the script accesses the database
- **THEN** the path resolves to `../database/db.sqlite` relative to the script location
- **AND** the connection succeeds

#### Scenario: Runtime code path resolution
- **GIVEN** the TypeScript server is running
- **WHEN** the server code connects to the database
- **THEN** the path resolves to `../database/db.sqlite` relative to the compiled output location
- **AND** the connection succeeds

#### Scenario: Path resolution from different working directories
- **GIVEN** the database script is run from a different working directory
- **WHEN** the script attempts to connect
- **THEN** the path still resolves correctly using `__dirname` based resolution
- **AND** the connection succeeds

### Requirement: Updated Script Paths

All database scripts (db-init, db-seed, db-verify) SHALL reference the new database location.

#### Scenario: db-init uses new path
- **GIVEN** the db-init.ts script exists
- **WHEN** examining the database path constant
- **THEN** it references `../database/db.sqlite`

#### Scenario: db-seed uses new path
- **GIVEN** the db-seed.ts script exists
- **WHEN** examining the database path constant
- **THEN** it references `../database/db.sqlite`

#### Scenario: db-verify uses new path
- **GIVEN** the db-verify.ts script exists
- **WHEN** examining the database path constant
- **THEN** it references `../database/db.sqlite`

### Requirement: Server Code Path Updates

The TypeScript server code SHALL use the updated database path in all connection and repository code.

#### Scenario: Repository connection uses new path
- **GIVEN** the SQLite repository implementation exists
- **WHEN** examining the database connection code
- **THEN** it references `../database/db.sqlite`

#### Scenario: No hardcoded paths exist
- **GIVEN** the TypeScript codebase exists
- **WHEN** searching for references to the old path `database/marketplace.db`
- **THEN** no results are found

#### Scenario: Shared path configuration
- **GIVEN** multiple files need the database path
- **WHEN** examining the code structure
- **THEN** a shared constant or configuration module exports the path
- **AND** all files import from this shared source

### Requirement: Verification Compatibility

The verification script SHALL validate the database at the new location.

#### Scenario: Verify database at new location
- **GIVEN** the database exists at `../database/db.sqlite`
- **WHEN** running `npm run db:verify`
- **THEN** the script finds and validates the database at the new location

#### Scenario: Detect missing database at new location
- **GIVEN** no database file exists at `../database/db.sqlite`
- **WHEN** running `npm run db:verify`
- **THEN** the script exits with an error message indicating the database is missing at the expected location

### Requirement: Test Compatibility

All existing tests SHALL continue to pass after the path migration.

#### Scenario: Unit tests with mocked paths
- **GIVEN** unit tests exist for use cases
- **WHEN** running the tests
- **THEN** all tests pass with the updated path references

#### Scenario: Integration tests with real database
- **GIVEN** integration tests use a test database
- **WHEN** running the tests
- **THEN** tests use either in-memory database or the new path
- **AND** all tests pass

#### Scenario: Repository tests with in-memory database
- **GIVEN** repository tests use in-memory SQLite
- **WHEN** running the tests
- **THEN** tests continue to use in-memory database for isolation
- **AND** all tests pass

### Requirement: Documentation Updates

Documentation SHALL reflect the new database location.

#### Scenario: README mentions correct path
- **GIVEN** the project has a README
- **WHEN** reading the database setup section
- **THEN** it references `database/db.sqlite` at the server level

#### Scenario: Code comments updated
- **GIVEN** code comments reference the database location
- **WHEN** reviewing the comments
- **THEN** they reference the new path `../database/db.sqlite` or the shared constant
