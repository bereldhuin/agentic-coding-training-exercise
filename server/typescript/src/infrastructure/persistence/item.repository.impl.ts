import type { IItemRepository, PaginatedResult, FilterOptions, SortOptions, PaginationOptions, SearchOptions } from '../../domain/repositories/item.repository.port.js';
import type { Item } from '../../domain/entities/item.entity.js';
import { itemFromRow, serializeImagesColumn } from '../../domain/entities/item.entity.js';
import { createCursor, parseCursor, type CursorData } from '../../shared/types.js';
import { getDatabase } from './sqlite.js';

/**
 * SQLite repository adapter for items
 */
export class SQLiteItemRepository implements IItemRepository {
  private db = getDatabase();

  /**
   * Create a new item
   */
  async create(data: Omit<Item, 'id' | 'created_at' | 'updated_at' | 'published_at'>): Promise<Item> {
    const now = new Date().toISOString();

    const stmt = this.db.prepare(`
      INSERT INTO items (
        title, description, price_cents, category, condition, status,
        is_featured, city, postal_code, country, delivery_available,
        images, created_at, updated_at
      ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `);

    const result = stmt.run(
      data.title,
      data.description ?? null,
      data.price_cents,
      data.category ?? null,
      data.condition,
      data.status,
      data.is_featured ? 1 : 0,
      data.city ?? null,
      data.postal_code ?? null,
      data.country,
      data.delivery_available ? 1 : 0,
      serializeImagesColumn(data.images),
      now,
      now
    );

    const createdItem = await this.findById(result.lastInsertRowid as number);
    if (!createdItem) {
      throw new Error('Failed to create item');
    }

    return createdItem;
  }

  /**
   * Find item by ID
   */
  async findById(id: number): Promise<Item | null> {
    const stmt = this.db.prepare('SELECT * FROM items WHERE id = ?');
    const row = stmt.get(id) as Record<string, unknown> | undefined;

    if (!row) {
      return null;
    }

    return itemFromRow(row);
  }

  /**
   * List all items with filters, sorting, and pagination
   */
  async findAll(options: {
    filters?: FilterOptions;
    sort?: SortOptions;
    pagination?: PaginationOptions;
  }): Promise<PaginatedResult<Item>> {
    const { filters, sort, pagination } = options;

    // Build WHERE clause
    const whereConditions: string[] = [];
    const params: unknown[] = [];

    if (filters) {
      if (filters.status) {
        whereConditions.push('status = ?');
        params.push(filters.status);
      }
      if (filters.category) {
        whereConditions.push('category = ?');
        params.push(filters.category);
      }
      if (filters.min_price_cents !== undefined) {
        whereConditions.push('price_cents >= ?');
        params.push(filters.min_price_cents);
      }
      if (filters.max_price_cents !== undefined) {
        whereConditions.push('price_cents <= ?');
        params.push(filters.max_price_cents);
      }
      if (filters.city) {
        whereConditions.push('city = ?');
        params.push(filters.city);
      }
      if (filters.postal_code) {
        whereConditions.push('postal_code = ?');
        params.push(filters.postal_code);
      }
      if (filters.is_featured !== undefined) {
        whereConditions.push('is_featured = ?');
        params.push(filters.is_featured ? 1 : 0);
      }
      if (filters.delivery_available !== undefined) {
        whereConditions.push('delivery_available = ?');
        params.push(filters.delivery_available ? 1 : 0);
      }
    }

    // Handle cursor pagination
    if (pagination?.cursor) {
      const cursorData = parseCursor(pagination.cursor);
      if (cursorData) {
        const sortField = sort?.field ?? 'created_at';
        const sortDir = sort?.direction ?? 'desc';

        if (sortDir === 'desc') {
          whereConditions.push(`(${sortField}, id) < (?, ?)`);
        } else {
          whereConditions.push(`(${sortField}, id) > (?, ?)`);
        }
        params.push(cursorData.created_at);
        params.push(cursorData.id);
      }
    }

    const whereClause = whereConditions.length > 0 ? `WHERE ${whereConditions.join(' AND ')}` : '';

    // Build ORDER BY clause
    const sortField = sort?.field ?? 'created_at';
    const sortDir = sort?.direction ?? 'desc';
    const orderClause = `ORDER BY ${sortField} ${sortDir.toUpperCase()}, id ${sortDir.toUpperCase()}`;

    // Build LIMIT clause
    const limit = pagination?.limit ?? 20;
    const limitClause = `LIMIT ${limit + 1}`; // Fetch one extra to check for next page

    // Execute query
    const query = `SELECT * FROM items ${whereClause} ${orderClause} ${limitClause}`;
    const stmt = this.db.prepare(query);
    const rows = stmt.all(...params) as Record<string, unknown>[];

    // Check if there's a next page
    const hasNextPage = rows.length > limit;
    const items = rows.slice(0, limit).map(row => itemFromRow(row));

    // Build next cursor
    let nextCursor: string | undefined;
    if (hasNextPage && items.length > 0) {
      const lastItem = items[items.length - 1];
      const cursorData: CursorData = {
        id: lastItem.id,
        created_at: lastItem.created_at
      };
      nextCursor = createCursor(cursorData);
    }

    return {
      items,
      next_cursor: nextCursor
    };
  }

