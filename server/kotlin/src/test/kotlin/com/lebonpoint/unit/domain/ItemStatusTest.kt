package com.lebonpoint.unit.domain

import com.lebonpoint.domain.valueobjects.ItemStatus
import com.lebonpoint.domain.valueobjects.toDatabaseString
import com.lebonpoint.domain.valueobjects.toItemStatus
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNull

/**
 * Unit tests for ItemStatus enum
 */
class ItemStatusTest {

    @Test
    fun `should parse status from string`() {
        assertEquals(ItemStatus.DRAFT, ItemStatus.fromString("draft"))
        assertEquals(ItemStatus.ACTIVE, ItemStatus.fromString("active"))
        assertEquals(ItemStatus.RESERVED, ItemStatus.fromString("reserved"))
        assertEquals(ItemStatus.SOLD, ItemStatus.fromString("sold"))
        assertEquals(ItemStatus.ARCHIVED, ItemStatus.fromString("archived"))
    }

    @Test
    fun `should parse status case-insensitively`() {
        assertEquals(ItemStatus.ACTIVE, ItemStatus.fromString("ACTIVE"))
        assertEquals(ItemStatus.DRAFT, ItemStatus.fromString("Draft"))
    }

    @Test
    fun `should return null for invalid status`() {
        assertNull(ItemStatus.fromString("invalid"))
        assertNull(ItemStatus.fromString(""))
        assertNull((null as String?)?.toItemStatus())
    }

    @Test
    fun `should convert to database string`() {
        assertEquals("draft", ItemStatus.DRAFT.toDatabaseString())
        assertEquals("active", ItemStatus.ACTIVE.toDatabaseString())
        assertEquals("sold", ItemStatus.SOLD.toDatabaseString())
    }

    @Test
    fun `should parse from database string`() {
        assertEquals(ItemStatus.DRAFT, "draft".toItemStatus())
        assertEquals(ItemStatus.ACTIVE, "active".toItemStatus())
        assertNull("invalid".toItemStatus())
    }
}
