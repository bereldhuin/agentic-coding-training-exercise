#!/usr/bin/env tsx
/**
 * Database Empty Script
 *
 * Empties all data from the database while preserving the schema.
 * This is useful when you want to clear the database without recreating it.
 */

import Database from 'better-sqlite3';
import path from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';
import fs from 'fs';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Database path - at server/database/db.sqlite (relative to repo root)
const DB_PATH = path.join(__dirname, '..', 'server', 'database', 'db.sqlite');

console.log('Database Empty Script');
console.log('====================');
console.log(`Database path: ${DB_PATH}`);

// Check if database exists
if (!fs.existsSync(DB_PATH)) {
  console.error('❌ Database not found. Please run npm run db:init first.');
  process.exit(1);
}

// Create database connection
const db = new Database(DB_PATH);

// Check current row count
const beforeCount = db.prepare('SELECT COUNT(*) as count FROM items').get() as { count: number };
console.log(`\nCurrent items count: ${beforeCount.count}`);

if (beforeCount.count === 0) {
  console.log('✓ Database is already empty');
  db.close();
  process.exit(0);
}

// Empty the tables
console.log('\nEmptying database...');
db.prepare('DELETE FROM items').run();
db.prepare('DELETE FROM items_fts').run();
console.log('✓ Deleted all items');
console.log('✓ Cleared FTS index');

// Verify empty state
const afterCount = db.prepare('SELECT COUNT(*) as count FROM items').get() as { count: number };
const ftsCount = db.prepare('SELECT COUNT(*) as count FROM items_fts').get() as { count: number };

console.log(`\n✅ Database emptied successfully!`);
console.log(`Items count: ${afterCount.count}`);
console.log(`FTS index count: ${ftsCount.count}`);

// Close database
db.close();

console.log('\nNext steps:');
console.log('  - Run "npm run db:seed" to populate with sample data');
console.log('  - Run "npm run db:verify" to verify the database state');
