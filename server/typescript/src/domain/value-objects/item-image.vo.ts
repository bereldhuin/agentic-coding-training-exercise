/**
 * Item image value object
 */
export interface ItemImage {
  url: string;
  alt?: string;
  sort_order?: number;
}

/**
 * Validate item image
 */
export function validateItemImage(image: unknown): image is ItemImage {
  if (typeof image !== 'object' || image === null) {
    return false;
  }

  const img = image as Record<string, unknown>;

  if (typeof img.url !== 'string' || img.url.trim().length === 0) {
    return false;
  }

  if (img.alt !== undefined && typeof img.alt !== 'string') {
    return false;
  }

  if (img.sort_order !== undefined && typeof img.sort_order !== 'number') {
    return false;
  }

  return true;
}
