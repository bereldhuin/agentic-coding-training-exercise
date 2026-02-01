package com.lebonpoint.domain.valueobjects

import kotlinx.serialization.Serializable

/**
 * Item condition enum
 * Represents the physical condition of a marketplace item
 */
@Serializable
enum class ItemCondition {
    NEW,
    LIKE_NEW,
    GOOD,
    FAIR,
    PARTS,
    UNKNOWN;

    companion object {
        /**
         * Parse from string (snake_case or lowercase)
         */
        fun fromString(value: String): ItemCondition? {
            return entries.find {
                it.name.equals(value, ignoreCase = true) ||
                it.name.lowercase().replace("_", "") == value.lowercase().replace("_", "")
            }
        }
    }
}

/**
 * Convert to snake_case string for database storage
 */
fun ItemCondition.toDatabaseString(): String {
    return name.lowercase()
}

/**
 * Parse ItemCondition from database string
 */
fun String.toItemCondition(): ItemCondition? {
    return ItemCondition.fromString(this)
}
