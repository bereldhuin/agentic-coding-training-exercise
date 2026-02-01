package com.lebonpoint.integration.repository

import com.lebonpoint.domain.entities.CreateItemData
import com.lebonpoint.domain.entities.Item
import com.lebonpoint.domain.entities.PatchItemData
import com.lebonpoint.domain.entities.ReplaceItemData
import com.lebonpoint.domain.repositories.ItemFilters
import com.lebonpoint.domain.repositories.SortDirection
import com.lebonpoint.domain.repositories.SortOptions
import com.lebonpoint.infrastructure.persistence.DatabaseConfig
import com.lebonpoint.infrastructure.persistence.SQLiteItemRepository
import com.lebonpoint.domain.valueobjects.ItemCondition
import com.lebonpoint.domain.valueobjects.ItemImage
import com.lebonpoint.domain.valueobjects.ItemStatus
import kotlinx.coroutines.runBlocking
import kotlinx.coroutines.test.runTest
import org.junit.After
import org.junit.Before
import org.junit.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

/**
 * Integration tests for SQLiteItemRepository
 * Uses in-memory SQLite database
 */
class SQLiteItemRepositoryTest {

    private lateinit var repository: SQLiteItemRepository

    @Before
    fun setup() = runBlocking {
        // Initialize in-memory database
        System.setProperty("DATABASE_PATH", ":memory:")
        DatabaseConfig.initializeSchema()
        repository = SQLiteItemRepository()
    }

    @After
    fun teardown() {
        DatabaseConfig.shutdown()
    }

    @Test
    fun `should create and retrieve item`() = runTest {
        // Arrange
        val data = CreateItemData(
            title = "Test Item",
            description = "Test Description",
            priceCents = 10000,
            category = "electronics",
            condition = ItemCondition.GOOD,
            images = listOf(ItemImage(url = "https://example.com/image.jpg", alt = "Test"))
        )

        // Act
        val created = repository.create(data)
        val retrieved = repository.findById(created.id)

        // Assert
        assertNotNull(retrieved)
        assertEquals(data.title, retrieved.title)
        assertEquals(data.description, retrieved.description)
        assertEquals(data.priceCents, retrieved.priceCents)
        assertEquals(data.category, retrieved.category)
        assertEquals(data.condition, retrieved.condition)
        assertEquals(1, retrieved.images.size)
        assertEquals("https://example.com/image.jpg", retrieved.images[0].url)
    }

    @Test
    fun `should return null for non-existent item`() = runTest {
        // Act
        val item = repository.findById(999)

        // Assert
        assertNull(item)
    }

    @Test
    fun `should find all items with no filters`() = runTest {
        // Arrange
        repository.create(CreateItemData(title = "Item 1", priceCents = 1000, condition = ItemCondition.NEW))
        repository.create(CreateItemData(title = "Item 2", priceCents = 2000, condition = ItemCondition.GOOD))
        repository.create(CreateItemData(title = "Item 3", priceCents = 3000, condition = ItemCondition.FAIR))

        // Act
        val page = repository.findAll()

        // Assert
        assertEquals(3, page.items.size)
    }

    @Test
    fun `should filter items by status`() = runTest {
        // Arrange
        repository.create(
            CreateItemData(
                title = "Item 1",
                priceCents = 1000,
                condition = ItemCondition.NEW,
                status = ItemStatus.ACTIVE
            )
        )
        repository.create(
            CreateItemData(
                title = "Item 2",
                priceCents = 2000,
                condition = ItemCondition.GOOD,
                status = ItemStatus.DRAFT
            )
        )
        repository.create(
            CreateItemData(
                title = "Item 3",
                priceCents = 3000,
                condition = ItemCondition.FAIR,
                status = ItemStatus.ACTIVE
            )
        )

        // Act
        val page = repository.findAll(
            filters = ItemFilters(status = "active")
        )

        // Assert
        assertEquals(2, page.items.size)
        assertTrue(page.items.all { it.status == ItemStatus.ACTIVE })
    }

    @Test
    fun `should filter items by price range`() = runTest {
        // Arrange
        repository.create(
            CreateItemData(
                title = "Item 1",
                priceCents = 1000,
                condition = ItemCondition.NEW
            )
        )
        repository.create(
            CreateItemData(
                title = "Item 2",
                priceCents = 5000,
                condition = ItemCondition.GOOD
            )
        )
        repository.create(
            CreateItemData(
                title = "Item 3",
                priceCents = 10000,
                condition = ItemCondition.FAIR
            )
        )

        // Act
        val page = repository.findAll(
            filters = ItemFilters(minPriceCents = 2000, maxPriceCents = 8000)
        )

        // Assert
        assertEquals(1, page.items.size)
        assertEquals(5000, page.items[0].priceCents)
    }

