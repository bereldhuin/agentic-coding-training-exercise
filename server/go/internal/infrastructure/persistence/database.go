package persistence

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Database wraps the database connection
type Database struct {
	db *sqlx.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dataSourceName string) (*Database, error) {
	// Open database connection
	db, err := sqlx.Connect("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	// SQLite doesn't benefit much from multiple connections due to its locking model,
	// but we set reasonable defaults for consistency with other databases
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(25)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime

	// Enable foreign key support
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Enable WAL mode for better concurrent read performance
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Set busy timeout to handle concurrent access
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}

	// Initialize FTS5 full-text search table
	if err := initFTS5(db); err != nil {
		return nil, fmt.Errorf("failed to initialize FTS5: %w", err)
	}

	return &Database{db: db}, nil
}

// initFTS5 creates the FTS5 virtual table for full-text search
func initFTS5(db *sqlx.DB) error {
	// Check if items_fts table already exists
	var count int
	if err := db.Get(&count, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='items_fts'"); err != nil {
		return fmt.Errorf("failed to check for items_fts table: %w", err)
	}

	// Create FTS5 virtual table if it doesn't exist
	if count == 0 {
		schema := `
			CREATE VIRTUAL TABLE items_fts USING fts5(
				title,
				description,
				content='items',
				content_rowid='id'
			)
		`
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to create items_fts table: %w", err)
		}

		// Populate FTS5 table with existing items
		if _, err := db.Exec(`
			INSERT INTO items_fts(rowid, title, description)
			SELECT id, title, coalesce(description, '') FROM items
		`); err != nil {
			return fmt.Errorf("failed to populate items_fts table: %w", err)
		}

		// Create triggers to keep FTS5 table in sync with items table
		triggers := []string{
			// INSERT trigger
			`CREATE TRIGGER items_ai AFTER INSERT ON items BEGIN
				INSERT INTO items_fts(rowid, title, description)
				VALUES (NEW.id, NEW.title, coalesce(NEW.description, ''));
			END`,
			// UPDATE trigger
			`CREATE TRIGGER items_au AFTER UPDATE ON items BEGIN
				UPDATE items_fts SET title = NEW.title, description = coalesce(NEW.description, '') WHERE rowid = NEW.id;
			END`,
			// DELETE trigger
			`CREATE TRIGGER items_ad AFTER DELETE ON items BEGIN
				DELETE FROM items_fts WHERE rowid = OLD.id;
			END`,
		}

		for _, trigger := range triggers {
			if _, err := db.Exec(trigger); err != nil {
				return fmt.Errorf("failed to create FTS5 trigger: %w", err)
			}
		}
	}

	return nil
}

// GetDB returns the underlying sqlx.DB
func (d *Database) GetDB() *sqlx.DB {
	return d.db
}

// Ping checks if the database is reachable
func (d *Database) Ping() error {
	return d.db.Ping()
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// BeginTx starts a new transaction
func (d *Database) BeginTx() (*sqlx.Tx, error) {
	return d.db.Beginx()
}

// ExecuteInTransaction executes a function within a transaction
// The function will be committed if it returns nil, otherwise rolled back
func (d *Database) ExecuteInTransaction(fn func(tx *sqlx.Tx) error) error {
	tx, err := d.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
