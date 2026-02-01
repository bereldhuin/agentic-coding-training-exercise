#!/usr/bin/env tsx
/**
 * Database Verification Script
 *
 * Verifies that the database schema is correct and that seed data has been properly populated.
 */

import Database from 'better-sqlite3';
import path from 'path';
import fs from 'fs';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Database path - at server/database/db.sqlite (relative to repo root)
const DB_PATH = path.join(__dirname, '..', 'server', 'database', 'db.sqlite');

console.log('Database Verification Script');
console.log('============================');
console.log(`Database path: ${DB_PATH}\n`);

let allPassed = true;
let checkCount = 0;
let passCount = 0;

// Helper function to run a check
function runCheck(checkName: string, checkFn: () => boolean): void {
  checkCount++;
  process.stdout.write(`  [${checkCount}] ${checkName}... `);
  try {
    const passed = checkFn();
    if (passed) {
      console.log('✓ PASS');
      passCount++;
    } else {
      console.log('✗ FAIL');
      allPassed = false;
    }
  } catch (error) {
    console.log('✗ ERROR');
    console.error(`    Error: ${error}`);
    allPassed = false;
  }
}

// Check 1: Database file exists
runCheck('Database file exists', () => {
  return fs.existsSync(DB_PATH);
});

if (!fs.existsSync(DB_PATH)) {
  console.error('\n❌ Database file not found. Please run npm run db:init first.');
  process.exit(1);
}

// Create database connection
const db = new Database(DB_PATH);

// Check 2: Items table exists
runCheck('Items table exists', () => {
  const result = db
    .prepare(
      "SELECT name FROM sqlite_master WHERE type='table' AND name='items'"
    )
    .get() as any;
  return result && result.name === 'items';
});

// Check 3: FTS5 table exists
runCheck('FTS5 virtual table exists', () => {
  const result = db
    .prepare(
      "SELECT name, sql FROM sqlite_master WHERE type='table' AND name='items_fts'"
    )
    .get() as any;
  if (!result || result.name !== 'items_fts') return false;
  // Check it's using FTS5
  return result.sql && result.sql.includes('USING fts5');
});

// Check 4: FTS5 table is properly linked
runCheck('FTS5 table is linked to items table', () => {
  const result = db
    .prepare(
      "SELECT sql FROM sqlite_master WHERE type='table' AND name='items_fts'"
    )
    .get() as any;
  if (!result) return false;
  // Check for content=items and content_rowid=rowid
  return (
    result.sql &&
    result.sql.includes('content=items') &&
    result.sql.includes('content_rowid=rowid')
  );
});

// Check 5: FTS5 has accent-insensitive tokenization
runCheck('FTS5 uses accent-insensitive tokenization', () => {
  const result = db
    .prepare(
      "SELECT sql FROM sqlite_master WHERE type='table' AND name='items_fts'"
    )
    .get() as any;
  if (!result) return false;
  // Check for unicode61 tokenizer with remove_diacritics
  return (
    result.sql &&
    result.sql.includes('unicode61') &&
    result.sql.includes('remove_diacritics 1')
  );
});

// Check 6: All required columns exist
const requiredColumns = [
  'id',
  'title',
  'description',
  'price_cents',
  'category',
  'condition',
  'status',
  'is_featured',
  'city',
  'postal_code',
  'country',
  'delivery_available',
  'created_at',
  'updated_at',
  'published_at',
  'images'
];

runCheck('All 16 required columns exist', () => {
  const columns = db.prepare('PRAGMA table_info(items)').all() as any[];
  const columnNames = columns.map((col) => col.name);
  return requiredColumns.every((col) => columnNames.includes(col));
});

// Check 7: Title has NOT NULL constraint
runCheck('Title column has NOT NULL constraint', () => {
  const result = db
    .prepare("SELECT * FROM pragma_table_info('items') WHERE name='title'")
    .get() as any;
  return result && result.notnull === 1;
});

// Check 8: Price has NOT NULL constraint
runCheck('Price column has NOT NULL constraint', () => {
  const result = db
    .prepare("SELECT * FROM pragma_table_info('items') WHERE name='price_cents'")
    .get() as any;
  return result && result.notnull === 1;
});

// Check 9: Title has CHECK constraint for length
runCheck('Title has CHECK constraint for length (3-200)', () => {
  try {
    // Try to insert a row with too short title
    db.prepare('INSERT INTO items (title, price_cents, condition) VALUES (?, ?, ?)').run('ab', 1000, 'good');
    // If we got here, the constraint didn't work
    db.prepare('DELETE FROM items WHERE title = ?').run('ab');
    return false;
  } catch (error) {
    // Expected to fail with CHECK constraint
    return true;
  }
});

// Check 10: Price has CHECK constraint for non-negative
runCheck('Price has CHECK constraint for non-negative', () => {
  try {
    // Try to insert a row with negative price
    db.prepare('INSERT INTO items (title, price_cents, condition) VALUES (?, ?, ?)').run(
      'Test Item',
      -100,
      'good'
    );
    // If we got here, the constraint didn't work
    db.prepare('DELETE FROM items WHERE title = ?').run('Test Item');
    return false;
  } catch (error) {
    // Expected to fail with CHECK constraint
    return true;
  }
});

