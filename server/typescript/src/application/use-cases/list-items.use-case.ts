import type { IItemRepository, FilterOptions, SortOptions, PaginationOptions } from '../../domain/repositories/item.repository.port.js';
import { ValidationService } from '../services/validation.service.js';

/**
 * List items use case
 */
export class ListItemsUseCase {
  constructor(private readonly itemRepository: IItemRepository) {}

  /**
   * Execute the use case
   */
  async execute(query: Record<string, unknown>): Promise<{
    items: unknown[];
    next_cursor?: string;
  }> {
    // Validate query parameters
    const validatedQuery = ValidationService.validateListItemsQuery(query);

    // Build filters
    const filters: FilterOptions = {};
    if (validatedQuery.status) filters.status = validatedQuery.status;
    if (validatedQuery.category) filters.category = validatedQuery.category;
    if (validatedQuery.city) filters.city = validatedQuery.city;
    if (validatedQuery.postal_code) filters.postal_code = validatedQuery.postal_code;
    if (validatedQuery.is_featured !== undefined) filters.is_featured = validatedQuery.is_featured;
    if (validatedQuery.delivery_available !== undefined) filters.delivery_available = validatedQuery.delivery_available;

    // Build sort options
    let sort: SortOptions | undefined;
    if (validatedQuery.sort) {
      const [field, direction] = validatedQuery.sort.split(':');
      sort = { field, direction: direction as 'asc' | 'desc' };
    }

    // Build pagination options
    const pagination: PaginationOptions = {
      limit: validatedQuery.limit
    };
    if (validatedQuery.cursor) {
      pagination.cursor = validatedQuery.cursor;
    }

    // Check if this is a full-text search
    if (validatedQuery.q) {
      const result = await this.itemRepository.search({
        query: validatedQuery.q,
        filters,
        pagination
      });
      return {
        items: result.items,
        next_cursor: result.next_cursor
      };
    }

    // Regular list with filters
    const result = await this.itemRepository.findAll({
      filters,
      sort,
      pagination
    });

    return {
      items: result.items,
      next_cursor: result.next_cursor
    };
  }
}
