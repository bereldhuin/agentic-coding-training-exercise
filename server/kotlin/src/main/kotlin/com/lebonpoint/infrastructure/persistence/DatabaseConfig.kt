package com.lebonpoint.infrastructure.persistence

import com.lebonpoint.shared.Configuration
import com.zaxxer.hikari.HikariConfig
import com.zaxxer.hikari.HikariDataSource
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.nio.file.Paths
import java.sql.Connection
import javax.sql.DataSource

/**
 * Database configuration and connection pool management
 */
object DatabaseConfig {
    private var dataSource: HikariDataSource? = null

    /**
     * Get or create the data source
     * Uses HikariCP for connection pooling
     */
    fun getDataSource(): DataSource {
        return dataSource ?: createDataSource().also { dataSource = it }
    }

    /**
     * Create a new HikariCP data source
     */
    private fun createDataSource(): HikariDataSource {
        val config = HikariConfig().apply {
            // JDBC URL for SQLite
            val dbPath = Configuration.databasePath
            jdbcUrl = if (dbPath == ":memory:") {
                "jdbc:sqlite::memory:"
            } else {
                // Resolve relative paths from project root
                val resolvedPath = if (dbPath.startsWith("/")) {
                    dbPath
                } else {
                    Paths.get(System.getProperty("user.dir")).resolve(dbPath).normalize().toString()
                }
                "jdbc:sqlite:$resolvedPath"
            }

            // SQLite driver
            driverClassName = "org.sqlite.JDBC"

            // Pool configuration
            maximumPoolSize = if (dbPath == ":memory:") 1 else Configuration.poolSize
            minimumIdle = if (dbPath == ":memory:") 1 else 1

            // Connection timeout
            connectionTimeout = 30000 // 30 seconds

            // SQLite doesn't support connection validation well
            connectionTestQuery = "SELECT 1"

            // Pool name
            poolName = "LeBonPointHikariPool"

            // SQLite specific: Disable connections being kept alive
            // SQLite uses file-based locking, so we need to be careful
            idleTimeout = if (dbPath == ":memory:") Long.MAX_VALUE else 600000 // 10 minutes
            maxLifetime = if (dbPath == ":memory:") Long.MAX_VALUE else 1800000 // 30 minutes
        }

        return HikariDataSource(config)
    }

    /**
     * Get a connection from the pool
     * Automatically closes when done (use with use block)
     */
    suspend fun <T> withConnection(block: suspend (Connection) -> T): T {
        return withContext(Dispatchers.IO) {
            getDataSource().connection.use { block(it) }
        }
    }

    /**
     * Close the connection pool
     * Call this on application shutdown
     */
    fun shutdown() {
        dataSource?.close()
        dataSource = null
    }

    /**
     * Get connection without auto-close (for advanced use cases)
     * Make sure to close it manually
     */
    fun getConnection(): Connection {
        return getDataSource().connection
    }

    /**
     * Initialize database schema if needed
     * This is called on application startup
     */
    suspend fun initializeSchema() {
        withConnection { conn ->
            // Check if items table exists
            val tableExists = conn.metaData.getTables(
                null, null, "items", null
            ).use { rs ->
                rs.next()
            }

            if (!tableExists) {
                createSchema(conn)
            }
        }
    }

    /**
     * Create database schema
     */
    private fun createSchema(conn: Connection) {
        conn.createStatement().use { stmt ->
            // Create items table
            stmt.execute("""
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
            """.trimIndent())

            // Create indexes
            stmt.execute("CREATE INDEX IF NOT EXISTS idx_items_status ON items(status)")
            stmt.execute("CREATE INDEX IF NOT EXISTS idx_items_category ON items(category)")
            stmt.execute("CREATE INDEX IF NOT EXISTS idx_items_condition ON items(condition)")
            stmt.execute("CREATE INDEX IF NOT EXISTS idx_items_city ON items(city)")
            stmt.execute("CREATE INDEX IF NOT EXISTS idx_items_price_cents ON items(price_cents)")

            // Create FTS5 virtual table
            stmt.execute("""
                CREATE VIRTUAL TABLE IF NOT EXISTS items_fts USING fts5(
                    title,
                    description,
                    content=items,
                    content_rowid=rowid,
                    tokenize='unicode61 remove_diacritics 1'
                )
            """.trimIndent())
        }
    }

    /**
     * Initialize for testing (in-memory database)
     */
    suspend fun initializeForTesting() {
        // Force in-memory database for tests
        System.setProperty("DATABASE_PATH", ":memory:")
        shutdown() // Close any existing connections
        initializeSchema()
    }
}
