## 1. Project Setup

- [ ] 1.1 Initialize npm project with `package.json`
- [ ] 1.2 Install dependencies: `better-sqlite3`
- [ ] 1.3 Install dev dependencies: `tsx`, `typescript`, `@types/better-sqlite3`
- [ ] 1.4 Configure TypeScript (tsconfig.json)
- [ ] 1.5 Add npm scripts: `db:init`, `db:seed`, `db:verify`

## 2. Database Schema

- [ ] 2.1 Create `scripts/db-init.ts` script
- [ ] 2.2 Define items table schema with all columns
- [ ] 2.3 Add CHECK constraints (title length, price_cents >= 0, condition enum, status enum)
- [ ] 2.4 Add DEFAULT values (status='draft', country='FR', is_featured=0, etc.)
- [ ] 2.5 Create items_fts FTS5 virtual table with external content
- [ ] 2.6 Configure FTS5 with unicode61 tokenizer and remove_diacritics

## 3. Seed Data Script

- [ ] 3.1 Create `scripts/db-seed.ts` script
- [ ] 3.2 Define arrays of sample data (categories, cities, conditions, statuses)
- [ ] 3.3 Implement random item generator function
- [ ] 3.4 Generate 15 random items with realistic data
- [ ] 3.5 Insert generated items into database
- [ ] 3.6 Ensure FTS5 index is populated after seeding

## 4. Verification Script

- [ ] 4.1 Create `scripts/db-verify.ts` script
- [ ] 4.2 Implement database existence check
- [ ] 4.3 Implement column presence verification
- [ ] 4.4 Implement constraint verification
- [ ] 4.5 Implement FTS5 table verification
- [ ] 4.6 Implement seed data count verification
- [ ] 4.7 Implement seed data quality checks
- [ ] 4.8 Output clear pass/fail messages

## 5. Testing

- [ ] 5.1 Test database initialization on empty state
- [ ] 5.2 Test database reinitialization (overwrites existing)
- [ ] 5.3 Test seed script creates exactly 15 items
- [ ] 5.4 Test verification script passes on good database
- [ ] 5.5 Test verification script fails on missing/invalid database
- [ ] 5.6 Test connection using sqlite3 CLI
- [ ] 5.7 Test connection from Node.js script

## 6. Documentation

- [ ] 6.1 Add README.md with setup and usage instructions
- [ ] 6.2 Document schema in a diagram or table
- [ ] 6.3 Document seed data format and how to customize
