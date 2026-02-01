/**
 * Error response format
 */
export interface ErrorResponse {
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}

/**
 * Cursor data for pagination
 */
export interface CursorData {
  id: number;
  created_at: string;
}

/**
 * Parse base64 cursor
 */
export function parseCursor(cursor: string): CursorData | null {
  try {
    const decoded = atob(cursor);
    return JSON.parse(decoded) as CursorData;
  } catch {
    return null;
  }
}

/**
 * Create base64 cursor
 */
export function createCursor(data: CursorData): string {
  return btoa(JSON.stringify(data));
}