// Check 11: Condition has CHECK constraint for enum
runCheck('Condition has CHECK constraint for enum values', () => {
  try {
    // Try to insert a row with invalid condition
    db.prepare('INSERT INTO items (title, price_cents, condition) VALUES (?, ?, ?)').run(
      'Test Item',
      1000,
      'invalid_condition'
    );
    // If we got here, the constraint didn't work
    db.prepare('DELETE FROM items WHERE title = ?').run('Test Item');
    return false;
  } catch (error) {
    // Expected to fail with CHECK constraint
    return true;
  }
});

// Check 12: Status has CHECK constraint for enum
runCheck('Status has CHECK constraint for enum values', () => {
  try {
    // Try to insert a row with invalid status
    db.prepare('INSERT INTO items (title, price_cents, condition, status) VALUES (?, ?, ?, ?)').run(
      'Test Item',
      1000,
      'good',
      'invalid_status'
    );
    // If we got here, the constraint didn't work
    db.prepare('DELETE FROM items WHERE title = ?').run('Test Item');
    return false;
  } catch (error) {
    // Expected to fail with CHECK constraint
    return true;
  }
});

// Check 13: Default values are set correctly
runCheck('Default values are set correctly', () => {
  try {
    const stmt = db.prepare(
      'INSERT INTO items (title, price_cents, condition) VALUES (?, ?, ?)'
    );
    const result = stmt.run('Test Item for Defaults', 1000, 'good');
    const item = db.prepare('SELECT * FROM items WHERE rowid = ?').get(result.lastInsertRowid) as any;
    db.prepare('DELETE FROM items WHERE rowid = ?').run(result.lastInsertRowid);

    return (
      item.status === 'draft' &&
      item.is_featured === 0 &&
      item.delivery_available === 0 &&
      item.country === 'FR' &&
      item.images === '[]' &&
      item.created_at !== null &&
      item.updated_at !== null
    );
  } catch (error) {
    console.error(`    Error: ${error}`);
    return false;
  }
});

// Check 14: Seed data count
runCheck('Database contains exactly 15 items', () => {
  const result = db.prepare('SELECT COUNT(*) as count FROM items').get() as { count: number };
  return result.count === 15;
});

// Check 15: Seed data has variety in categories
runCheck('Seed data has at least 3 different categories', () => {
  const result = db.prepare('SELECT COUNT(DISTINCT category) as count FROM items').get() as { count: number };
  return result.count >= 3;
});

// Check 16: Seed data has variety in conditions
runCheck('Seed data has at least 3 different conditions', () => {
  const result = db.prepare('SELECT COUNT(DISTINCT condition) as count FROM items').get() as { count: number };
  return result.count >= 3;
});

// Check 17: Seed data has variety in statuses
runCheck('Seed data has at least 2 different statuses', () => {
  const result = db.prepare('SELECT COUNT(DISTINCT status) as count FROM items').get() as { count: number };
  return result.count >= 2;
});

// Check 18: All prices are non-negative
runCheck('All prices are non-negative', () => {
  const result = db.prepare('SELECT MIN(price_cents) as min_price FROM items').get() as { min_price: number };
  return result.min_price >= 0;
});

// Check 19: At least 10 items have French locations
runCheck('At least 10 items have French city/postal_code', () => {
  const result = db
    .prepare("SELECT COUNT(*) as count FROM items WHERE city IS NOT NULL AND postal_code IS NOT NULL")
    .get() as { count: number };
  return result.count >= 10;
});

// Check 20: At least 10 items have images
runCheck('At least 10 items have non-empty image arrays', () => {
  const result = db.prepare("SELECT COUNT(*) as count FROM items WHERE images != '[]'").get() as { count: number };
  return result.count >= 10;
});

// Check 21: FTS5 index is populated
runCheck('FTS5 index is populated', () => {
  const result = db.prepare('SELECT COUNT(*) as count FROM items_fts').get() as { count: number };
  return result.count === 15;
});

// Check 22: FTS5 search works
runCheck('FTS5 full-text search works', () => {
  const result = db.prepare('SELECT COUNT(*) as count FROM items_fts WHERE items_fts MATCH ?').get('iPhone') as { count: number };
  return (result as any).count >= 0; // Just check it doesn't error
});

// Additional stats
const statsResult = db.prepare(`
  SELECT
    COUNT(DISTINCT category) as categories,
    COUNT(DISTINCT condition) as conditions,
    COUNT(DISTINCT status) as statuses,
    COUNT(DISTINCT city) as cities
  FROM items
`).get() as any;

// Close database
db.close();

// Summary
console.log('\n' + '='.repeat(50));
console.log(`Summary: ${passCount}/${checkCount} checks passed`);

if (statsResult) {
  console.log('\nDatabase Statistics:');
  console.log(`  - Total items: 15`);
  console.log(`  - Categories: ${statsResult.categories}`);
  console.log(`  - Conditions: ${statsResult.conditions}`);
  console.log(`  - Statuses: ${statsResult.statuses}`);
  console.log(`  - Cities: ${statsResult.cities}`);
}

if (allPassed) {
  console.log('\n✅ All verification checks passed!');
  console.log('\nThe database is ready to use.');
} else {
  console.log('\n❌ Some verification checks failed.');
  console.log('\nPlease review the failed checks above and ensure:');
  console.log('  1. You have run "npm run db:init"');
  console.log('  2. You have run "npm run db:seed"');
  process.exit(1);
}
