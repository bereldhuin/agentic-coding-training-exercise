import { join, dirname } from 'node:path';
import { existsSync, mkdirSync, copyFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

/**
 * Shared database configuration
 *
 * The database is located at ../database/db.sqlite relative to the typescript service,
 * which resolves to /server/database/db.sqlite at the project level.
 * This allows the database to be shared across multiple services in a multi-service architecture.
 */

/**
 * Get the database path
 * Returns ':memory:' for test environments, otherwise returns the path to db.sqlite
 */
export function getDatabasePath(): string {
  // Check if we're in a test environment
  if (process.env.NODE_ENV === 'test' || process.env.VITEST === 'true') {
    return ':memory:'; // Use in-memory database for tests
  }

  // The database is in the parent directory's database folder
  // From /server/typescript/src/shared/database.ts
  // We need to go up 4 levels to reach /server/, then into database/
  // src/shared -> src -> typescript -> server -> database
  const dbPath = join(__dirname, '../../../', 'database', 'db.sqlite');

  return dbPath;
}

/**
 * Ensure the database directory exists
 * Creates the directory if it doesn't exist
 */
export function ensureDatabaseDirectory(): void {
  const dbPath = getDatabasePath();
  const dbDir = dirname(dbPath);

  if (!existsSync(dbDir)) {
    mkdirSync(dbDir, { recursive: true });
  }
}

/**
 * Migrate database from old location to new location
 * Copies database/marketplace.db to ../database/db.sqlite if it exists
 */
export function migrateDatabase(): boolean {
  // Old database path was at /server/typescript/database/marketplace.db
  // From src/shared/, go up 2 levels to reach typescript/, then into database/
  const oldDbPath = join(__dirname, '../../database/marketplace.db');
  const newDbPath = getDatabasePath();

  // Skip if old database doesn't exist
  if (!existsSync(oldDbPath)) {
    return false;
  }

  // Skip if new database already exists
  if (existsSync(newDbPath)) {
    console.warn('⚠️  Warning: Both old and new database files exist.');
    console.warn(`   Old: ${oldDbPath}`);
    console.warn(`   New: ${newDbPath}`);
    console.warn('   Using existing database at new location.');
    return false;
  }

  // Copy old database to new location
  console.log('Migrating database from old location...');
  console.log(`  From: ${oldDbPath}`);
  console.log(`  To: ${newDbPath}`);

  ensureDatabaseDirectory();
  copyFileSync(oldDbPath, newDbPath);

  console.log('✓ Database migrated successfully');
  console.log('⚠️  You can now safely remove the old database file:');
  console.log(`   rm ${oldDbPath}`);

  return true;
}

/**
 * Check if old database file exists
 */
export function hasOldDatabase(): boolean {
  // Old database path was at /server/typescript/database/marketplace.db
  const oldDbPath = join(__dirname, '../../database/marketplace.db');
  return existsSync(oldDbPath);
}
