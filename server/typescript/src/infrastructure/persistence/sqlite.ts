import Database from 'better-sqlite3';
import { getDatabasePath, ensureDatabaseDirectory } from '../../shared/database.js';

/**
 * Database configuration
 */
interface DatabaseConfig {
  readonly filename: string;
  readonly readonly?: boolean;
  readonly fileMustExist?: boolean;
}

/**
 * Get database path
 * Uses the shared database configuration from src/shared/database.ts
 * The database is located at ../database/db.sqlite relative to the typescript service
 */
function getDatabaseConfigPath(): string {
  return getDatabasePath();
}

/**
 * Get database configuration
 */
function getDatabaseConfig(): DatabaseConfig {
  const filename = getDatabaseConfigPath();

  // Ensure directory exists for file-based databases
  if (filename !== ':memory:') {
    ensureDatabaseDirectory();
  }

  return {
    filename,
    fileMustExist: false // Allow creation
  };
}

/**
 * Create database connection
 */
export function createDatabaseConnection(): Database.Database {
  const config = getDatabaseConfig();
  const options: Database.Options = {};
  if (config.readonly !== undefined) {
    options.readonly = config.readonly;
  }
  if (config.fileMustExist !== undefined) {
    options.fileMustExist = config.fileMustExist;
  }

  const db = new Database(config.filename, options);

  // Enable WAL mode for better concurrency (not for in-memory)
  if (config.filename !== ':memory:') {
    db.pragma('journal_mode = WAL');
  }

  // Enable foreign keys
  db.pragma('foreign_keys = ON');

  return db;
}

/**
 * Get singleton database instance
 */
let dbInstance: Database.Database | null = null;

export function getDatabase(): Database.Database {
  if (!dbInstance) {
    dbInstance = createDatabaseConnection();
  }
  return dbInstance;
}

/**
 * Close database connection
 */
export function closeDatabase(): void {
  if (dbInstance) {
    dbInstance.close();
    dbInstance = null;
  }
}

/**
 * Reset database (for tests)
 */
export function resetDatabase(): void {
  closeDatabase();
}