  /**
   * Full-text search using FTS5
   */
  async search(options: SearchOptions): Promise<PaginatedResult<Item>> {
    const { query, filters, pagination } = options;

    // Build WHERE clause for filters
    const whereConditions: string[] = [];
    const params: unknown[] = [];

    if (filters) {
      if (filters.status) {
        whereConditions.push('items.status = ?');
        params.push(filters.status);
      }
      if (filters.category) {
        whereConditions.push('items.category = ?');
        params.push(filters.category);
      }
      if (filters.min_price_cents !== undefined) {
        whereConditions.push('items.price_cents >= ?');
        params.push(filters.min_price_cents);
      }
      if (filters.max_price_cents !== undefined) {
        whereConditions.push('items.price_cents <= ?');
        params.push(filters.max_price_cents);
      }
      if (filters.city) {
        whereConditions.push('items.city = ?');
        params.push(filters.city);
      }
      if (filters.postal_code) {
        whereConditions.push('items.postal_code = ?');
        params.push(filters.postal_code);
      }
      if (filters.is_featured !== undefined) {
        whereConditions.push('items.is_featured = ?');
        params.push(filters.is_featured ? 1 : 0);
      }
      if (filters.delivery_available !== undefined) {
        whereConditions.push('items.delivery_available = ?');
        params.push(filters.delivery_available ? 1 : 0);
      }
    }

    // Handle cursor pagination
    if (pagination?.cursor) {
      const cursorData = parseCursor(pagination.cursor);
      if (cursorData) {
        whereConditions.push('(items.created_at, items.id) < (?, ?)');
        params.push(cursorData.created_at);
        params.push(cursorData.id);
      }
    }

    const filtersClause = whereConditions.length > 0 ? `AND ${whereConditions.join(' AND ')}` : '';

    // Build LIMIT clause
    const limit = pagination?.limit ?? 20;
    const limitClause = `LIMIT ${limit + 1}`;

    // Escape and prepare FTS search query
    // Simple escaping - replace quotes with double quotes
    const ftsQuery = query.replace(/"/g, '""');

    // Execute FTS search with BM25 ranking
    const searchQuery = `
      SELECT items.* FROM items
      INNER JOIN items_fts ON items.id = items_fts.rowid
      WHERE items_fts MATCH ?
      ${filtersClause}
      ORDER BY bm25(items_fts), items.created_at DESC, items.id DESC
      ${limitClause}
    `;

    const stmt = this.db.prepare(searchQuery);
    const rows = stmt.all(ftsQuery, ...params) as Record<string, unknown>[];

    // Check if there's a next page
    const hasNextPage = rows.length > limit;
    const items = rows.slice(0, limit).map(row => itemFromRow(row));

    // Build next cursor
    let nextCursor: string | undefined;
    if (hasNextPage && items.length > 0) {
      const lastItem = items[items.length - 1];
      const cursorData: CursorData = {
        id: lastItem.id,
        created_at: lastItem.created_at
      };
      nextCursor = createCursor(cursorData);
    }

    return {
      items,
      next_cursor: nextCursor
    };
  }

  /**
   * Update (replace) an item
   */
  async update(id: number, data: Omit<Item, 'id' | 'created_at' | 'updated_at' | 'published_at'>): Promise<Item | null> {
    const now = new Date().toISOString();

    const stmt = this.db.prepare(`
      UPDATE items SET
        title = ?,
        description = ?,
        price_cents = ?,
        category = ?,
        condition = ?,
        status = ?,
        is_featured = ?,
        city = ?,
        postal_code = ?,
        country = ?,
        delivery_available = ?,
        images = ?,
        updated_at = ?
      WHERE id = ?
    `);

    const result = stmt.run(
      data.title,
      data.description ?? null,
      data.price_cents,
      data.category ?? null,
      data.condition,
      data.status,
      data.is_featured ? 1 : 0,
      data.city ?? null,
      data.postal_code ?? null,
      data.country,
      data.delivery_available ? 1 : 0,
      serializeImagesColumn(data.images),
      now,
      id
    );

    if (result.changes === 0) {
      return null;
    }

    return this.findById(id);
  }

  /**
   * Partial update an item
   */
  async patch(id: number, data: Partial<Omit<Item, 'id' | 'created_at' | 'updated_at' | 'published_at'>>): Promise<Item | null> {
    const now = new Date().toISOString();

    // Build dynamic UPDATE statement
    const updates: string[] = [];
    const params: unknown[] = [];

    if (data.title !== undefined) {
      updates.push('title = ?');
      params.push(data.title);
    }
    if (data.description !== undefined) {
      updates.push('description = ?');
      params.push(data.description);
    }
    if (data.price_cents !== undefined) {
      updates.push('price_cents = ?');
      params.push(data.price_cents);
    }
    if (data.category !== undefined) {
      updates.push('category = ?');
      params.push(data.category);
    }
    if (data.condition !== undefined) {
      updates.push('condition = ?');
      params.push(data.condition);
    }
    if (data.status !== undefined) {
      updates.push('status = ?');
      params.push(data.status);
    }
    if (data.is_featured !== undefined) {
      updates.push('is_featured = ?');
      params.push(data.is_featured ? 1 : 0);
    }
    if (data.city !== undefined) {
      updates.push('city = ?');
      params.push(data.city);
    }
    if (data.postal_code !== undefined) {
      updates.push('postal_code = ?');
      params.push(data.postal_code);
    }
    if (data.country !== undefined) {
      updates.push('country = ?');
      params.push(data.country);
    }
    if (data.delivery_available !== undefined) {
      updates.push('delivery_available = ?');
      params.push(data.delivery_available ? 1 : 0);
    }
    if (data.images !== undefined) {
      updates.push('images = ?');
      params.push(serializeImagesColumn(data.images));
    }

    updates.push('updated_at = ?');
    params.push(now);

    params.push(id);

    const stmt = this.db.prepare(`UPDATE items SET ${updates.join(', ')} WHERE id = ?`);
    const result = stmt.run(...params);

    if (result.changes === 0) {
      return null;
    }

    return this.findById(id);
  }

  /**
   * Delete an item
   */
  async delete(id: number): Promise<boolean> {
    const stmt = this.db.prepare('DELETE FROM items WHERE id = ?');
    const result = stmt.run(id);

    return result.changes > 0;
  }

  /**
   * Check if an item exists
   */
  async exists(id: number): Promise<boolean> {
    const stmt = this.db.prepare('SELECT 1 FROM items WHERE id = ? LIMIT 1');
    const result = stmt.get(id);

    return result !== undefined;
  }
}
