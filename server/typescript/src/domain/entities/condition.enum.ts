/**
 * Item condition enum
 */
export enum Condition {
  NEW = 'new',
  LIKE_NEW = 'like_new',
  GOOD = 'good',
  FAIR = 'fair',
  PARTS = 'parts',
  UNKNOWN = 'unknown'
}

export const CONDITION_VALUES = Object.values(Condition) as readonly string[];
