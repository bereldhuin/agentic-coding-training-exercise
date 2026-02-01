package com.lebonpoint.shared

import kotlinx.serialization.Serializable
import kotlinx.serialization.encodeToString
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.json.Json
import java.util.Base64

/**
 * Cursor data for pagination
 */
@Serializable
data class CursorData(
    val id: Int
)

/**
 * Encode cursor data to base64 string
 */
fun encodeCursor(id: Int): String {
    val data = CursorData(id)
    val json = Json.encodeToString(data)
    return Base64.getEncoder().encodeToString(json.toByteArray())
}

/**
 * Decode cursor string to cursor data
 * Returns null for invalid cursors
 */
fun decodeCursor(cursor: String?): CursorData? {
    if (cursor.isNullOrBlank()) return null

    return try {
        val json = String(Base64.getDecoder().decode(cursor))
        Json.decodeFromString<CursorData>(json)
    } catch (e: Exception) {
        null
    }
}
