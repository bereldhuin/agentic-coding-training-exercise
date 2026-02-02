package com.lebonpoint.infrastructure.persistence

import com.lebonpoint.domain.entities.CreateItemData
import com.lebonpoint.domain.entities.Item
import com.lebonpoint.domain.entities.PatchItemData
import com.lebonpoint.domain.entities.ReplaceItemData
import com.lebonpoint.domain.entities.applyPatch
import com.lebonpoint.domain.repositories.ItemFilters
import com.lebonpoint.domain.repositories.ItemPage
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.domain.repositories.SortDirection
import com.lebonpoint.domain.repositories.SortOptions
import com.lebonpoint.domain.valueobjects.ItemCondition
import com.lebonpoint.domain.valueobjects.ItemStatus
import com.lebonpoint.domain.valueobjects.toItemCondition
import com.lebonpoint.domain.valueobjects.toItemStatus
import com.lebonpoint.domain.valueobjects.toDatabaseString
import com.lebonpoint.domain.valueobjects.serializeToJson
import com.lebonpoint.domain.valueobjects.deserializeToImages
import com.lebonpoint.shared.decodeCursor
import com.lebonpoint.shared.encodeCursor
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.json.Json
import java.sql.Connection
import java.sql.Types
import java.time.Instant
import java.time.format.DateTimeFormatter

/**
 * SQLite implementation of ItemRepository
 * Uses JDBC with connection pooling via HikariCP
 */
