/**
 * Location value object
 */
export interface Location {
  city?: string;
  postal_code?: string;
  country: string;
}

/**
 * Create a location with defaults
 */
export function createLocation(input?: Partial<Location>): Location {
  return {
    city: input?.city,
    postal_code: input?.postal_code,
    country: input?.country ?? 'FR'
  };
}
