package com.lebonpoint.unit.shared

import com.lebonpoint.shared.decodeCursor
import com.lebonpoint.shared.encodeCursor
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

/**
 * Unit tests for Cursor codec
 */
class CursorCodecTest {

    @Test
    fun `should encode cursor to base64`() {
        val encoded = encodeCursor(42)
        assertNotNull(encoded)
        assertTrue(encoded.isNotEmpty())
    }

    @Test
    fun `should decode cursor from base64`() {
        val original = encodeCursor(123)
        val decoded = decodeCursor(original)
        val safeDecoded = assertNotNull(decoded)
        assertEquals(123, safeDecoded.id)
    }

    @Test
    fun `should handle round-trip encoding`() {
        val id = 456
        val encoded = encodeCursor(id)
        val decoded = decodeCursor(encoded)
        val safeDecoded = assertNotNull(decoded)
        assertEquals(id, safeDecoded.id)
    }

    @Test
    fun `should return null for invalid cursor`() {
        assertNull(decodeCursor(null))
        assertNull(decodeCursor(""))
        assertNull(decodeCursor("invalid-base64"))
        assertNull(decodeCursor("invalid json"))
    }
}
