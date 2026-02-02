"""SQLite implementation of Item repository."""

import sqlite3
from datetime import datetime, timezone
from pathlib import Path

from ...domain.entities.item import (
    Item,
    CreateItemData,
    UpdateItemData,
    PatchItemData,
    item_from_row,
    serialize_images_column,
)
from ...domain.repositories.item_repository import (
    ItemRepository,
    FilterOptions,
    SortOptions,
    SortField,
    SortDirection,
    PaginationOptions,
    SearchOptions,
    PaginatedResult,
)
from ...shared.config import config
from ...shared.cursor import encode_cursor, decode_cursor, CursorData


class SQLiteItemRepository(ItemRepository):
    """SQLite repository adapter for items."""

    def __init__(self, db_path: Path | None = None) -> None:
        """
        Initialize the repository.

        Args:
            db_path: Path to SQLite database. If None, uses config.database_path
        """
        self._db_path = db_path or config.database_path

    def _get_connection(self) -> sqlite3.Connection:
        """
        Create a new database connection.

        Returns:
            SQLite connection with row factory configured
        """
        conn = sqlite3.connect(str(self._db_path))
        conn.row_factory = sqlite3.Row
        return conn

    async def create(self, data: CreateItemData) -> Item:
        """Create a new item."""
        now = datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")

        images_json = serialize_images_column(data.images or [])

        with self._get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute(
                """
                INSERT INTO items (
                    title, description, price_cents, category, condition, status,
                    is_featured, city, postal_code, country, delivery_available,
                    images, created_at, updated_at
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """,
                (
                    data.title,
                    data.description,
                    data.price_cents,
                    data.category,
                    data.condition.value,
                    data.status.value,
                    1 if data.is_featured else 0,
                    data.city,
                    data.postal_code,
                    data.country,
                    1 if data.delivery_available else 0,
                    images_json,
                    now,
                    now,
                ),
            )
            conn.commit()
            item_id = cursor.lastrowid

        created = await self.find_by_id(item_id)
        if created is None:
            raise RuntimeError("Failed to create item")
        return created

    async def find_by_id(self, id: int) -> Item | None:
        """Find an item by ID."""
        with self._get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM items WHERE id = ?", (id,))
            row = cursor.fetchone()

            if row is None:
                return None

            return item_from_row(dict(row))

    async def find_all(
        self,
        filters: FilterOptions | None = None,
        sort: SortOptions | None = None,
        pagination: PaginationOptions | None = None,
    ) -> PaginatedResult:
        """List all items with filtering, sorting, and pagination."""
        where_conditions: list[str] = []
        params: list = []

        if filters:
            if filters.status:
                where_conditions.append("status = ?")
                params.append(filters.status)
            if filters.category:
                where_conditions.append("category = ?")
                params.append(filters.category)
            if filters.city:
                where_conditions.append("city = ?")
                params.append(filters.city)
            if filters.postal_code:
                where_conditions.append("postal_code = ?")
                params.append(filters.postal_code)
            if filters.is_featured is not None:
                where_conditions.append("is_featured = ?")
                params.append(1 if filters.is_featured else 0)
            if filters.delivery_available is not None:
                where_conditions.append("delivery_available = ?")
                params.append(1 if filters.delivery_available else 0)

        # Handle cursor pagination
        if pagination and pagination.cursor:
            cursor_data = decode_cursor(pagination.cursor)
            if cursor_data:
                sort_field = sort.field if sort else SortField.CREATED_AT
                sort_dir = sort.direction if sort else SortDirection.DESC

                if sort_dir == SortDirection.DESC:
                    where_conditions.append(f"({sort_field.value}, id) < (?, ?)")
                else:
                    where_conditions.append(f"({sort_field.value}, id) > (?, ?)")
                params.extend([cursor_data.created_at, cursor_data.id])

        where_clause = f"WHERE {' AND '.join(where_conditions)}" if where_conditions else ""

        # Build ORDER BY clause
        sort_field = sort.field if sort else SortField.CREATED_AT
        sort_dir = sort.direction if sort else SortDirection.DESC
        order_clause = f"ORDER BY {sort_field.value} {sort_dir.value.upper()}, id {sort_dir.value.upper()}"

        # Build LIMIT clause
        limit = pagination.limit if pagination else 20
        limit_clause = f"LIMIT {limit + 1}"

        # Execute query
        query = f"SELECT * FROM items {where_clause} {order_clause} {limit_clause}"
        with self._get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute(query, params)
            rows = cursor.fetchall()

        # Check for next page
        has_next_page = len(rows) > limit
        items = [item_from_row(dict(row)) for row in rows[:limit]]

        # Build next cursor
        next_cursor: str | None = None
        if has_next_page and items:
            last_item = items[-1]
            cursor_data = CursorData(id=last_item.id, created_at=last_item.created_at)
            next_cursor = encode_cursor(cursor_data)

        return PaginatedResult(items=items, next_cursor=next_cursor)

    async def search(self, options: SearchOptions) -> PaginatedResult:
        """Full-text search using FTS5."""
        where_conditions: list[str] = []
        params: list = []

        if options.filters:
            if options.filters.status:
                where_conditions.append("items.status = ?")
                params.append(options.filters.status)
            if options.filters.category:
                where_conditions.append("items.category = ?")
                params.append(options.filters.category)
            if options.filters.city:
                where_conditions.append("items.city = ?")
                params.append(options.filters.city)
            if options.filters.postal_code:
                where_conditions.append("items.postal_code = ?")
                params.append(options.filters.postal_code)
            if options.filters.is_featured is not None:
                where_conditions.append("items.is_featured = ?")
                params.append(1 if options.filters.is_featured else 0)
            if options.filters.delivery_available is not None:
                where_conditions.append("items.delivery_available = ?")
                params.append(1 if options.filters.delivery_available else 0)

        # Handle cursor pagination
        if options.pagination and options.pagination.cursor:
            cursor_data = decode_cursor(options.pagination.cursor)
            if cursor_data:
                where_conditions.append("(items.created_at, items.id) < (?, ?)")
                params.extend([cursor_data.created_at, cursor_data.id])

        filters_clause = f"AND {' AND '.join(where_conditions)}" if where_conditions else ""

        # Build LIMIT clause
        limit = options.pagination.limit if options.pagination else 20
        limit_clause = f"LIMIT {limit + 1}"

        # Escape FTS search query
        fts_query = options.query.replace('"', '""')

        # Execute FTS search with BM25 ranking
        search_query = f"""
            SELECT items.* FROM items
            INNER JOIN items_fts ON items.id = items_fts.rowid
            WHERE items_fts MATCH ?
            {filters_clause}
            ORDER BY bm25(items_fts), items.created_at DESC, items.id DESC
            {limit_clause}
        """

        with self._get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute(search_query, [fts_query] + params)
            rows = cursor.fetchall()

        # Check for next page
        has_next_page = len(rows) > limit
        items = [item_from_row(dict(row)) for row in rows[:limit]]

        # Build next cursor
        next_cursor: str | None = None
        if has_next_page and items:
            last_item = items[-1]
            cursor_data = CursorData(id=last_item.id, created_at=last_item.created_at)
            next_cursor = encode_cursor(cursor_data)

        return PaginatedResult(items=items, next_cursor=next_cursor)

    async def update(self, id: int, data: UpdateItemData) -> Item | None:
        """Update (replace) an item."""
        now = datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")

        images_json = serialize_images_column(data.images or [])

        with self._get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute(
                """
                UPDATE items SET
                    title = ?, description = ?, price_cents = ?, category = ?,
                    condition = ?, status = ?, is_featured = ?, city = ?,
                    postal_code = ?, country = ?, delivery_available = ?,
                    images = ?, updated_at = ?
                WHERE id = ?
                """,
                (
                    data.title,
                    data.description,
                    data.price_cents,
                    data.category,
                    data.condition.value,
                    data.status.value,
                    1 if data.is_featured else 0,
                    data.city,
                    data.postal_code,
                    data.country,
                    1 if data.delivery_available else 0,
                    images_json,
                    now,
                    id,
                ),
            )
            conn.commit()

            if cursor.rowcount == 0:
                return None

        return await self.find_by_id(id)

    async def patch(self, id: int, data: PatchItemData) -> Item | None:
        """Partially update an item."""
        now = datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")

        updates: list[str] = []
        params: list = []

        if data.title is not None:
            updates.append("title = ?")
            params.append(data.title)
        if data.description is not None:
            updates.append("description = ?")
            params.append(data.description)
        if data.price_cents is not None:
            updates.append("price_cents = ?")
            params.append(data.price_cents)
        if data.category is not None:
            updates.append("category = ?")
            params.append(data.category)
        if data.condition is not None:
            updates.append("condition = ?")
            params.append(data.condition.value)
        if data.status is not None:
            updates.append("status = ?")
            params.append(data.status.value)
        if data.is_featured is not None:
            updates.append("is_featured = ?")
            params.append(1 if data.is_featured else 0)
        if data.city is not None:
            updates.append("city = ?")
            params.append(data.city)
        if data.postal_code is not None:
            updates.append("postal_code = ?")
            params.append(data.postal_code)
        if data.country is not None:
            updates.append("country = ?")
            params.append(data.country)
        if data.delivery_available is not None:
            updates.append("delivery_available = ?")
            params.append(1 if data.delivery_available else 0)
        if data.images is not None:
            updates.append("images = ?")
            params.append(serialize_images_column(data.images))

        updates.append("updated_at = ?")
        params.append(now)
        params.append(id)

        with self._get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute(f"UPDATE items SET {', '.join(updates)} WHERE id = ?", params)
            conn.commit()

            if cursor.rowcount == 0:
                return None

        return await self.find_by_id(id)

    async def delete(self, id: int) -> bool:
        """Delete an item."""
        with self._get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute("DELETE FROM items WHERE id = ?", (id,))
            conn.commit()
            return cursor.rowcount > 0

    async def exists(self, id: int) -> bool:
        """Check if an item exists."""
        with self._get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT 1 FROM items WHERE id = ? LIMIT 1", (id,))
            return cursor.fetchone() is not None
