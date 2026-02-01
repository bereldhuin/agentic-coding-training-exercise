import { Condition } from './condition.enum.js';
import { Status } from './status.enum.js';
import type { ItemImage } from '../value-objects/item-image.vo.js';
import type { Location } from '../value-objects/location.vo.js';

/**
 * Item entity
 */
export interface Item {
  id: number;
  title: string;
  description?: string;
  price_cents: number;
  category?: string;
  condition: Condition;
  status: Status;
  is_featured: boolean;
  city?: string;
  postal_code?: string;
  country: string;
  delivery_available: boolean;
  created_at: string;
  updated_at: string;
  published_at?: string | null;
  images: ItemImage[];
}

/**
 * Item data for creation (without id and timestamps)
 */
export interface CreateItemData {
  title: string;
  description?: string;
  price_cents: number;
  category?: string;
  condition: Condition;
  status?: Status;
  is_featured?: boolean;
  city?: string;
  postal_code?: string;
  country?: string;
  delivery_available?: boolean;
  images?: ItemImage[];
}

/**
 * Item data for update (all optional except id)
 */
export interface UpdateItemData {
  title?: string;
  description?: string;
  price_cents?: number;
  category?: string;
  condition?: Condition;
  status?: Status;
  is_featured?: boolean;
  city?: string;
  postal_code?: string;
  country?: string;
  delivery_available?: boolean;
  images?: ItemImage[];
}

/**
 * Item data for PUT (full replace)
 */
export interface ReplaceItemData {
  title: string;
  description?: string;
  price_cents: number;
  category?: string;
  condition: Condition;
  status: Status;
  is_featured: boolean;
  city?: string;
  postal_code?: string;
  country: string;
  delivery_available: boolean;
  images?: ItemImage[];
}

/**
 * Create item with defaults
 */
export function createItemWithDefaults(data: CreateItemData): Omit<Item, 'id' | 'created_at' | 'updated_at' | 'published_at'> {
  return {
    title: data.title,
    description: data.description,
    price_cents: data.price_cents,
    category: data.category,
    condition: data.condition,
    status: data.status ?? Status.DRAFT,
    is_featured: data.is_featured ?? false,
    city: data.city,
    postal_code: data.postal_code,
    country: data.country ?? 'FR',
    delivery_available: data.delivery_available ?? false,
    images: data.images ?? []
  };
}

/**
 * Convert database row to Item entity
 */
export function itemFromRow(row: Record<string, unknown>): Item {
  return {
    id: row.id as number,
    title: row.title as string,
    description: (row.description as string | undefined) ?? undefined,
    price_cents: row.price_cents as number,
    category: (row.category as string | undefined) ?? undefined,
    condition: row.condition as Condition,
    status: row.status as Status,
    is_featured: Boolean(row.is_featured),
    city: (row.city as string | undefined) ?? undefined,
    postal_code: (row.postal_code as string | undefined) ?? undefined,
    country: row.country as string,
    delivery_available: Boolean(row.delivery_available),
    created_at: row.created_at as string,
    updated_at: row.updated_at as string,
    published_at: (row.published_at as string | null | undefined) ?? undefined,
    images: parseImagesColumn(row.images as string | null)
  };
}

/**
 * Parse images from JSON column
 */
function parseImagesColumn(value: string | null): ItemImage[] {
  if (!value) {
    return [];
  }

  try {
    const parsed = JSON.parse(value);
    if (Array.isArray(parsed)) {
      return parsed;
    }
    return [];
  } catch {
    return [];
  }
}

/**
 * Serialize images to JSON column
 */
export function serializeImagesColumn(images: ItemImage[]): string {
  return JSON.stringify(images);
}
