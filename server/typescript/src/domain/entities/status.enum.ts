/**
 * Item status enum
 */
export enum Status {
  DRAFT = 'draft',
  ACTIVE = 'active',
  RESERVED = 'reserved',
  SOLD = 'sold',
  ARCHIVED = 'archived'
}

export const STATUS_VALUES = Object.values(Status) as readonly string[];
