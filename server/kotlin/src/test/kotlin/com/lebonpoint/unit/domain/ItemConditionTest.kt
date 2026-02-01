package com.lebonpoint.unit.domain

import com.lebonpoint.domain.valueobjects.ItemCondition
import com.lebonpoint.domain.valueobjects.toDatabaseString
import com.lebonpoint.domain.valueobjects.toItemCondition
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull

/**
 * Unit tests for ItemCondition enum
 */
class ItemConditionTest {

    @Test
    fun `should parse condition from string`() {
        assertEquals(ItemCondition.NEW, ItemCondition.fromString("new"))
        assertEquals(ItemCondition.LIKE_NEW, ItemCondition.fromString("like_new"))
        assertEquals(ItemCondition.GOOD, ItemCondition.fromString("good"))
        assertEquals(ItemCondition.FAIR, ItemCondition.fromString("fair"))
        assertEquals(ItemCondition.PARTS, ItemCondition.fromString("parts"))
        assertEquals(ItemCondition.UNKNOWN, ItemCondition.fromString("unknown"))
    }

    @Test
    fun `should parse condition case-insensitively`() {
        assertEquals(ItemCondition.NEW, ItemCondition.fromString("NEW"))
        assertEquals(ItemCondition.GOOD, ItemCondition.fromString("Good"))
        assertEquals(ItemCondition.LIKE_NEW, ItemCondition.fromString("LIKE_NEW"))
    }

    @Test
    fun `should return null for invalid condition`() {
        assertNull(ItemCondition.fromString("invalid"))
        assertNull(ItemCondition.fromString(""))
        assertNull((null as String?)?.toItemCondition())
    }

    @Test
    fun `should convert to database string`() {
        assertEquals("new", ItemCondition.NEW.toDatabaseString())
        assertEquals("like_new", ItemCondition.LIKE_NEW.toDatabaseString())
        assertEquals("good", ItemCondition.GOOD.toDatabaseString())
    }

    @Test
    fun `should parse from database string`() {
        assertEquals(ItemCondition.NEW, "new".toItemCondition())
        assertEquals(ItemCondition.GOOD, "good".toItemCondition())
        assertNull("invalid".toItemCondition())
    }
}
