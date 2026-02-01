## 1. Update Database Configuration

- [ ] 1.1 Create shared database path constant in `src/shared/database.ts`
- [ ] 1.2 Export `DATABASE_PATH` as `../database/db.sqlite`
- [ ] 1.3 Add helper function to ensure parent database directory exists

## 2. Update Database Scripts

- [ ] 2.1 Update `scripts/db-init.ts` to use new path
- [ ] 2.2 Add directory creation logic for `../database/` if it doesn't exist
- [ ] 2.3 Add migration logic to copy old database if it exists
- [ ] 2.4 Update `scripts/db-seed.ts` to use new path
- [ ] 2.5 Update `scripts/db-verify.ts` to use new path
- [ ] 2.6 Update verification messages to reference correct location

## 3. Update TypeScript Server Code

- [ ] 3.1 Update `src/infrastructure/persistence/sqlite.ts` to use new path
- [ ] 3.2 Update `src/infrastructure/persistence/item.repository.impl.ts` if it has hardcoded paths
- [ ] 3.3 Replace any hardcoded path references with shared constant import
- [ ] 3.4 Verify all files import from `src/shared/database.ts`

## 4. Update Tests

- [ ] 4.1 Update any test fixtures that reference the old path
- [ ] 4.2 Ensure in-memory database tests still work correctly
- [ ] 4.3 Update integration tests to use new path (or in-memory)
- [ ] 4.4 Run all tests to verify nothing broke

## 5. Update Package.json Scripts

- [ ] 5.1 Review npm scripts in package.json
- [ ] 5.2 Update any scripts that reference the old database path
- [ ] 5.3 Verify scripts work with new path

## 6. Verification

- [ ] 6.1 Run `npm run db:init` - creates database at new location
- [ ] 6.2 Run `npm run db:seed` - seeds data to new location
- [ ] 6.3 Run `npm run db:verify` - verifies database at new location
- [ ] 6.4 Start the TypeScript server - connects to database at new location
- [ ] 6.5 Test API endpoints - all work correctly with new path
- [ ] 6.6 Run full test suite - all tests pass

## 7. Migration Cleanup (Optional)

- [ ] 7.1 Add warning message if old database file is detected
- [ ] 7.2 Provide instructions for manually removing old database
- [ ] 7.3 Update documentation to reflect new structure

## 8. Documentation

- [ ] 8.1 Update README.md with new database location
- [ ] 8.2 Update any setup/development guides
- [ ] 8.3 Document the directory structure change
