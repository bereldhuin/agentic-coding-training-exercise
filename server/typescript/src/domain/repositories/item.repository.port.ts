import type { Item, CreateItemData, UpdateItemData, ReplaceItemData } from '../entities/item.entity.js';

/**
 * Pagination options
 */
export interface PaginationOptions {
  limit: number;
  cursor?: string;
}

/**
 * Filter options for listing items
 */
export interface FilterOptions {
  status?: string;
  category?: string;
  city?: string;
  postal_code?: string;
  is_featured?: boolean;
  delivery_available?: boolean;
}

/**
 * Sort options
 */
export interface SortOptions {
  field: string;
  direction: 'asc' | 'desc';
}

/**
 * Full-text search options
 */
export interface SearchOptions {
  query: string;
  filters?: FilterOptions;
  pagination?: PaginationOptions;
}

/**
 * Paginated result
 */
export interface PaginatedResult<T> {
  items: T[];
  next_cursor?: string;
  total?: number;
}

/**
 * Item repository port interface
 * This is the port that the application layer depends on.
 * Adapters implement this interface to provide persistence.
 */
export interface IItemRepository {
  /**
   * Create a new item
   */
  create(data: Omit<Item, 'id' | 'created_at' | 'updated_at' | 'published_at'>): Promise<Item>;

  /**
   * Find item by ID
   */
  findById(id: number): Promise<Item | null>;

  /**
   * List all items with filters, sorting, and pagination
   */
  findAll(options: {
    filters?: FilterOptions;
    sort?: SortOptions;
    pagination?: PaginationOptions;
  }): Promise<PaginatedResult<Item>>;

  /**
   * Full-text search
   */
  search(options: SearchOptions): Promise<PaginatedResult<Item>>;

  /**
   * Update (replace) an item
   */
  update(id: number, data: Omit<Item, 'id' | 'created_at' | 'updated_at' | 'published_at'>): Promise<Item | null>;

  /**
   * Partial update an item
   */
  patch(id: number, data: Partial<UpdateItemData>): Promise<Item | null>;

  /**
   * Delete an item
   */
  delete(id: number): Promise<boolean>;

  /**
   * Check if an item exists
   */
  exists(id: number): Promise<boolean>;
}
