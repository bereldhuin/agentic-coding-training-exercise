package com.lebonpoint.unit.application

import com.lebonpoint.application.usecases.CreateItemUseCase
import com.lebonpoint.domain.entities.CreateItemData
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.domain.valueobjects.ItemCondition
import com.lebonpoint.domain.valueobjects.ItemStatus
import com.lebonpoint.shared.ValidationException
import io.mockk.*
import kotlinx.coroutines.test.runTest
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

/**
 * Unit tests for CreateItemUseCase
 */
class CreateItemUseCaseTest {

    @Test
    fun `should create item with valid data`() = runTest {
        // Arrange
        val repository = mockk<ItemRepository>()
        val useCase = CreateItemUseCase(repository)

        val data = CreateItemData(
            title = "Test Item",
            priceCents = 10000,
            condition = ItemCondition.GOOD
        )

        val expectedItem = com.lebonpoint.domain.entities.Item(
            id = 1,
            title = "Test Item",
            description = null,
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

        coEvery { repository.create(any()) } returns expectedItem

        // Act
        val result = useCase.execute(data)

        // Assert
        assertEquals(expectedItem, result)
        coVerify { repository.create(data) }
    }

    @Test
    fun `should throw validation error for empty title`() = runTest {
        // Arrange
        val repository = mockk<ItemRepository>()
        val useCase = CreateItemUseCase(repository)

        val data = CreateItemData(
            title = "",
            priceCents = 10000,
            condition = ItemCondition.GOOD
        )

        // Act & Assert
        val exception = assertFailsWith<ValidationException> {
            useCase.execute(data)
        }

        assertEquals("validation_error", exception.code)
        assertTrue(exception.details.containsKey("title"))
    }

    @Test
    fun `should throw validation error for short title`() = runTest {
        // Arrange
        val repository = mockk<ItemRepository>()
        val useCase = CreateItemUseCase(repository)

        val data = CreateItemData(
            title = "ab", // Less than 3 characters
            priceCents = 10000,
            condition = ItemCondition.GOOD
        )

        // Act & Assert
        assertFailsWith<ValidationException> {
            useCase.execute(data)
        }
    }

    @Test
    fun `should throw validation error for negative price`() = runTest {
        // Arrange
        val repository = mockk<ItemRepository>()
        val useCase = CreateItemUseCase(repository)

        val data = CreateItemData(
            title = "Valid Title",
            priceCents = -100,
            condition = ItemCondition.GOOD
        )

        // Act & Assert
        assertFailsWith<ValidationException> {
            useCase.execute(data)
        }
    }
}