class SQLiteItemRepository(
    private val dbConfig: DatabaseConfig = DatabaseConfig
) : ItemRepository {

    companion object {
        private val ISO_FORMATTER = DateTimeFormatter.ISO_OFFSET_DATE_TIME

        /**
         * Parse ISO 8601 timestamp string to Instant
         */
        fun parseInstant(value: String?): Instant? {
            return value?.let {
                try {
                    Instant.parse(it)
                } catch (e: Exception) {
                    null
                }
            }
        }

        /**
         * Format Instant to ISO 8601 string
         */
        fun formatInstant(instant: Instant?): String? {
            return instant?.toString()
        }
    }

    override suspend fun create(data: CreateItemData): Item = withContext(Dispatchers.IO) {
        dbConfig.withConnection { conn ->
            val now = Instant.now()

            conn.prepareStatement("""
                INSERT INTO items (
                    title, description, price_cents, category, condition, status,
                    is_featured, city, postal_code, country, delivery_available,
                    created_at, updated_at, published_at, images
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """.trimIndent()).use { stmt ->
                stmt.setString(1, data.title)
                stmt.setString(2, data.description)
                stmt.setInt(3, data.priceCents)
                stmt.setString(4, data.category)
                stmt.setString(5, data.condition.toDatabaseString())
                stmt.setString(6, data.status.toDatabaseString())
                stmt.setInt(7, if (data.isFeatured) 1 else 0)
                stmt.setString(8, data.city)
                stmt.setString(9, data.postalCode)
                stmt.setString(10, data.country)
                stmt.setInt(11, if (data.deliveryAvailable) 1 else 0)
                stmt.setString(12, formatInstant(now))
                stmt.setString(13, formatInstant(now))
                stmt.setString(14, null) // published_at
                stmt.setString(15, data.images.serializeToJson())

                stmt.executeUpdate()

                val generatedId = conn.createStatement().use { keyStmt ->
                    keyStmt.executeQuery("SELECT last_insert_rowid()").use { rs ->
                        if (rs.next()) {
                            rs.getInt(1)
                        } else {
                            throw IllegalStateException("Failed to get generated ID")
                        }
                    }
                }

                upsertFullTextIndex(conn = conn, id = generatedId, title = data.title, description = data.description)

                Item(
                    id = generatedId,
                    title = data.title,
                    description = data.description,
                    priceCents = data.priceCents,
                    category = data.category,
                    condition = data.condition,
                    status = data.status,
                    isFeatured = data.isFeatured,
                    city = data.city,
                    postalCode = data.postalCode,
                    country = data.country,
                    deliveryAvailable = data.deliveryAvailable,
                    createdAt = now,
                    updatedAt = now,
                    publishedAt = null,
                    images = data.images
                )
            }
        }
    }

    override suspend fun findById(id: Int): Item? = withContext(Dispatchers.IO) {
        dbConfig.withConnection { conn ->
            conn.prepareStatement("SELECT * FROM items WHERE id = ?").use { stmt ->
                stmt.setInt(1, id)
                stmt.executeQuery().use { rs ->
                    if (rs.next()) {
                        mapRowToItem(rs)
                    } else {
                        null
                    }
                }
            }
        }
    }

    override suspend fun findAll(
        filters: ItemFilters,
        sort: SortOptions,
        limit: Int,
        cursor: String?
    ): ItemPage = withContext(Dispatchers.IO) {
        dbConfig.withConnection { conn ->
            val whereClause = buildWhereClause(filters)
            val orderByClause = buildOrderByClause(sort)
            val cursorFilter = buildCursorFilter(cursor)

            val combinedWhere = listOfNotNull(
                whereClause,
                cursorFilter
            ).let { clauses ->
                if (clauses.isNotEmpty()) "WHERE ${clauses.joinToString(" AND ")}" else ""
            }

            val sql = """
                SELECT * FROM items
                $combinedWhere
                $orderByClause
                LIMIT ?
            """.trimIndent()

            conn.prepareStatement(sql).use { stmt ->
                var paramIndex = 1
                paramIndex = setFilterParameters(stmt, filters, paramIndex)
                stmt.setInt(paramIndex, limit.coerceIn(1, 100))

                stmt.executeQuery().use { rs ->
                    val items = mutableListOf<Item>()
                    while (rs.next()) {
                        items.add(mapRowToItem(rs))
                    }

                    val nextCursor = items.lastOrNull()?.let {
                        if (items.size == limit) encodeCursor(it.id) else null
                    }

                    ItemPage(items, nextCursor)
                }
            }
        }
    }

    override suspend fun search(
        query: String,
        filters: ItemFilters,
        sort: SortOptions,
        limit: Int,
        cursor: String?
    ): ItemPage = withContext(Dispatchers.IO) {
        dbConfig.withConnection { conn ->
            val whereClause = buildWhereClause(filters)
            val cursorFilter = buildCursorFilter(cursor)

            val combinedWhere = listOfNotNull(
                "items_fts MATCH ?",
                whereClause,
                cursorFilter
            ).let { clauses ->
                "WHERE ${clauses.joinToString(" AND ")}"
            }

            // Use BM25 ranking for relevance
            val orderByClause = if (sort.field == "created_at" && sort.direction == SortDirection.DESC) {
                // Keep default sort by relevance
                "ORDER BY bm25(items_fts) ASC, items.id ASC"
            } else {
                buildOrderByClause(sort)
            }

            val sql = """
                SELECT items.* FROM items
                INNER JOIN items_fts ON items.id = items_fts.rowid
                $combinedWhere
                $orderByClause
                LIMIT ?
            """.trimIndent()

            conn.prepareStatement(sql).use { stmt ->
                var paramIndex = 1
                stmt.setString(paramIndex++, query) // FTS search term
                paramIndex = setFilterParameters(stmt, filters, paramIndex)
                stmt.setInt(paramIndex, limit.coerceIn(1, 100))

                stmt.executeQuery().use { rs ->
                    val items = mutableListOf<Item>()
                    while (rs.next()) {
                        items.add(mapRowToItem(rs))
                    }

                    val nextCursor = items.lastOrNull()?.let {
                        if (items.size == limit) encodeCursor(it.id) else null
                    }

                    ItemPage(items, nextCursor)
                }
            }
        }
    }

    override suspend fun update(id: Int, data: ReplaceItemData): Item? = withContext(Dispatchers.IO) {
        dbConfig.withConnection { conn ->
            val now = Instant.now()

            conn.prepareStatement("""
                UPDATE items SET
                    title = ?, description = ?, price_cents = ?, category = ?, condition = ?,
                    status = ?, is_featured = ?, city = ?, postal_code = ?, country = ?,
                    delivery_available = ?, updated_at = ?, images = ?
                WHERE id = ?
            """.trimIndent()).use { stmt ->
                stmt.setString(1, data.title)
                stmt.setString(2, data.description)
                stmt.setInt(3, data.priceCents)
                stmt.setString(4, data.category)
                stmt.setString(5, data.condition.toDatabaseString())
                stmt.setString(6, data.status.toDatabaseString())
                stmt.setInt(7, if (data.isFeatured) 1 else 0)
                stmt.setString(8, data.city)
                stmt.setString(9, data.postalCode)
                stmt.setString(10, data.country)
                stmt.setInt(11, if (data.deliveryAvailable) 1 else 0)
                stmt.setString(12, formatInstant(now))
                stmt.setString(13, data.images.serializeToJson())
                stmt.setInt(14, id)

                val rows = stmt.executeUpdate()
                if (rows > 0) {
                    upsertFullTextIndex(conn = conn, id = id, title = data.title, description = data.description)
                    conn.prepareStatement("SELECT * FROM items WHERE id = ?").use { selectStmt ->
                        selectStmt.setInt(1, id)
                        selectStmt.executeQuery().use { rs ->
                            if (rs.next()) {
                                mapRowToItem(rs)
                            } else {
                                null
                            }
                        }
                    }
                } else {
                    null
                }
            }
        }
    }

    override suspend fun patch(id: Int, data: PatchItemData): Item? = withContext(Dispatchers.IO) {
        val existing = findById(id) ?: return@withContext null

        val updated = existing.applyPatch(data)

        dbConfig.withConnection { conn ->
            val setClauses = mutableListOf<String>()
            val values = mutableListOf<Any?>()

            data.title?.let { setClauses.add("title = ?"); values.add(it) }
            data.description?.let { setClauses.add("description = ?"); values.add(it) }
            data.priceCents?.let { setClauses.add("price_cents = ?"); values.add(it) }
            data.category?.let { setClauses.add("category = ?"); values.add(it) }
            data.condition?.let { setClauses.add("condition = ?"); values.add(it.toDatabaseString()) }
            data.status?.let { setClauses.add("status = ?"); values.add(it.toDatabaseString()) }
            data.isFeatured?.let { setClauses.add("is_featured = ?"); values.add(if (it) 1 else 0) }
            data.city?.let { setClauses.add("city = ?"); values.add(it) }
            data.postalCode?.let { setClauses.add("postal_code = ?"); values.add(it) }
            data.country?.let { setClauses.add("country = ?"); values.add(it) }
            data.deliveryAvailable?.let { setClauses.add("delivery_available = ?"); values.add(if (it) 1 else 0) }
            data.images?.let { setClauses.add("images = ?"); values.add(it.serializeToJson()) }

            setClauses.add("updated_at = ?")
            values.add(formatInstant(Instant.now()))

            if (setClauses.isEmpty()) {
                return@withConnection existing
            }

            val sql = "UPDATE items SET ${setClauses.joinToString(", ")} WHERE id = ?"
            values.add(id)

            conn.prepareStatement(sql).use { stmt ->
                values.forEachIndexed { index, value ->
                    when (value) {
                        is String -> stmt.setString(index + 1, value)
                        is Int -> stmt.setInt(index + 1, value)
                        null -> stmt.setNull(index + 1, java.sql.Types.VARCHAR)
                    }
                }

                stmt.executeUpdate()
                upsertFullTextIndex(conn = conn, id = id, title = updated.title, description = updated.description)
            }

            updated
        }
    }

    override suspend fun delete(id: Int): Boolean = withContext(Dispatchers.IO) {
        dbConfig.withConnection { conn ->
            conn.prepareStatement("DELETE FROM items WHERE id = ?").use { stmt ->
                stmt.setInt(1, id)
                val deleted = stmt.executeUpdate() > 0
                if (deleted) {
                    deleteFullTextIndex(conn, id)
                }
                deleted
            }
        }
    }

    /**
     * Map database row to Item entity
     */
    private fun mapRowToItem(rs: java.sql.ResultSet): Item {
        return Item(
            id = rs.getInt("id"),
            title = rs.getString("title"),
            description = rs.getString("description"),
            priceCents = rs.getInt("price_cents"),
            category = rs.getString("category"),
            condition = rs.getString("condition").toItemCondition() ?: ItemCondition.UNKNOWN,
            status = rs.getString("status").toItemStatus() ?: ItemStatus.DRAFT,
            isFeatured = rs.getInt("is_featured") == 1,
            city = rs.getString("city"),
            postalCode = rs.getString("postal_code"),
            country = rs.getString("country"),
            deliveryAvailable = rs.getInt("delivery_available") == 1,
            createdAt = parseInstant(rs.getString("created_at")) ?: Instant.EPOCH,
            updatedAt = parseInstant(rs.getString("updated_at")) ?: Instant.EPOCH,
            publishedAt = parseInstant(rs.getString("published_at")),
            images = rs.getString("images").deserializeToImages()
        )
    }

    /**
     * Build WHERE clause from filters
     */
    private fun buildWhereClause(filters: ItemFilters): String? {
        val conditions = mutableListOf<String>()

        filters.status?.let { conditions.add("status = ?") }
        filters.category?.let { conditions.add("category = ?") }
        filters.city?.let { conditions.add("city = ?") }
        filters.postalCode?.let { conditions.add("postal_code = ?") }
        filters.isFeatured?.let { conditions.add("is_featured = ?") }
        filters.deliveryAvailable?.let { conditions.add("delivery_available = ?") }

        return conditions.takeIf { it.isNotEmpty() }?.joinToString(" AND ")
    }

    /**
     * Build ORDER BY clause from sort options
     */
    private fun buildOrderByClause(sort: SortOptions): String {
        val direction = if (sort.direction == SortDirection.ASC) "ASC" else "DESC"
        return "ORDER BY ${sort.field} $direction"
    }

    /**
     * Build cursor filter for pagination
     */
    private fun buildCursorFilter(cursor: String?): String? {
        val cursorData = decodeCursor(cursor) ?: return null
        return "id > ${cursorData.id}"
    }

    /**
     * Set filter parameters on prepared statement
     */
    private fun setFilterParameters(
        stmt: java.sql.PreparedStatement,
        filters: ItemFilters,
        startIndex: Int
    ): Int {
        var index = startIndex

        filters.status?.let { stmt.setString(index++, it) }
        filters.category?.let { stmt.setString(index++, it) }
        filters.city?.let { stmt.setString(index++, it) }
        filters.postalCode?.let { stmt.setString(index++, it) }
        filters.isFeatured?.let { stmt.setInt(index++, if (it) 1 else 0) }
        filters.deliveryAvailable?.let { stmt.setInt(index++, if (it) 1 else 0) }

        return index
    }

    private fun upsertFullTextIndex(
        conn: Connection,
        id: Int,
        title: String?,
        description: String?
    ) {
        conn.prepareStatement("""
            INSERT OR REPLACE INTO items_fts(rowid, title, description)
            VALUES (?, ?, ?)
        """.trimIndent()).use { stmt ->
            stmt.setInt(1, id)
            stmt.setString(2, title)
            if (description != null) {
                stmt.setString(3, description)
            } else {
                stmt.setNull(3, Types.VARCHAR)
            }
            stmt.executeUpdate()
        }
    }

    private fun deleteFullTextIndex(conn: Connection, id: Int) {
        conn.prepareStatement("DELETE FROM items_fts WHERE rowid = ?").use { stmt ->
            stmt.setInt(1, id)
            stmt.executeUpdate()
        }
    }
}

/**
 * Helper function for non-null list of strings
 */
private fun <T> listOfNotNull(vararg elements: T?): List<T> {
    return elements.filterNotNull().toList()
}