    @Test
    fun `should sort items by price descending`() = runTest {
        // Arrange
        repository.create(
            CreateItemData(
                title = "Item 1",
                priceCents = 1000,
                condition = ItemCondition.NEW
            )
        )
        repository.create(
            CreateItemData(
                title = "Item 2",
                priceCents = 5000,
                condition = ItemCondition.GOOD
            )
        )
        repository.create(
            CreateItemData(
                title = "Item 3",
                priceCents = 3000,
                condition = ItemCondition.FAIR
            )
        )

        // Act
        val page = repository.findAll(
            sort = SortOptions(field = "price_cents", direction = SortDirection.DESC)
        )

        // Assert
        assertEquals(3, page.items.size)
        assertEquals(5000, page.items[0].priceCents)
        assertEquals(3000, page.items[1].priceCents)
        assertEquals(1000, page.items[2].priceCents)
    }

    @Test
    fun `should paginate items with cursor`() = runTest {
        // Arrange
        for (i in 1..25) {
            repository.create(
                CreateItemData(
                    title = "Item $i",
                    priceCents = i * 100,
                    condition = ItemCondition.GOOD
                )
            )
        }

        // Act - First page
        val page1 = repository.findAll(limit = 10)

        // Assert
        assertEquals(10, page1.items.size)
        assertNotNull(page1.nextCursor)

        // Act - Second page
        val page2 = repository.findAll(limit = 10, cursor = page1.nextCursor)

        // Assert
        assertEquals(10, page2.items.size)
        assertTrue(page2.items[0].id > page1.items.last().id)
    }

    @Test
    fun `should search items with full-text search`() = runTest {
        // Arrange
        repository.create(
            CreateItemData(
                title = "iPhone 13 Pro",
                priceCents = 65000,
                condition = ItemCondition.GOOD,
                description = "Excellent condition"
            )
        )
        repository.create(
            CreateItemData(
                title = "Samsung Galaxy",
                priceCents = 40000,
                condition = ItemCondition.NEW,
                description = "Brand new"
            )
        )
        repository.create(
            CreateItemData(
                title = "MacBook Pro",
                priceCents = 120000,
                condition = ItemCondition.LIKE_NEW,
                description = "Like new"
            )
        )

        // Act
        val page = repository.search(query = "iPhone")

        // Assert
        assertEquals(1, page.items.size)
        assertEquals("iPhone 13 Pro", page.items[0].title)
    }

    @Test
    fun `should update item`() = runTest {
        // Arrange
        val created = repository.create(
            CreateItemData(
                title = "Original Title",
                priceCents = 10000,
                condition = ItemCondition.GOOD
            )
        )

        // Act
        val updated = repository.update(
            created.id,
            ReplaceItemData(
                title = "Updated Title",
                description = "Updated Description",
                priceCents = 15000,
                category = "electronics",
                condition = ItemCondition.LIKE_NEW,
                status = ItemStatus.ACTIVE,
                isFeatured = true,
                city = "Paris",
                postalCode = "75001",
                country = "FR",
                deliveryAvailable = true,
                images = emptyList()
            )
        )

        // Assert
        assertNotNull(updated)
        assertEquals("Updated Title", updated.title)
        assertEquals(15000, updated.priceCents)
        assertEquals(ItemCondition.LIKE_NEW, updated.condition)
    }

    @Test
    fun `should patch item`() = runTest {
        // Arrange
        val created = repository.create(
            CreateItemData(
                title = "Title",
                priceCents = 10000,
                condition = ItemCondition.GOOD
            )
        )

        // Act
        val patched = repository.patch(
            created.id,
            PatchItemData(title = "Patched Title", priceCents = 20000)
        )

        // Assert
        assertNotNull(patched)
        assertEquals("Patched Title", patched.title)
        assertEquals(20000, patched.priceCents)
        assertEquals(ItemCondition.GOOD, patched.condition) // Unchanged
    }

    @Test
    fun `should delete item`() = runTest {
        // Arrange
        val created = repository.create(
            CreateItemData(
                title = "To Delete",
                priceCents = 10000,
                condition = ItemCondition.GOOD
            )
        )

        // Act
        val deleted = repository.delete(created.id)

        // Assert
        assertTrue(deleted)
        assertNull(repository.findById(created.id))
    }

    @Test
    fun `should return false when deleting non-existent item`() = runTest {
        // Act
        val deleted = repository.delete(999)

        // Assert
        assertEquals(false, deleted)
    }
}
