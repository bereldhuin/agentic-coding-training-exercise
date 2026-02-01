package com.lebonpoint.unit.domain

import com.lebonpoint.domain.entities.CreateItemData
import com.lebonpoint.domain.entities.Item
import com.lebonpoint.domain.entities.PatchItemData
import com.lebonpoint.domain.entities.applyPatch
import com.lebonpoint.domain.valueobjects.ItemCondition
import com.lebonpoint.domain.valueobjects.ItemImage
import com.lebonpoint.domain.valueobjects.ItemStatus
import java.time.Instant
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

/**
 * Unit tests for Item entity
 */
class ItemEntityTest {

    @Test
    fun `should create Item with all fields`() {
        val item = Item(
            id = 1,
            title = "Test Item",
            description = "Test Description",
            priceCents = 10000,
            category = "electronics",
            condition = ItemCondition.NEW,
            status = ItemStatus.ACTIVE,
            isFeatured = true,
            city = "Paris",
            postalCode = "75001",
            country = "FR",
            deliveryAvailable = true,
            createdAt = Instant.EPOCH,
            updatedAt = Instant.EPOCH,
            publishedAt = Instant.EPOCH,
            images = listOf(
                ItemImage(url = "https://example.com/image.jpg", alt = "Test Image", sortOrder = 0)
            )
        )

        assertEquals(1, item.id)
        assertEquals("Test Item", item.title)
        assertEquals("Test Description", item.description)
        assertEquals(10000, item.priceCents)
        assertEquals(ItemCondition.NEW, item.condition)
        assertEquals(ItemStatus.ACTIVE, item.status)
        assertTrue(item.isFeatured)
    }

    @Test
    fun `should create CreateItemData with defaults`() {
        val data = CreateItemData(
            title = "Test Item",
            priceCents = 5000,
            condition = ItemCondition.GOOD
        )

        assertEquals("Test Item", data.title)
        assertEquals(5000, data.priceCents)
        assertEquals(ItemCondition.GOOD, data.condition)
        assertEquals(ItemStatus.DRAFT, data.status)
        assertEquals(false, data.isFeatured)
        assertEquals("FR", data.country)
        assertEquals(false, data.deliveryAvailable)
        assertTrue(data.images.isEmpty())
    }

    @Test
    fun `should apply patch data to item`() {
        val item = Item(
            id = 1,
            title = "Original Title",
            description = "Original Description",
            priceCents = 10000,
            category = null,
            condition = ItemCondition.GOOD,
            status = ItemStatus.DRAFT,
            isFeatured = false,
            city = null,
            postalCode = null,
            country = "FR",
            deliveryAvailable = false,
            createdAt = java.time.Instant.EPOCH,
            updatedAt = java.time.Instant.EPOCH,
            images = emptyList()
        )

        val patch = PatchItemData(
            title = "Updated Title",
            priceCents = 15000,
            status = ItemStatus.ACTIVE
        )

        val updated = item.applyPatch(patch)

        assertEquals("Updated Title", updated.title)
        assertEquals("Original Description", updated.description) // Unchanged
        assertEquals(15000, updated.priceCents)
        assertEquals(ItemStatus.ACTIVE, updated.status)
        assertEquals(ItemCondition.GOOD, updated.condition) // Unchanged
    }
}
