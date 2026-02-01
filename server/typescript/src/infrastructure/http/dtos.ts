import type { Item } from '../../domain/entities/item.entity.js';

/**
 * Item response DTO
 */
export interface ItemResponseDTO {
  id: number;
  title: string;
  description?: string;
  price_cents: number;
  category?: string;
  condition: string;
  status: string;
  is_featured: boolean;
  city?: string;
  postal_code?: string;
  country: string;
  delivery_available: boolean;
  created_at: string;
  updated_at: string;
  published_at?: string | null;
  images: Array<{
    url: string;
    alt?: string;
    sort_order?: number;
  }>;
}

/**
 * Create item request DTO
 */
export interface CreateItemRequestDTO {
  title: string;
  description?: string;
  price_cents: number;
  category?: string;
  condition: string;
  status?: string;
  is_featured?: boolean;
  city?: string;
  postal_code?: string;
  country?: string;
  delivery_available?: boolean;
  images?: Array<{
    url: string;
    alt?: string;
    sort_order?: number;
  }>;
}

/**
 * Update item request DTO (PUT)
 */
export interface UpdateItemRequestDTO {
  title: string;
  description?: string;
  price_cents: number;
  category?: string;
  condition: string;
  status: string;
  is_featured: boolean;
  city?: string;
  postal_code?: string;
  country: string;
  delivery_available: boolean;
  images?: Array<{
    url: string;
    alt?: string;
    sort_order?: number;
  }>;
}

/**
 * Patch item request DTO
 */
export interface PatchItemRequestDTO {
  title?: string;
  description?: string;
  price_cents?: number;
  category?: string;
  condition?: string;
  status?: string;
  is_featured?: boolean;
  city?: string;
  postal_code?: string;
  country?: string;
  delivery_available?: boolean;
  images?: Array<{
    url: string;
    alt?: string;
    sort_order?: number;
  }>;
}

/**
 * List items response DTO
 */
export interface ListItemsResponseDTO {
  items: ItemResponseDTO[];
  next_cursor?: string;
}

/**
 * Error response DTO
 */
export interface ErrorResponseDTO {
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}

/**
 * Convert Item entity to ItemResponseDTO
 */
export function toItemResponseDTO(item: Item): ItemResponseDTO {
  return {
    id: item.id,
    title: item.title,
    description: item.description,
    price_cents: item.price_cents,
    category: item.category,
    condition: item.condition,
    status: item.status,
    is_featured: item.is_featured,
    city: item.city,
    postal_code: item.postal_code,
    country: item.country,
    delivery_available: item.delivery_available,
    created_at: item.created_at,
    updated_at: item.updated_at,
    published_at: item.published_at,
    images: item.images
  };
}
