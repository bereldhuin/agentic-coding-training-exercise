import Database from 'better-sqlite3';
import { getDatabase } from './sqlite.js';

/**
 * Initialize database schema
 * This is called when using an in-memory database for tests
 */
export function initializeDatabaseSchema(db: Database.Database): void {
  // Create items table
  db.exec(`
    CREATE TABLE IF NOT EXISTS items (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      title TEXT NOT NULL CHECK(length(title) >= 3 AND length(title) <= 200),
      description TEXT,
      price_cents INTEGER NOT NULL CHECK(price_cents >= 0),
      category TEXT,
      condition TEXT NOT NULL CHECK(condition IN ('new', 'like_new', 'good', 'fair', 'parts', 'unknown')),
      status TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'active', 'reserved', 'sold', 'archived')),
      is_featured INTEGER NOT NULL DEFAULT 0,
      city TEXT,
      postal_code TEXT,
      country TEXT NOT NULL DEFAULT 'FR',
      delivery_available INTEGER NOT NULL DEFAULT 0,
      created_at TEXT NOT NULL DEFAULT (datetime('now')),
      updated_at TEXT NOT NULL DEFAULT (datetime('now')),
      published_at TEXT,
      images TEXT NOT NULL DEFAULT '[]'
    )
  `);

  // Create indexes
  db.exec(`
    CREATE INDEX IF NOT EXISTS idx_items_status ON items(status);
    CREATE INDEX IF NOT EXISTS idx_items_category ON items(category);
    CREATE INDEX IF NOT EXISTS idx_items_condition ON items(condition);
    CREATE INDEX IF NOT EXISTS idx_items_city ON items(city);
    CREATE INDEX IF NOT EXISTS idx_items_price_cents ON items(price_cents);
  `);

  // Create FTS5 virtual table
  db.exec(`
    CREATE VIRTUAL TABLE IF NOT EXISTS items_fts USING fts5(
      title,
      description,
      content=items,
      content_rowid=rowid,
      tokenize='unicode61 remove_diacritics 1'
    )
  `);

  // Create triggers to keep FTS table in sync with items table
  db.exec(`
    CREATE TRIGGER IF NOT EXISTS items_ai AFTER INSERT ON items BEGIN
      INSERT INTO items_fts(rowid, title, description)
      VALUES (new.rowid, new.title, new.description);
    END;
  `);

  db.exec(`
    CREATE TRIGGER IF NOT EXISTS items_au AFTER UPDATE ON items BEGIN
      UPDATE items_fts
      SET title = new.title, description = new.description
      WHERE rowid = new.rowid;
    END;
  `);

  db.exec(`
    CREATE TRIGGER IF NOT EXISTS items_ad AFTER DELETE ON items BEGIN
      DELETE FROM items_fts WHERE rowid = old.rowid;
    END;
  `);
}

/**
 * Initialize database if needed
 */
export function ensureDatabaseInitialized(): void {
  const db = getDatabase();

  // Check if items table exists
  const tableExists = db.prepare(`
    SELECT name FROM sqlite_master
    WHERE type='table' AND name='items'
  `).get();

  if (!tableExists) {
    initializeDatabaseSchema(db);
  }
}
