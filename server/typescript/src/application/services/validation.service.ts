import { z } from 'zod';
import { CONDITION_VALUES, STATUS_VALUES } from '../../domain/entities/index.js';
import type { ItemImage } from '../../domain/value-objects/item-image.vo.js';
import { ValidationError } from '../../shared/errors.js';

/**
 * Item image schema
 */
const itemImageSchema: z.ZodType<ItemImage> = z.object({
  url: z.string().min(1, 'Image URL is required'),
  alt: z.string().optional(),
  sort_order: z.number().int().optional()
});

/**
 * Create item schema
 */
export const createItemSchema = z.object({
  title: z.string().min(3, 'Title must be at least 3 characters').max(200, 'Title must be at most 200 characters'),
  description: z.string().optional(),
  price_cents: z.number().int().nonnegative('Price must be non-negative'),
  category: z.string().optional(),
  condition: z.enum(CONDITION_VALUES as readonly [string, ...string[]]),
  status: z.enum(STATUS_VALUES as readonly [string, ...string[]]).optional().default('draft'),
  is_featured: z.boolean().optional().default(false),
  city: z.string().optional(),
  postal_code: z.string().optional(),
  country: z.string().optional().default('FR'),
  delivery_available: z.boolean().optional().default(false),
  images: z.array(itemImageSchema).optional().default([])
});

/**
 * Update (PUT) item schema - all fields except id are required
 */
export const updateItemSchema = z.object({
  title: z.string().min(3, 'Title must be at least 3 characters').max(200, 'Title must be at most 200 characters'),
  description: z.string().optional(),
  price_cents: z.number().int().nonnegative('Price must be non-negative'),
  category: z.string().optional(),
  condition: z.enum(CONDITION_VALUES as readonly [string, ...string[]]),
  status: z.enum(STATUS_VALUES as readonly [string, ...string[]]),
  is_featured: z.boolean(),
  city: z.string().optional(),
  postal_code: z.string().optional(),
  country: z.string(),
  delivery_available: z.boolean(),
  images: z.array(itemImageSchema).optional().default([])
});

/**
 * Patch item schema - all fields optional
 */
export const patchItemSchema = z.object({
  title: z.string().min(3, 'Title must be at least 3 characters').max(200, 'Title must be at most 200 characters').optional(),
  description: z.string().optional(),
  price_cents: z.number().int().nonnegative('Price must be non-negative').optional(),
  category: z.string().optional(),
  condition: z.enum(CONDITION_VALUES as readonly [string, ...string[]]).optional(),
  status: z.enum(STATUS_VALUES as readonly [string, ...string[]]).optional(),
  is_featured: z.boolean().optional(),
  city: z.string().optional(),
  postal_code: z.string().optional(),
  country: z.string().optional(),
  delivery_available: z.boolean().optional(),
  images: z.array(itemImageSchema).optional()
}).refine(data => Object.keys(data).length > 0, {
  message: 'At least one field must be provided for patch'
});

/**
 * Query parameter schemas for listing items
 */
export const listItemsQuerySchema = z.object({
  status: z.enum(STATUS_VALUES as readonly [string, ...string[]]).optional(),
  category: z.string().optional(),
  min_price_cents: z.coerce.number().int().nonnegative().optional(),
  max_price_cents: z.coerce.number().int().nonnegative().optional(),
  city: z.string().optional(),
  postal_code: z.string().optional(),
  is_featured: z.coerce.boolean().optional(),
  delivery_available: z.coerce.boolean().optional(),
  sort: z.string().regex(/^\w+:(asc|desc)$/, 'Sort must be in format "field:direction"').optional(),
  limit: z.coerce.number().int().min(1).max(100).optional().default(20),
  cursor: z.string().optional(),
  q: z.string().optional() // Full-text search query
});

/**
 * Validation service
 */
export class ValidationService {
  /**
   * Validate create item data
   */
  static validateCreateItem(data: unknown): z.infer<typeof createItemSchema> {
    try {
      return createItemSchema.parse(data);
    } catch (error) {
      if (error instanceof z.ZodError) {
        throw new ValidationError('Validation failed', this.formatZodError(error));
      }
      throw error;
    }
  }

  /**
   * Validate update item data
   */
  static validateUpdateItem(data: unknown): z.infer<typeof updateItemSchema> {
    try {
      return updateItemSchema.parse(data);
    } catch (error) {
      if (error instanceof z.ZodError) {
        throw new ValidationError('Validation failed', this.formatZodError(error));
      }
      throw error;
    }
  }

  /**
   * Validate patch item data
   */
  static validatePatchItem(data: unknown): z.infer<typeof patchItemSchema> {
    try {
      return patchItemSchema.parse(data);
    } catch (error) {
      if (error instanceof z.ZodError) {
        throw new ValidationError('Validation failed', this.formatZodError(error));
      }
      throw error;
    }
  }

  /**
   * Validate list items query parameters
   */
  static validateListItemsQuery(data: unknown): z.infer<typeof listItemsQuerySchema> {
    try {
      return listItemsQuerySchema.parse(data);
    } catch (error) {
      if (error instanceof z.ZodError) {
        throw new ValidationError('Validation failed', this.formatZodError(error));
      }
      throw error;
    }
  }

  /**
   * Format Zod error to details object
   */
  private static formatZodError(error: z.ZodError): Record<string, unknown> {
    const details: Record<string, unknown> = {};
    const issues = error.issues as Array<{ path: (string | number)[]; message: string }>;
    issues.forEach((err) => {
      const path = err.path.join('.');
      details[path] = err.message;
    });
    return details;
  }
}
