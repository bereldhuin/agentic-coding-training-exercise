package com.lebonpoint.domain.valueobjects

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

/**
 * Item status enum
 * Represents the lifecycle status of a marketplace item
 */
@Serializable
enum class ItemStatus {
    @SerialName("draft")
    DRAFT,
    @SerialName("active")
    ACTIVE,
    @SerialName("reserved")
    RESERVED,
    @SerialName("sold")
    SOLD,
    @SerialName("archived")
    ARCHIVED;

    companion object {
        /**
         * Parse from string (snake_case or lowercase)
         */
        fun fromString(value: String): ItemStatus? {
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
fun ItemStatus.toDatabaseString(): String {
    return name.lowercase()
}

/**
 * Parse ItemStatus from database string
 */
fun String.toItemStatus(): ItemStatus? {
    return ItemStatus.fromString(this)
}
