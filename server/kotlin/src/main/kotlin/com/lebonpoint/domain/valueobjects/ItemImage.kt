package com.lebonpoint.domain.valueobjects

import kotlinx.serialization.Serializable

/**
 * Item image value object
 * Represents a single image associated with a marketplace item
 */
@Serializable
data class ItemImage(
    var url: String,
    var alt: String? = null,
    var sortOrder: Int = 0
)

/**
 * Serialize list of ItemImage to JSON string for database storage
 */
fun List<ItemImage>.serializeToJson(): String {
    return kotlinx.serialization.json.Json.encodeToString(
        kotlinx.serialization.builtins.ListSerializer(ItemImage.serializer()),
        this
    )
}

/**
 * Parse JSON string to list of ItemImage
 */
fun String.deserializeToImages(): List<ItemImage> {
    return if (this.isBlank()) {
        emptyList()
    } else {
        try {
            kotlinx.serialization.json.Json.decodeFromString(
                kotlinx.serialization.builtins.ListSerializer(ItemImage.serializer()),
                this
            )
        } catch (e: Exception) {
            emptyList()
        }
    }
}
